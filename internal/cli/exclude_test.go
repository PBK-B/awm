package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func runTestGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func inTestDir(t *testing.T, dir string) {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(cwd); err != nil {
			t.Fatal(err)
		}
	})
}

func testReadFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestGitExcludePreservesContentAndIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	runTestGit(t, dir, "init")
	inTestDir(t, dir)
	path := runTestGit(t, dir, "rev-parse", "--path-format=absolute", "--git-path", "info/exclude")
	// Comments and similar paths must not be mistaken for existing rules.
	original := "# .awm-metadata.local.json\n.agents/tmp/backup/\ncustom-file"
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}
	if err := ensureGitExclude(); err != nil {
		t.Fatal(err)
	}
	want := original + "\n\n# awm local files\n.awm-metadata.local.json\n.agents/tmp/\n"
	if got := testReadFile(t, path); got != want {
		t.Fatalf("exclude = %q, want %q", got, want)
	}
	if err := ensureGitExclude(); err != nil {
		t.Fatal(err)
	}
	if got := testReadFile(t, path); got != want {
		t.Fatalf("second call changed exclude: %q", got)
	}
	if err := os.WriteFile(path, []byte(".awm-metadata.local.json\r\n.agents/tmp/\r\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := ensureGitExclude(); err != nil {
		t.Fatal(err)
	}
	if got := testReadFile(t, path); got != ".awm-metadata.local.json\r\n.agents/tmp/\r\n" {
		t.Fatalf("existing CRLF rules changed: %q", got)
	}
}

func TestInitAndRestoreLocalExcludes(t *testing.T) {
	for _, existingIgnore := range []bool{false, true} {
		name := "without_gitignore"
		if existingIgnore {
			name = "with_gitignore"
		}
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			runTestGit(t, dir, "init")
			inTestDir(t, dir)
			if existingIgnore {
				if err := os.WriteFile(".gitignore", []byte("custom-file"), 0644); err != nil {
					t.Fatal(err)
				}
			}
			if err := cmdInit([]string{"--name", "test"}); err != nil {
				t.Fatal(err)
			}
			assertLocalFilesIgnored(t, dir)
			path := runTestGit(t, dir, "rev-parse", "--path-format=absolute", "--git-path", "info/exclude")
			for _, command := range []string{"install", "restore"} {
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				if err := Run([]string{command}); err != nil {
					t.Fatal(err)
				}
				assertLocalFilesIgnored(t, dir)
			}
			if existingIgnore {
				if got := testReadFile(t, ".gitignore"); got != "custom-file" {
					t.Fatalf("gitignore changed: %q", got)
				}
			} else if _, err := os.Stat(".gitignore"); !os.IsNotExist(err) {
				t.Fatalf("unexpected gitignore: %v", err)
			}
		})
	}
}

func assertLocalFilesIgnored(t *testing.T, dir string) {
	t.Helper()
	for _, path := range []string{".awm-metadata.local.json", ".agents/tmp/test.txt"} {
		if got := runTestGit(t, dir, "check-ignore", path); got != path {
			t.Fatalf("check-ignore = %q, want %q", got, path)
		}
	}
}

func TestGitExcludeLinkedRepositories(t *testing.T) {
	for _, layout := range []string{"worktree", "submodule"} {
		t.Run(layout, func(t *testing.T) {
			parent := t.TempDir()
			repo := filepath.Join(parent, "repo")
			if err := os.Mkdir(repo, 0755); err != nil {
				t.Fatal(err)
			}
			runTestGit(t, repo, "init")
			runTestGit(t, repo, "-c", "user.name=Test", "-c", "user.email=test@example.com", "commit", "--allow-empty", "-m", "initial")
			var linked string
			if layout == "worktree" {
				linked = filepath.Join(parent, "linked")
				runTestGit(t, repo, "worktree", "add", "--detach", linked)
			} else {
				linked = filepath.Join(repo, "child")
				runTestGit(t, repo, "-c", "protocol.file.allow=always", "submodule", "add", repo, "child")
			}
			inTestDir(t, linked)
			if err := ensureGitExclude(); err != nil {
				t.Fatal(err)
			}
			assertLocalFilesIgnored(t, linked)
			path := runTestGit(t, linked, "rev-parse", "--path-format=absolute", "--git-path", "info/exclude")
			if !strings.Contains(testReadFile(t, path), "# awm local files") {
				t.Fatal("rules missing from Git-resolved exclude")
			}
		})
	}
}

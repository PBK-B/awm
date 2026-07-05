package editor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDetectPrefersExplicit(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "fake-editor")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
	t.Setenv("AWM_EDITOR", "missing-editor")
	got, err := Detect("fake-editor")
	if err != nil {
		t.Fatal(err)
	}
	if got != bin {
		t.Fatalf("Detect = %q, want %q", got, bin)
	}
}

func TestDetectExplicitMissingDoesNotFallback(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "env-editor")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
	t.Setenv("AWM_EDITOR", "env-editor")

	if _, err := Detect("missing-editor"); err == nil {
		t.Fatal("expected missing explicit editor to return an error")
	}
}

func TestDetectUsesEnvironment(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "env-editor")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
	t.Setenv("AWM_EDITOR", "env-editor")
	got, err := Detect("")
	if err != nil {
		t.Fatal(err)
	}
	if got != bin {
		t.Fatalf("Detect = %q, want %q", got, bin)
	}
}

func TestOpenWaitsForEditor(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "slow-editor")
	marker := filepath.Join(dir, "done")
	script := "#!/bin/sh\nsleep 0.1\ntouch \"" + marker + "\"\n"
	if err := os.WriteFile(bin, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	start := time.Now()
	if err := Open("slow-editor"); err != nil {
		t.Fatal(err)
	}
	if time.Since(start) < 100*time.Millisecond {
		t.Fatal("Open returned before editor process completed")
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("marker not created: %v", err)
	}
}

func TestOpenReportsEditorFailure(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "fail-editor")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\necho tool failed >&2\nexit 42\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	err := Open("fail-editor")
	if err == nil {
		t.Fatal("expected editor failure")
	}
	msg := err.Error()
	if !strings.Contains(msg, `awm open: editor "fail-editor" failed`) {
		t.Fatalf("error missing editor context: %s", msg)
	}
	if !strings.Contains(msg, "exit status 42") {
		t.Fatalf("error missing exit status: %s", msg)
	}
}

func TestOpenPassesArgsAndEnv(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "capture-editor")
	out := filepath.Join(dir, "out")
	script := "#!/bin/sh\nprintf 'env=%s args=%s\\n' \"$AWM_TEST_VALUE\" \"$*\" > \"" + out + "\"\n"
	if err := os.WriteFile(bin, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	err := OpenWithOptions(OpenOptions{Editor: "capture-editor", Args: []string{"--foo", "bar"}, Env: []string{"AWM_TEST_VALUE=ok"}})
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if !strings.Contains(got, "env=ok") || !strings.Contains(got, "args=--foo bar") {
		t.Fatalf("unexpected captured output: %q", got)
	}
}

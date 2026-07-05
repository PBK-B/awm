package cli

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"awmcli/internal/editor"
	"awmcli/internal/gitutil"
	"awmcli/internal/templates"
	"awmcli/internal/workspace"
)

func cmdInit(args []string) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	name := fs.String("name", filepath.Base(mustCwd()), "workspace name")
	openAfter := fs.Bool("open", false, "open editor after init")
	ed := fs.String("editor", "", "editor to open after init")
	yes := fs.Bool("yes", false, "answer yes to git init prompt")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureGitRepo(*yes); err != nil {
		return err
	}
	if err := ensureDirs(); err != nil {
		return err
	}
	if err := ensureInitFiles(*name); err != nil {
		return err
	}
	fmt.Println("Workspace initialized.")
	if *openAfter || *ed != "" {
		return editor.Open(*ed)
	}
	return nil
}

func ensureGitRepo(autoYes bool) error {
	if gitutil.IsRepo() {
		return nil
	}
	if !autoYes {
		ok, err := confirm("Current directory is not a git repository. Initialize git repository here? [y/N]: ")
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("aborted: awm requires a git repository")
		}
	}
	return gitutil.Run("init")
}

func confirm(prompt string) (bool, error) {
	fmt.Print(prompt)
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return false, err
	}
	s := strings.ToLower(strings.TrimSpace(input))
	switch s {
	case "y", "yes":
		return true, nil
	case "", "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid input: %s", s)
	}
}

func ensureDirs() error {
	for _, dir := range []string{workspace.SkillsDir, workspace.PatchesDir, workspace.WorkspaceDir, workspace.ArrangeDir, workspace.TmpDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func ensureInitFiles(name string) error {
	if !workspace.Exists(workspace.MetadataPath) {
		if err := workspace.WriteMetadata(workspace.DefaultMetadata(name)); err != nil {
			return err
		}
	}
	if !workspace.Exists(workspace.LockPath) {
		if err := workspace.WriteLock(workspace.DefaultLock(name)); err != nil {
			return err
		}
	}
	files := map[string]string{
		workspace.ReadmeFile:    templates.Readme(name),
		workspace.AgentsFile:    templates.Agents(name),
		workspace.PatchesReadme: templates.PatchesReadme(),
	}
	for path, content := range files {
		if !workspace.Exists(path) {
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				return err
			}
		}
	}
	return ensureGitignore()
}

func ensureGitignore() error {
	if !workspace.Exists(workspace.GitignorePath) {
		return os.WriteFile(workspace.GitignorePath, []byte(templates.Gitignore()), 0644)
	}
	b, err := os.ReadFile(workspace.GitignorePath)
	if err != nil {
		return err
	}
	s := string(b)
	appendLines := []string{}
	for _, line := range []string{workspace.LocalPath, workspace.TmpDir + "/"} {
		if !strings.Contains(s, line) {
			appendLines = append(appendLines, line)
		}
	}
	if len(appendLines) == 0 {
		return nil
	}
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	s += "\n# awm local files\n" + strings.Join(appendLines, "\n") + "\n"
	return os.WriteFile(workspace.GitignorePath, []byte(s), 0644)
}

func mustCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "workspace"
	}
	return cwd
}

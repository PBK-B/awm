package editor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"awmcli/internal/workspace"
)

type OpenOptions struct {
	Editor string
	Args   []string
	Env    []string
}

func Open(preferred string) error {
	return OpenWithOptions(OpenOptions{Editor: preferred})
}

func OpenWithOptions(opts OpenOptions) error {
	ed, err := Detect(opts.Editor)
	if err != nil {
		return err
	}
	args := opts.Args
	if len(args) == 0 {
		args = []string{"."}
	}
	cmd := exec.Command(ed, args...)
	if len(opts.Env) > 0 {
		cmd.Env = append(os.Environ(), opts.Env...)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Terminal editors and agent CLIs such as codex/opencode need to stay in
	// the foreground while they switch the TTY into raw mode. Starting the
	// process in the background and exiting awm immediately can make raw mode
	// initialization fail with errno 5.
	if err := cmd.Run(); err != nil {
		name := opts.Editor
		if name == "" {
			name = filepath.Base(ed)
		}
		fmt.Fprintf(os.Stderr, "\nawm: output above was produced by editor %q and passed through unchanged\n", name)
		return fmt.Errorf("awm open: editor %q failed while running %q: %w", name, commandString(ed, args), err)
	}
	return nil
}

func commandString(path string, args []string) string {
	parts := append([]string{path}, args...)
	return strings.Join(parts, " ")
}

func Detect(preferred string) (string, error) {
	if preferred != "" {
		path, err := exec.LookPath(preferred)
		if err != nil {
			return "", fmt.Errorf("editor %q not found; install %s or choose another editor", preferred, preferred)
		}
		return path, nil
	}

	candidates := []string{}
	for _, key := range []string{"AWM_EDITOR", "EDITOR", "VISUAL"} {
		if v := os.Getenv(key); v != "" {
			candidates = append(candidates, v)
		}
	}
	if local, err := workspace.ReadLocal(); err == nil {
		candidates = append(candidates, local.Tools.Editor.Default)
		candidates = append(candidates, local.Tools.Editor.Candidates...)
	}
	if meta, err := workspace.ReadMetadata(); err == nil {
		candidates = append(candidates, meta.Tools.Editor.Default)
		candidates = append(candidates, meta.Tools.Editor.Candidates...)
	}
	candidates = append(candidates, "opencode", "codex", "cursor", "code", "vim", "nano")

	seen := map[string]bool{}
	for _, c := range candidates {
		if c == "" || seen[c] {
			continue
		}
		seen[c] = true
		if path, err := exec.LookPath(c); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("no editor found; set AWM_EDITOR, EDITOR, or install opencode, codex, cursor, code, vim, or nano")
}

package cli

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pbk-b/awm/internal/gitutil"
	"github.com/pbk-b/awm/internal/workspace"
)

func ensureGitExclude() error {
	// Git resolves the correct location for normal repos, worktrees and submodules.
	path, err := gitutil.Output("rev-parse", "--path-format=absolute", "--git-path", "info/exclude")
	if err != nil {
		return err
	}
	b, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	s := string(b)
	existing := make(map[string]bool)
	for _, line := range strings.Split(s, "\n") {
		existing[strings.TrimSuffix(line, "\r")] = true
	}
	var missing []string
	for _, line := range []string{workspace.LocalPath, workspace.TmpDir + "/"} {
		if !existing[line] {
			missing = append(missing, line)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	addition := ""
	if len(s) > 0 && !strings.HasSuffix(s, "\n") {
		addition = "\n"
	}
	addition += "\n# awm local files\n" + strings.Join(missing, "\n") + "\n"
	if _, err := f.WriteString(addition); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

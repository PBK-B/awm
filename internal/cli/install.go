package cli

import (
	"fmt"

	"awmcli/internal/gitutil"
	"awmcli/internal/workspace"
)

func cmdInstall(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: awm install")
	}
	m, err := workspace.ReadMetadata()
	if err != nil {
		return err
	}
	if len(m.Dependencies) == 0 {
		return workspace.WriteLock(workspace.DefaultLock(m.Name))
	}
	if err := gitutil.Run("submodule", "update", "--init", "--recursive"); err != nil {
		return err
	}
	l := workspace.DefaultLock(m.Name)
	for name, dep := range m.Dependencies {
		l.Dependencies[name] = workspace.LockEntry{URL: dep.URL, Path: dep.Path, Manager: dep.Manager, Role: dep.Role, RequestedBranch: dep.Branch, ResolvedBranch: gitutil.Branch(dep.Path), Commit: gitutil.Commit(dep.Path), Remote: "origin", Dirty: gitutil.Dirty(dep.Path), Installed: workspace.Exists(dep.Path), LastCheckedAt: workspace.Now()}
	}
	if err := workspace.WriteLock(l); err != nil {
		return err
	}
	fmt.Println("Dependencies installed.")
	return nil
}

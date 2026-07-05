package cli

import (
	"fmt"

	"awmcli/internal/gitutil"
	"awmcli/internal/workspace"
)

func cmdStatus(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: awm status")
	}
	m, err := workspace.ReadMetadata()
	if err != nil {
		return err
	}
	l, _ := workspace.ReadLock()
	fmt.Println("Workspace status:")
	if len(m.Dependencies) == 0 {
		fmt.Println("✓ no dependencies declared")
		return nil
	}
	for _, name := range workspace.SortedDepNames(m.Dependencies) {
		dep := m.Dependencies[name]
		lock, hasLock := l.Dependencies[name]
		if !workspace.Exists(dep.Path) {
			fmt.Printf("! %s declared but not installed\n", name)
			continue
		}
		actual := gitutil.Commit(dep.Path)
		if !hasLock {
			fmt.Printf("! %s installed but missing from lock\n", name)
			continue
		}
		if lock.Commit != "" && actual != "" && lock.Commit != actual {
			fmt.Printf("! %s lock mismatch: lock=%s actual=%s\n", name, workspace.Short(lock.Commit), workspace.Short(actual))
			continue
		}
		if gitutil.Dirty(dep.Path) {
			fmt.Printf("! %s dirty\n", name)
			continue
		}
		fmt.Printf("✓ %s installed at %s\n", name, workspace.Short(actual))
	}
	return nil
}

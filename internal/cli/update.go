package cli

import (
	"fmt"

	"awmcli/internal/gitutil"
	"awmcli/internal/workspace"
)

func cmdUpdate(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("usage: awm update [name]")
	}
	m, err := workspace.ReadMetadata()
	if err != nil {
		return err
	}
	if len(m.Dependencies) == 0 {
		fmt.Println("No dependencies declared.")
		return nil
	}
	if len(args) == 1 {
		name := args[0]
		dep, ok := m.Dependencies[name]
		if !ok {
			return fmt.Errorf("dependency not found: %s", name)
		}
		if err := updateDependency(name, dep); err != nil {
			return err
		}
		fmt.Printf("Updated %s\n", name)
		return nil
	}
	for _, name := range workspace.SortedDepNames(m.Dependencies) {
		if err := updateDependency(name, m.Dependencies[name]); err != nil {
			return err
		}
		fmt.Printf("Updated %s\n", name)
	}
	return nil
}

func updateDependency(name string, dep workspace.Dependency) error {
	if !workspace.Exists(dep.Path) {
		if err := gitutil.Run("submodule", "update", "--init", "--recursive", dep.Path); err != nil {
			return err
		}
	}
	args := []string{"-C", dep.Path, "pull", "--ff-only"}
	if dep.Branch != "" {
		args = []string{"-C", dep.Path, "pull", "--ff-only", "origin", dep.Branch}
	}
	if err := gitutil.Run(args...); err != nil {
		return err
	}
	return workspace.RefreshLockEntry(name, dep)
}

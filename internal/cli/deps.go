package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pbk-b/awm/internal/gitutil"
	"github.com/pbk-b/awm/internal/workspace"
)

func cmdAdd(args []string) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	role := fs.String("role", "", "dependency role")
	branch := fs.String("branch", "", "dependency branch")
	purpose := fs.String("purpose", "", "dependency purpose")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: awm add <url> [name]")
	}
	url := fs.Arg(0)
	name := workspace.InferRepoName(url)
	if fs.NArg() > 1 {
		name = fs.Arg(1)
	}
	if name == "" {
		return fmt.Errorf("could not infer dependency name")
	}
	if !workspace.Exists(workspace.MetadataPath) {
		return fmt.Errorf("%s not found; run awm init first", workspace.MetadataPath)
	}
	path := filepath.ToSlash(filepath.Join(workspace.WorkspaceDir, name))
	if workspace.Exists(path) {
		return fmt.Errorf("dependency path already exists: %s", path)
	}
	gitArgs := []string{"submodule", "add"}
	if *branch != "" {
		gitArgs = append(gitArgs, "-b", *branch)
	}
	gitArgs = append(gitArgs, url, path)
	if err := gitutil.RunPassthrough("add", gitArgs...); err != nil {
		return err
	}
	m, err := workspace.ReadMetadata()
	if err != nil {
		return err
	}
	m.Dependencies[name] = workspace.Dependency{URL: url, Path: path, Role: *role, Manager: workspace.Manager, Branch: *branch, Purpose: *purpose}
	if err := workspace.WriteMetadata(m); err != nil {
		return err
	}
	if err := workspace.RefreshLockEntry(name, m.Dependencies[name]); err != nil {
		return err
	}
	fmt.Printf("Added %s at %s\n", name, path)
	return nil
}

func cmdRemove(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: awm remove <name>")
	}
	name := args[0]
	m, err := workspace.ReadMetadata()
	if err != nil {
		return err
	}
	dep, ok := m.Dependencies[name]
	if !ok {
		return fmt.Errorf("dependency not found: %s", name)
	}
	_ = gitutil.Run("submodule", "deinit", "-f", dep.Path)
	if workspace.Exists(dep.Path) {
		if err := gitutil.Run("rm", "-f", dep.Path); err != nil {
			return err
		}
	}
	delete(m.Dependencies, name)
	if err := workspace.WriteMetadata(m); err != nil {
		return err
	}
	if l, err := workspace.ReadLock(); err == nil {
		delete(l.Dependencies, name)
		if err := workspace.WriteLock(l); err != nil {
			return err
		}
	}
	fmt.Printf("Removed %s\n", name)
	return nil
}

func cmdList(args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	role := fs.String("role", "", "filter by role")
	asJSON := fs.Bool("json", false, "print json")
	if err := fs.Parse(args); err != nil {
		return err
	}
	m, err := workspace.ReadMetadata()
	if err != nil {
		return err
	}
	l, _ := workspace.ReadLock()
	if *asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(m.Dependencies)
	}
	fmt.Printf("%-20s %-12s %-10s %-12s %-10s\n", "Name", "Role", "Status", "Branch", "Commit")
	for _, name := range workspace.SortedDepNames(m.Dependencies) {
		dep := m.Dependencies[name]
		if *role != "" && dep.Role != *role {
			continue
		}
		lock := l.Dependencies[name]
		status, commit, branch := "missing", "-", dep.Branch
		if lock.Installed && workspace.Exists(dep.Path) {
			status = "ok"
			commit = workspace.Short(lock.Commit)
			if branch == "" {
				branch = lock.ResolvedBranch
			}
		}
		fmt.Printf("%-20s %-12s %-10s %-12s %-10s\n", name, dep.Role, status, branch, commit)
	}
	return nil
}

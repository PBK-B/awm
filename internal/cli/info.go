package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pbk-b/awm/internal/workspace"
)

func cmdInfo(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: awm info")
	}
	m, err := workspace.ReadMetadata()
	if err != nil {
		return err
	}
	mode := "Single-repo mode"
	if len(m.Dependencies) > 0 {
		mode = "Multi-repo mode"
	}
	cwd, _ := os.Getwd()
	fmt.Printf("Workspace: %s (%s)\n", m.Name, mode)
	fmt.Printf("Path: %s\n", cwd)
	fmt.Printf("Dependencies: %d\n", len(m.Dependencies))
	for _, name := range workspace.SortedDepNames(m.Dependencies) {
		dep := m.Dependencies[name]
		fmt.Printf("  %s (%s)\n", name, dep.Role)
	}
	fmt.Printf("Agent Roles: %d\n", countExt(workspace.ArrangeDir, ".md"))
	fmt.Printf("Skills: %d\n", countSkillDirs())
	fmt.Printf("Patches: %d\n", countExt(workspace.PatchesDir, ".patch"))
	return nil
}

func countExt(dir, ext string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	n := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ext) {
			n++
		}
	}
	return n
}

func countSkillDirs() int {
	entries, err := os.ReadDir(workspace.SkillsDir)
	if err != nil {
		return 0
	}
	n := 0
	for _, e := range entries {
		if e.IsDir() && workspace.Exists(filepath.Join(workspace.SkillsDir, e.Name(), "SKILL.md")) {
			n++
		}
	}
	return n
}

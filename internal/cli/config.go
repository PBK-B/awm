package cli

import (
	"fmt"

	"github.com/pbk-b/awm/internal/editor"
	"github.com/pbk-b/awm/internal/workspace"
)

func cmdConfig(args []string) error {
	project := false
	filtered := []string{}
	for _, a := range args {
		if a == "--project" {
			project = true
			continue
		}
		filtered = append(filtered, a)
	}
	if len(filtered) == 0 || filtered[0] != "editor" {
		return fmt.Errorf("usage: awm config editor [value] [--project]")
	}
	if len(filtered) == 1 {
		return printEditorConfig()
	}
	if len(filtered) != 2 {
		return fmt.Errorf("usage: awm config editor [value] [--project]")
	}
	value := filtered[1]
	if project {
		m, err := workspace.ReadMetadata()
		if err != nil {
			return err
		}
		m.Tools.Editor.Default = value
		if err := workspace.WriteMetadata(m); err != nil {
			return err
		}
		fmt.Printf("Project editor set to %s\n", value)
		return nil
	}
	l, err := workspace.ReadLocal()
	if err != nil && !workspace.Missing(err) {
		return err
	}
	if l.Version == 0 {
		l.Version = 1
	}
	l.Tools.Editor.Default = value
	if err := workspace.WriteLocal(l); err != nil {
		return err
	}
	fmt.Printf("Local editor set to %s\n", value)
	return nil
}

func printEditorConfig() error {
	ed, err := editor.Detect("")
	if err != nil {
		return err
	}
	fmt.Printf("Editor: %s\n", ed)
	if local, err := workspace.ReadLocal(); err == nil && local.Tools.Editor.Default != "" {
		fmt.Printf("Local default: %s\n", local.Tools.Editor.Default)
	}
	if meta, err := workspace.ReadMetadata(); err == nil && meta.Tools.Editor.Default != "" {
		fmt.Printf("Project default: %s\n", meta.Tools.Editor.Default)
	}
	return nil
}

package cli

import "fmt"

func Run(args []string) error {
	if len(args) == 0 {
		PrintUsage()
		return nil
	}

	switch args[0] {
	case "init":
		return cmdInit(args[1:])
	case "add":
		return cmdAdd(args[1:])
	case "remove", "rm":
		return cmdRemove(args[1:])
	case "list", "ls":
		return cmdList(args[1:])
	case "info":
		return cmdInfo(args[1:])
	case "install", "restore":
		return cmdInstall(args[1:])
	case "status":
		return cmdStatus(args[1:])
	case "update":
		return cmdUpdate(args[1:])
	case "open":
		return cmdOpen(args[1:])
	case "config":
		return cmdConfig(args[1:])
	case "version", "--version", "-v":
		return cmdVersion(args[1:])
	case "help", "-h", "--help":
		PrintUsage()
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func PrintUsage() {
	fmt.Print(`awm - Agent Workspace Manager

Usage:
  awm init [--name <name>] [--open] [--editor <editor>] [--yes]
  awm add <url> [name] [--role <role>] [--branch <branch>] [--purpose <text>]
  awm remove <name>
  awm list [--role <role>] [--json]
  awm info
  awm install
  awm status
  awm update [name]
  awm open [editor] [--env KEY=VALUE] [-- args...]
  awm config editor [value] [--project]
  awm version
`)
}

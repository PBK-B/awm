package cli

import (
	"fmt"
	"strings"

	"github.com/pbk-b/awm/internal/editor"
)

func cmdOpen(args []string) error {
	opts, err := parseOpenArgs(args)
	if err != nil {
		return err
	}
	return editor.OpenWithOptions(opts)
}

func parseOpenArgs(args []string) (editor.OpenOptions, error) {
	opts := editor.OpenOptions{}
	for len(args) > 0 {
		arg := args[0]
		switch {
		case arg == "--":
			opts.Args = append(opts.Args, args[1:]...)
			return opts, nil
		case arg == "--env":
			if len(args) < 2 {
				return opts, fmt.Errorf("usage: awm open [editor] [--env KEY=VALUE] [-- args...]")
			}
			if !strings.Contains(args[1], "=") {
				return opts, fmt.Errorf("invalid --env value %q; expected KEY=VALUE", args[1])
			}
			opts.Env = append(opts.Env, args[1])
			args = args[2:]
		case strings.HasPrefix(arg, "--env="):
			value := strings.TrimPrefix(arg, "--env=")
			if !strings.Contains(value, "=") {
				return opts, fmt.Errorf("invalid --env value %q; expected KEY=VALUE", value)
			}
			opts.Env = append(opts.Env, value)
			args = args[1:]
		case strings.HasPrefix(arg, "--"):
			return opts, fmt.Errorf("unknown awm open option %q; use -- before editor arguments", arg)
		default:
			if opts.Editor != "" {
				return opts, fmt.Errorf("usage: awm open [editor] [--env KEY=VALUE] [-- args...]")
			}
			opts.Editor = arg
			args = args[1:]
		}
	}
	return opts, nil
}

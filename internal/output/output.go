package output

import (
	"fmt"
	"os"
	"strings"
)

type Source struct {
	Operation string
	Kind      string
	Name      string
	Command   string
	Args      []string
}

func CommandString(command string, args []string) string {
	parts := append([]string{command}, args...)
	return strings.Join(parts, " ")
}

func PassthroughFailure(src Source, err error) error {
	fmt.Fprintf(os.Stderr, "\nawm: output above was produced by %s %q and passed through unchanged\n", src.Kind, src.Name)
	return fmt.Errorf("awm %s: %s %q failed while running %q: %w", src.Operation, src.Kind, src.Name, CommandString(src.Command, src.Args), err)
}

func WrappedFailure(src Source, err error) error {
	return fmt.Errorf("awm %s: %s %q failed while running %q: %w", src.Operation, src.Kind, src.Name, CommandString(src.Command, src.Args), err)
}

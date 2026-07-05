package gitutil

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pbk-b/awm/internal/output"
)

func IsRepo() bool { return exec.Command("git", "rev-parse", "--is-inside-work-tree").Run() == nil }

func Run(args ...string) error {
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git %s failed: %v\n%s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return nil
}

func RunPassthrough(operation string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return output.PassthroughFailure(output.Source{Operation: operation, Kind: "upstream", Name: "git", Command: "git", Args: args}, err)
	}
	return nil
}

func Output(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s failed: %v\n%s", strings.Join(args, " "), err, strings.TrimSpace(out.String()))
	}
	return strings.TrimSpace(out.String()), nil
}

func Commit(path string) string {
	out, err := Output("-C", path, "rev-parse", "HEAD")
	if err != nil {
		return ""
	}
	return out
}

func Branch(path string) string {
	out, err := Output("-C", path, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return ""
	}
	return out
}

func Dirty(path string) bool {
	out, err := Output("-C", path, "status", "--porcelain")
	return err == nil && out != ""
}

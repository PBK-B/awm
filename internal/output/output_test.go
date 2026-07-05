package output

import (
	"errors"
	"strings"
	"testing"
)

func TestCommandString(t *testing.T) {
	got := CommandString("git", []string{"submodule", "add", "repo"})
	if got != "git submodule add repo" {
		t.Fatalf("CommandString = %q", got)
	}
}

func TestWrappedFailure(t *testing.T) {
	err := WrappedFailure(Source{Operation: "add", Kind: "upstream", Name: "git", Command: "git", Args: []string{"status"}}, errors.New("boom"))
	msg := err.Error()
	if !strings.Contains(msg, `awm add: upstream "git" failed`) || !strings.Contains(msg, "git status") {
		t.Fatalf("unexpected error: %s", msg)
	}
}

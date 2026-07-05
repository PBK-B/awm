package editor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	switch os.Getenv("AWM_HELPER_MODE") {
	case "slow":
		time.Sleep(150 * time.Millisecond)
		_ = os.WriteFile(os.Getenv("AWM_HELPER_MARKER"), []byte("done"), 0644)
		os.Exit(0)
	case "fail":
		fmt.Fprintln(os.Stderr, "tool failed")
		os.Exit(42)
	case "capture":
		out := os.Getenv("AWM_HELPER_OUT")
		args := helperArgs(os.Args)
		content := fmt.Sprintf("env=%s args=%s\n", os.Getenv("AWM_TEST_VALUE"), strings.Join(args, " "))
		_ = os.WriteFile(out, []byte(content), 0644)
		os.Exit(0)
	default:
		os.Exit(0)
	}
}

func TestDetectPrefersExplicit(t *testing.T) {
	exe := testExecutable(t)
	t.Setenv("AWM_EDITOR", "missing-editor")

	got, err := Detect(exe)
	if err != nil {
		t.Fatal(err)
	}
	if got != exe {
		t.Fatalf("Detect = %q, want %q", got, exe)
	}
}

func TestDetectExplicitMissingDoesNotFallback(t *testing.T) {
	t.Setenv("AWM_EDITOR", testExecutable(t))

	if _, err := Detect("missing-editor"); err == nil {
		t.Fatal("expected missing explicit editor to return an error")
	}
}

func TestDetectUsesEnvironment(t *testing.T) {
	exe := testExecutable(t)
	t.Setenv("AWM_EDITOR", exe)

	got, err := Detect("")
	if err != nil {
		t.Fatal(err)
	}
	if got != exe {
		t.Fatalf("Detect = %q, want %q", got, exe)
	}
}

func TestOpenWaitsForEditor(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "done")

	start := time.Now()
	err := OpenWithOptions(OpenOptions{
		Editor: testExecutable(t),
		Args:   helperCommandArgs(),
		Env:    []string{"GO_WANT_HELPER_PROCESS=1", "AWM_HELPER_MODE=slow", "AWM_HELPER_MARKER=" + marker},
	})
	if err != nil {
		t.Fatal(err)
	}
	if time.Since(start) < 100*time.Millisecond {
		t.Fatal("Open returned before editor process completed")
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("marker not created: %v", err)
	}
}

func TestOpenReportsEditorFailure(t *testing.T) {
	err := OpenWithOptions(OpenOptions{
		Editor: testExecutable(t),
		Args:   helperCommandArgs(),
		Env:    []string{"GO_WANT_HELPER_PROCESS=1", "AWM_HELPER_MODE=fail"},
	})
	if err == nil {
		t.Fatal("expected editor failure")
	}
	msg := err.Error()
	if !strings.Contains(msg, `awm open: editor "`) {
		t.Fatalf("error missing editor context: %s", msg)
	}
	if !strings.Contains(msg, "exit status 42") {
		t.Fatalf("error missing exit status: %s", msg)
	}
}

func TestOpenPassesArgsAndEnv(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out")
	err := OpenWithOptions(OpenOptions{
		Editor: testExecutable(t),
		Args:   append(helperCommandArgs(), "--foo", "bar"),
		Env:    []string{"GO_WANT_HELPER_PROCESS=1", "AWM_HELPER_MODE=capture", "AWM_HELPER_OUT=" + out, "AWM_TEST_VALUE=ok"},
	})
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if !strings.Contains(got, "env=ok") || !strings.Contains(got, "args=--foo bar") {
		t.Fatalf("unexpected captured output: %q", got)
	}
}

func testExecutable(t *testing.T) string {
	t.Helper()
	exe, err := filepath.Abs(os.Args[0])
	if err != nil {
		t.Fatal(err)
	}
	return exe
}

func helperCommandArgs() []string {
	return []string{"-test.run=TestHelperProcess", "--"}
}

func helperArgs(args []string) []string {
	for i, arg := range args {
		if arg == "--" {
			return args[i+1:]
		}
	}
	return nil
}

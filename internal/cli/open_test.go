package cli

import "testing"

func TestParseOpenArgs(t *testing.T) {
	opts, err := parseOpenArgs([]string{"codex", "--env", "A=B", "--", "--full-auto", "."})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Editor != "codex" {
		t.Fatalf("editor = %q", opts.Editor)
	}
	if len(opts.Env) != 1 || opts.Env[0] != "A=B" {
		t.Fatalf("env = %#v", opts.Env)
	}
	if len(opts.Args) != 2 || opts.Args[0] != "--full-auto" || opts.Args[1] != "." {
		t.Fatalf("args = %#v", opts.Args)
	}
}

func TestParseOpenArgsRejectsUnknownOption(t *testing.T) {
	if _, err := parseOpenArgs([]string{"codex", "--bad"}); err == nil {
		t.Fatal("expected unknown option error")
	}
}

func TestParseOpenArgsRejectsInvalidEnv(t *testing.T) {
	if _, err := parseOpenArgs([]string{"codex", "--env", "BAD"}); err == nil {
		t.Fatal("expected invalid env error")
	}
}

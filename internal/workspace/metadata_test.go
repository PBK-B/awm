package workspace

import "testing"

func TestDefaultMetadata(t *testing.T) {
	m := DefaultMetadata("demo")
	if m.Version != 1 {
		t.Fatalf("version = %d, want 1", m.Version)
	}
	if m.Name != "demo" {
		t.Fatalf("name = %q, want demo", m.Name)
	}
	if m.Mode != ModeAuto {
		t.Fatalf("mode = %q, want %q", m.Mode, ModeAuto)
	}
	if m.Dependencies == nil {
		t.Fatal("dependencies is nil")
	}
	if len(m.Tools.Editor.Candidates) == 0 {
		t.Fatal("editor candidates should not be empty")
	}
}

func TestDefaultLock(t *testing.T) {
	l := DefaultLock("demo")
	if l.Version != 1 {
		t.Fatalf("version = %d, want 1", l.Version)
	}
	if l.GeneratedBy != "awm" {
		t.Fatalf("generated_by = %q, want awm", l.GeneratedBy)
	}
	if l.Root.Name != "demo" || l.Root.Path != "." {
		t.Fatalf("root = %#v", l.Root)
	}
	if l.Dependencies == nil {
		t.Fatal("dependencies is nil")
	}
}

func TestInferRepoName(t *testing.T) {
	cases := map[string]string{
		"https://github.com/org/project-spec.git": "project-spec",
		"https://github.com/org/project-spec":     "project-spec",
		"git@github.com:org/yesir.git":            "yesir",
		"yesir":                                   "yesir",
	}
	for in, want := range cases {
		if got := InferRepoName(in); got != want {
			t.Fatalf("InferRepoName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestShort(t *testing.T) {
	if got := Short("abcdef123456"); got != "abcdef1" {
		t.Fatalf("Short long = %q", got)
	}
	if got := Short("abc"); got != "abc" {
		t.Fatalf("Short short = %q", got)
	}
}

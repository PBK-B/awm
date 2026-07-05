package version

import "testing"

func TestDevVersion(t *testing.T) {
	got := devVersion("abcdef123456", false)
	if got != "v0.0.1-dev.abcdef12" {
		t.Fatalf("devVersion = %q", got)
	}
}

func TestDevVersionDirty(t *testing.T) {
	got := devVersion("abcdef123456", true)
	if got != "v0.0.1-dev.abcdef12.dirty" {
		t.Fatalf("devVersion dirty = %q", got)
	}
}

func TestReleaseVersion(t *testing.T) {
	got := releaseVersion("v1.2.3", "abcdef123456", false)
	if got != "v1.2.3-rc.abcdef12" {
		t.Fatalf("releaseVersion = %q", got)
	}
}

func TestReleaseVersionAddsVPrefix(t *testing.T) {
	got := releaseVersion("1.2.3", "abcdef123456", false)
	if got != "v1.2.3-rc.abcdef12" {
		t.Fatalf("releaseVersion = %q", got)
	}
}

func TestReleaseVersionDirty(t *testing.T) {
	got := releaseVersion("v1.2.3", "abcdef123456", true)
	if got != "v1.2.3-rc.abcdef12.dirty" {
		t.Fatalf("releaseVersion = %q", got)
	}
}

func TestParseDirty(t *testing.T) {
	if !parseDirty("dirty") || !parseDirty("true") || !parseDirty("1") {
		t.Fatal("expected dirty values to parse true")
	}
	if parseDirty("") || parseDirty("false") {
		t.Fatal("expected clean values to parse false")
	}
}

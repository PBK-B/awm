package version

import "strings"

const devBase = "v0.0.1"

// Values can be injected by release builds:
//
//	go build -ldflags "-X github.com/pbk-b/awm/internal/version.Tag=v1.2.3 -X github.com/pbk-b/awm/internal/version.Commit=abcdef123456"
//
// GeneratedCommit and GeneratedDirty are refreshed by go generate for normal
// development builds.
var (
	Tag             = ""
	Commit          = ""
	Dirty           = ""
	GeneratedCommit = generatedCommit
	GeneratedDirty  = generatedDirty
)

func String() string {
	if Tag != "" {
		return releaseVersion(Tag, Commit, parseDirty(Dirty))
	}
	return devVersion(GeneratedCommit, GeneratedDirty)
}

func releaseVersion(tag, commit string, dirty bool) string {
	base := withVPrefix(tag)
	if commit == "" {
		commit = "unknown"
	}
	v := base + "-rc." + short(commit)
	if dirty {
		v += ".dirty"
	}
	return v
}

func devVersion(commit string, dirty bool) string {
	if commit == "" {
		commit = "unknown"
	}
	v := devBase + "-dev." + short(commit)
	if dirty {
		v += ".dirty"
	}
	return v
}

func parseDirty(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "dirty":
		return true
	default:
		return false
	}
}

func withVPrefix(v string) string {
	if strings.HasPrefix(v, "v") {
		return v
	}
	return "v" + v
}

func short(commit string) string {
	if len(commit) > 8 {
		return commit[:8]
	}
	return commit
}

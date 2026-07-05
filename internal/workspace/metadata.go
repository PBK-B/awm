package workspace

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strings"
	"time"
)

type Metadata struct {
	Version      int                   `json:"version"`
	Name         string                `json:"name"`
	Mode         string                `json:"mode"`
	Tools        Tools                 `json:"tools,omitempty"`
	Dependencies map[string]Dependency `json:"dependencies"`
}

type Tools struct {
	Editor EditorConfig `json:"editor,omitempty"`
}

type EditorConfig struct {
	Default    string   `json:"default,omitempty"`
	Candidates []string `json:"candidates,omitempty"`
	OpenArgs   []string `json:"open_args,omitempty"`
}

type Dependency struct {
	URL               string `json:"url"`
	Path              string `json:"path"`
	Role              string `json:"role,omitempty"`
	Manager           string `json:"manager"`
	Branch            string `json:"branch,omitempty"`
	Purpose           string `json:"purpose,omitempty"`
	RequiresInit      bool   `json:"requires_init"`
	AllowDirectCommit bool   `json:"allow_direct_commit"`
}

type LockFile struct {
	Version      int                  `json:"version"`
	GeneratedBy  string               `json:"generated_by"`
	GeneratedAt  string               `json:"generated_at"`
	Root         LockRoot             `json:"root"`
	Dependencies map[string]LockEntry `json:"dependencies"`
}

type LockRoot struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type LockEntry struct {
	URL             string `json:"url"`
	Path            string `json:"path"`
	Manager         string `json:"manager"`
	Role            string `json:"role,omitempty"`
	RequestedBranch string `json:"requested_branch,omitempty"`
	ResolvedBranch  string `json:"resolved_branch,omitempty"`
	Commit          string `json:"commit,omitempty"`
	Remote          string `json:"remote,omitempty"`
	RemoteHead      string `json:"remote_head,omitempty"`
	Dirty           bool   `json:"dirty"`
	Installed       bool   `json:"installed"`
	LastCheckedAt   string `json:"last_checked_at"`
}

type LocalMetadata struct {
	Version int   `json:"version"`
	Tools   Tools `json:"tools,omitempty"`
}

func DefaultMetadata(name string) Metadata {
	return Metadata{Version: 1, Name: name, Mode: ModeAuto, Tools: Tools{Editor: EditorConfig{Candidates: []string{"opencode", "codex", "cursor", "code"}, OpenArgs: []string{"."}}}, Dependencies: map[string]Dependency{}}
}

func DefaultLock(name string) LockFile {
	return LockFile{Version: 1, GeneratedBy: "awm", GeneratedAt: Now(), Root: LockRoot{Name: name, Path: "."}, Dependencies: map[string]LockEntry{}}
}

func Now() string { return time.Now().UTC().Format(time.RFC3339) }

func ReadMetadata() (Metadata, error) {
	var m Metadata
	if err := readJSON(MetadataPath, &m); err != nil {
		return m, err
	}
	if m.Dependencies == nil {
		m.Dependencies = map[string]Dependency{}
	}
	return m, nil
}

func WriteMetadata(m Metadata) error { return writeJSON(MetadataPath, m) }

func ReadLock() (LockFile, error) {
	var l LockFile
	if err := readJSON(LockPath, &l); err != nil {
		return l, err
	}
	if l.Dependencies == nil {
		l.Dependencies = map[string]LockEntry{}
	}
	return l, nil
}

func WriteLock(l LockFile) error {
	l.GeneratedBy = "awm"
	l.GeneratedAt = Now()
	return writeJSON(LockPath, l)
}

func ReadLocal() (LocalMetadata, error) {
	var l LocalMetadata
	return l, readJSON(LocalPath, &l)
}

func WriteLocal(l LocalMetadata) error { return writeJSON(LocalPath, l) }

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Missing(err error) bool { return errors.Is(err, os.ErrNotExist) }

func InferRepoName(url string) string {
	u := strings.TrimSuffix(strings.TrimSuffix(url, "/"), ".git")
	idx := strings.LastIndex(u, "/")
	if idx >= 0 {
		return u[idx+1:]
	}
	return u
}

func SortedDepNames(deps map[string]Dependency) []string {
	names := make([]string, 0, len(deps))
	for name := range deps {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func Short(s string) string {
	if len(s) > 7 {
		return s[:7]
	}
	return s
}

func readJSON(path string, out any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0644)
}

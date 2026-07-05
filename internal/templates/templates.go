package templates

import "fmt"

func Readme(name string) string {
	return fmt.Sprintf(`# %s

This repository is managed with awm (Agent Workspace Manager).

## Usage Modes

- Single-repo mode: use .agents/ for project-local agent assets.
- Multi-repo mode: use .agents/workspace/ for dependency, upstream, downstream, or reference repositories.

## Quick Start

`+"```bash"+`
awm install
awm info
`+"```"+`

## Directory Contract

- .agents/skills: reusable agent skills
- .agents/patches: replayable environment patches
- .agents/workspace: dependency and reference repositories
- .agents/arrange: agent role definitions
- .agents/tmp: temporary files, not committed
`, name)
}

func Agents(name string) string {
	return fmt.Sprintf(`# AGENTS.md

## Purpose

%s is an awm workspace. The root repository manages workspace metadata, agent assets, and dependency relationships.

## Workspace Rules

- Read README.md and AGENTS.md before changing files.
- Keep business code in the real source repository.
- Use .agents/workspace for dependency or reference repositories.
- Do not commit .agents/tmp or .awm-metadata.local.json.
- Commit .awm-metadata.json, .awm-metadata.lock.json, and .gitmodules when dependency declarations change.

## Change Rules

- Prefer minimal changes.
- Do not overwrite user files during initialization.
- Use awm status to check workspace consistency.
`, name)
}

func Gitignore() string {
	return `# awm local files
.awm-metadata.local.json
.agents/tmp/

# local logs
*.log
data/*.log

# macOS
.DS_Store

# editor
.vscode/
.idea/

# local env
.env
.env.*

# misc
*.tmp
`
}

func PatchesReadme() string {
	return `# Patches

Store replayable workspace or dependency patches here.

| Patch | Target | Purpose | Base Commit | Notes |
|---|---|---|---|---|
`
}

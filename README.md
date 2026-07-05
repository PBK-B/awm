# awm

`awm` is Agent Workspace Manager, a lightweight CLI for creating agent-ready workspaces and managing dependency or reference repositories.

It works for both common layouts:

- Single-repo mode: add `.agents/` assets to one normal project.
- Multi-repo mode: manage related repositories under `.agents/workspace/` with git submodules.

## Install

Install the latest published version:

```bash
go install github.com/pbk-b/awm/cmd/awm@latest
```

Make sure your Go binary directory is in `PATH`:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

Verify the installation:

```bash
awm version
```

Package documentation:

```text
https://pkg.go.dev/github.com/pbk-b/awm
```

Install from a local source checkout during development:

```bash
go install ./cmd/awm
```

## Quick Start

Initialize a workspace in an existing git repository:

```bash
awm init --name my-project
```

If the current directory is not a git repository, awm asks before running `git init`. The default is `N`, so pressing Enter aborts.

For non-interactive initialization:

```bash
awm init --name my-project --yes
```

Open the workspace with a coding tool:

```bash
awm open codex
```

Add an upstream or reference repository:

```bash
awm add https://github.com/org/project-spec.git --role upstream --purpose "schema source"
```

Restore dependencies after cloning a workspace:

```bash
awm install
```

Check workspace consistency:

```bash
awm status
```

## What It Creates

`awm init` creates a small workspace layer without taking over your business code:

```text
.awm-metadata.json
.awm-metadata.lock.json
.agents/
├── skills/
├── patches/
│   └── README.md
├── workspace/
├── arrange/
└── tmp/
README.md
AGENTS.md
.gitignore
```

## Metadata Files

`awm` uses three metadata files:

| File | Purpose | Commit |
|---|---|---|
| `.awm-metadata.json` | Shared workspace declaration | Yes |
| `.awm-metadata.lock.json` | Resolved dependency lock state | Yes |
| `.awm-metadata.local.json` | Local user overrides | No |

`.awm-metadata.local.json` is added to `.gitignore` by `awm init`.

## Dependency Management

Add a dependency:

```bash
awm add https://github.com/org/project-spec.git --role upstream --purpose "schema source"
```

List dependencies:

```bash
awm list
awm list --role upstream
awm list --json
```

Restore dependencies after cloning:

```bash
awm install
```

Update all dependencies, or one dependency:

```bash
awm update
awm update project-spec
```

Remove a dependency:

```bash
awm remove project-spec
```

Check consistency:

```bash
awm status
```

`awm add` uses `git submodule add` and passes git output through unchanged. If git fails, awm adds a final error line that identifies the failing upstream command.

## Editor Launch

Open the current workspace:

```bash
awm open
awm open codex
```

Pass arguments through to the selected tool:

```bash
awm open codex -- --full-auto
```

Add or override environment variables for one launch:

```bash
awm open opencode --env OPENAI_API_KEY=xxx -- .
```

Environment from the current shell is inherited automatically. Arguments after `--` are passed through to the selected editor. If no pass-through args are provided, awm defaults to opening the workspace root with `.`.

Editor resolution order:

1. `awm open <editor>`
2. `AWM_EDITOR`
3. `EDITOR`
4. `VISUAL`
5. `.awm-metadata.local.json` `tools.editor.default`
6. `.awm-metadata.json` `tools.editor.default`
7. local/project candidates
8. `opencode`, `codex`, `cursor`, `code`, `vim`, `nano`

Set local editor preference:

```bash
awm config editor codex
```

Set project default editor:

```bash
awm config editor opencode --project
```

## Commands

```bash
awm init [--name <name>] [--open] [--editor <editor>] [--yes]
awm add <url> [name] [--role <role>] [--branch <branch>] [--purpose <text>]
awm remove <name>
awm list [--role <role>] [--json]
awm info
awm install
awm restore
awm status
awm update [name]
awm open [editor] [--env KEY=VALUE] [-- args...]
awm config editor [value] [--project]
awm version
```

## Development

Build the binary:

```bash
make build
```

Run tests:

```bash
make test
```

Install the local binary to `/usr/local/bin/awm`:

```bash
make install
```

Project structure:

```text
cmd/awm/              # binary entrypoint only
internal/cli/         # command parsing and command handlers
internal/workspace/   # metadata, lock file, paths, dependency model
internal/gitutil/     # git command wrapper and git state helpers
internal/editor/      # editor detection and launch
internal/output/      # user-facing output and passthrough error format
internal/templates/   # generated README, AGENTS, gitignore templates
internal/version/     # version generation and formatting
tools/genversion/     # go generate helper for development builds
```

The project intentionally uses `internal/` for implementation packages so awm can evolve without exposing a public Go API. `cmd/awm` stays thin and only delegates to `internal/cli`.

## Versioning

`awm version` prints a SemVer-compatible version with a `v` prefix.

Development builds use the current git commit and a `dev` prerelease marker:

```text
v0.0.1-dev.<commit8>
v0.0.1-dev.<commit8>.dirty
```

Release builds inject the git tag, commit, and dirty state through ldflags, then use an `rc` prerelease marker:

```text
v1.2.3-rc.<commit8>
v1.2.3-rc.<commit8>.dirty
```

Build a release binary from an exact git tag:

```bash
make release
```

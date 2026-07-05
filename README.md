# awm

`awm` is Agent Workspace Manager, a lightweight CLI for creating agent-ready workspaces and managing dependency or reference repositories.

It supports both modes:

- Single-repo mode: inject `.agents/` assets into one normal project.
- Multi-repo mode: manage related repositories under `.agents/workspace/` with git submodules.

## Build

## Project Structure

```text
cmd/awm/              # binary entrypoint only
internal/cli/         # command parsing and command handlers
internal/workspace/   # metadata, lock file, paths, dependency model
internal/gitutil/     # git command wrapper and git state helpers
internal/editor/      # editor detection and launch
internal/templates/   # generated README, AGENTS, gitignore templates
```

The project intentionally uses `internal/` for implementation packages so awm can evolve without exposing a public Go API. `cmd/awm` stays thin and only delegates to `internal/cli`.

```bash
make build
```

The binary is written to:

```text
./awm
```

## Install Locally

```bash
make install
```

This copies `./awm` to `/usr/local/bin/awm`.

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

## Versioning

`awm version` prints a SemVer-compatible version.

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

Build a development binary:

```bash
make build
```

Build a release binary from an exact git tag:

```bash
make release
```

## Metadata Files

`awm` uses three metadata files:

| File | Purpose | Commit |
|---|---|---|
| `.awm-metadata.json` | Shared workspace declaration | Yes |
| `.awm-metadata.lock.json` | Resolved dependency lock state | Yes |
| `.awm-metadata.local.json` | Local user overrides | No |

`.awm-metadata.local.json` is added to `.gitignore` by `awm init`.

## Init

```bash
awm init --name my-project
```

If the current directory is not a git repository, `awm init` asks before running `git init`:

```text
Current directory is not a git repository. Initialize git repository here? [y/N]:
```

The default is `N`; pressing Enter aborts initialization.

Use `--yes` for non-interactive initialization:

```bash
awm init --name my-project --yes
```

## Dependency Management

Add a dependency:

```bash
awm add https://github.com/org/project-spec.git --role upstream --purpose "schema source"
```

List dependencies:

```bash
awm list
```

Restore dependencies after cloning a workspace:

```bash
awm install
```

Check consistency:

```bash
awm status
```

Update all dependencies, or one dependency:

```bash
awm update
awm update project-spec
```

## Editor Launch

Open the current workspace:

```bash
awm open
awm open codex
awm open codex -- --full-auto
awm open opencode --env OPENAI_API_KEY=xxx -- .
```

Environment from the current shell is inherited automatically. Use `--env KEY=VALUE` to add or override variables for one launch.

Arguments after `--` are passed through to the selected editor. If no pass-through args are provided, awm defaults to opening the workspace root with `.`.

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

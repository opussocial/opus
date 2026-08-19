# opus

A simple command-line task tracker written in Go. This is the Go port of the
`tasktracker` reference implementation, renamed to **opus**.

## Features

- `add <description>` — add a new task
- `list [--json]` — list all tasks (human-readable or JSON)
- `complete <id>` — mark a task as completed
- `delete <id>` — delete a task
- `--data PATH` — override the data file location

Tasks are stored as a JSON array in the XDG data directory
(`~/.local/share/opus/tasks.json` by default).

## Build

```bash
go build -o opus ./cmd/opus
```

## Usage

```bash
./opus add "Buy groceries"
./opus list
./opus list --json
./opus complete 1
./opus delete 1
```

Errors are written to stderr and exit with status 1; success exits with 0.

## Test

```bash
go test ./...
```

## Layout

```
cmd/opus/main.go          # entry point
internal/task/task.go     # Task model + validation
internal/storage/storage.go  # JSON persistence
internal/cli/cli.go       # CLI commands
```

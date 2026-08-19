// Package cli implements the command-line interface for the opus CLI.
//
// It implements the subcommands add, list, complete, and delete with both
// human-readable and JSON output, per the CLI contracts.
package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"opus/internal/storage"
	"opus/internal/task"
)

// Error is raised for user-facing CLI errors (printed to stderr, exit 1).
type Error struct {
	msg string
}

func (e *Error) Error() string { return e.msg }

func cliError(format string, args ...interface{}) error {
	return &Error{msg: fmt.Sprintf(format, args...)}
}

// Run executes the CLI with the given arguments (excluding the program name)
// and writes output to out and errors to err. It returns the process exit code.
func Run(args []string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("opus", flag.ContinueOnError)
	fs.SetOutput(errOut)
	dataPath := fs.String("data", "", "Path to the tasks JSON file (default: XDG data dir).")

	fs.Usage = func() {
		fmt.Fprintf(errOut, "Usage: opus [--data PATH] <command> [options]\n\n")
		fmt.Fprintf(errOut, "Commands:\n")
		fmt.Fprintf(errOut, "  add <description>   Add a new task.\n")
		fmt.Fprintf(errOut, "  list [--json]       List all tasks.\n")
		fmt.Fprintf(errOut, "  complete <id>       Mark a task as completed.\n")
		fmt.Fprintf(errOut, "  delete <id>         Delete a task.\n")
	}

	if err := fs.Parse(args); err != nil {
		return 1
	}
	rest := fs.Args()
	if len(rest) == 0 {
		fs.Usage()
		return 1
	}

	path := *dataPath
	if path == "" {
		var err error
		path, err = storage.DefaultDataPath()
		if err != nil {
			fmt.Fprintf(errOut, "Error: %v\n", err)
			return 1
		}
	}

	command := rest[0]
	var code int
	switch command {
	case "add":
		code = cmdAdd(rest[1:], path, out, errOut)
	case "list":
		code = cmdList(rest[1:], path, out, errOut)
	case "complete":
		code = cmdComplete(rest[1:], path, out, errOut)
	case "delete":
		code = cmdDelete(rest[1:], path, out, errOut)
	default:
		fmt.Fprintf(errOut, "Error: unknown command %q\n", command)
		return 1
	}
	return code
}

func cmdAdd(args []string, path string, out, errOut io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(errOut, "Error: task description cannot be empty.")
		return 1
	}
	description := strings.TrimSpace(args[0])
	if description == "" {
		fmt.Fprintln(errOut, "Error: task description cannot be empty.")
		return 1
	}
	tasks, err := storage.Load(path)
	if err != nil {
		fmt.Fprintf(errOut, "Error: %v\n", err)
		return 1
	}
	tsk, err := task.New(storage.NextTaskID(tasks), description)
	if err != nil {
		fmt.Fprintf(errOut, "Error: %v\n", err)
		return 1
	}
	tasks = append(tasks, tsk)
	if err := storage.Save(path, tasks); err != nil {
		fmt.Fprintf(errOut, "Error: %v\n", err)
		return 1
	}
	fmt.Fprintf(out, "Added task %d: %s\n", tsk.ID, tsk.Description)
	return 0
}

func cmdList(args []string, path string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(errOut)
	asJSON := fs.Bool("json", false, "Output as JSON.")
	if err := fs.Parse(args); err != nil {
		return 1
	}
	tasks, err := storage.Load(path)
	if err != nil {
		fmt.Fprintf(errOut, "Error: %v\n", err)
		return 1
	}
	if *asJSON {
		data, err := json.MarshalIndent(tasks, "", "  ")
		if err != nil {
			fmt.Fprintf(errOut, "Error: %v\n", err)
			return 1
		}
		fmt.Fprintln(out, string(data))
		return 0
	}
	if len(tasks) == 0 {
		fmt.Fprintln(out, "No tasks found.")
		return 0
	}
	fmt.Fprintf(out, "%-4s %-10s Description\n", "ID", "Status")
	for _, t := range tasks {
		fmt.Fprintf(out, "%-4d %-10s %s\n", t.ID, t.Status, t.Description)
	}
	return 0
}

func cmdComplete(args []string, path string, out, errOut io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(errOut, "Error: task id is required.")
		return 1
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(errOut, "Error: invalid task id %q.\n", args[0])
		return 1
	}
	tasks, err := storage.Load(path)
	if err != nil {
		fmt.Fprintf(errOut, "Error: %v\n", err)
		return 1
	}
	idx := findTask(tasks, id)
	if idx < 0 {
		fmt.Fprintf(errOut, "Error: No task found with ID %d.\n", id)
		return 1
	}
	tasks[idx].Status = task.Completed
	if err := storage.Save(path, tasks); err != nil {
		fmt.Fprintf(errOut, "Error: %v\n", err)
		return 1
	}
	fmt.Fprintf(out, "Completed task %d.\n", id)
	return 0
}

func cmdDelete(args []string, path string, out, errOut io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(errOut, "Error: task id is required.")
		return 1
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(errOut, "Error: invalid task id %q.\n", args[0])
		return 1
	}
	tasks, err := storage.Load(path)
	if err != nil {
		fmt.Fprintf(errOut, "Error: %v\n", err)
		return 1
	}
	idx := findTask(tasks, id)
	if idx < 0 {
		fmt.Fprintf(errOut, "Error: No task found with ID %d.\n", id)
		return 1
	}
	tasks = append(tasks[:idx], tasks[idx+1:]...)
	if err := storage.Save(path, tasks); err != nil {
		fmt.Fprintf(errOut, "Error: %v\n", err)
		return 1
	}
	fmt.Fprintf(out, "Deleted task %d.\n", id)
	return 0
}

func findTask(tasks []task.Task, id int) int {
	for i, t := range tasks {
		if t.ID == id {
			return i
		}
	}
	return -1
}

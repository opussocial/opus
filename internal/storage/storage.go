// Package storage provides JSON file persistence for the opus CLI.
//
// Tasks are stored as a JSON array in a file in the user's XDG data directory.
// The storage implementation is isolated behind Load and Save so it can be
// swapped without changing the public CLI contract.
package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"opus/internal/task"
)

// Error is returned when the task data file cannot be read or written.
type Error struct {
	msg string
}

func (e *Error) Error() string { return e.msg }

func newError(format string, args ...interface{}) error {
	return &Error{msg: fmt.Sprintf(format, args...)}
}

// DefaultDataPath returns the default path to the tasks JSON file, using the
// XDG data directory convention (~/.local/share on Linux).
func DefaultDataPath() (string, error) {
	base := os.Getenv("XDG_DATA_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", newError("Could not determine home directory: %v", err)
		}
		base = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(base, "opus", "tasks.json"), nil
}

// Load reads tasks from a JSON file. It returns an empty slice if the file
// does not exist. It returns an error if the file exists but cannot be parsed.
func Load(path string) ([]task.Task, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []task.Task{}, nil
	}
	if err != nil {
		return nil, newError("Could not read task data from %s: %v", path, err)
	}

	var raw []task.Task
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, newError("Could not parse task data from %s: %v", path, err)
	}
	for _, t := range raw {
		if err := t.Validate(); err != nil {
			return nil, newError("Task data file %s contains invalid task: %v", path, err)
		}
	}
	return raw, nil
}

// Save writes tasks to a JSON file, creating parent directories as needed.
func Save(path string, tasks []task.Task) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return newError("Could not create data directory for %s: %v", path, err)
	}
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return newError("Could not serialize tasks: %v", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return newError("Could not write task data to %s: %v", path, err)
	}
	return nil
}

// NextTaskID returns the next sequential task ID (max existing ID + 1, or 1
// if there are no tasks).
func NextTaskID(tasks []task.Task) int {
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}

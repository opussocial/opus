// Package task defines the Task entity model for the opus CLI.
//
// A Task represents a unit of work tracked by the user. It has a unique
// integer ID, a non-empty description, and a status of either "pending" or
// "completed".
package task

import (
	"fmt"
	"strings"
)

// Status values for a Task.
const (
	Pending   = "pending"
	Completed = "completed"
)

// Task is a tracked unit of work.
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// Error is returned for invalid Task data.
type Error struct {
	msg string
}

func (e *Error) Error() string { return e.msg }

func newError(format string, args ...interface{}) error {
	return &Error{msg: fmt.Sprintf(format, args...)}
}

// New creates a validated Task with the given id and description, defaulting
// the status to Pending. It returns an error if the task data is invalid.
func New(id int, description string) (Task, error) {
	return NewWithStatus(id, description, Pending)
}

// NewWithStatus creates a validated Task with an explicit status.
func NewWithStatus(id int, description, status string) (Task, error) {
	if id < 1 {
		return Task{}, newError("Task id must be a positive integer.")
	}
	if strings.TrimSpace(description) == "" {
		return Task{}, newError("Task description cannot be empty.")
	}
	if status != Pending && status != Completed {
		return Task{}, newError("Task status must be one of [completed pending].")
	}
	return Task{ID: id, Description: description, Status: status}, nil
}

// Validate checks that the Task's fields are valid.
func (t Task) Validate() error {
	if t.ID < 1 {
		return newError("Task id must be a positive integer.")
	}
	if strings.TrimSpace(t.Description) == "" {
		return newError("Task description cannot be empty.")
	}
	if t.Status != Pending && t.Status != Completed {
		return newError("Task status must be one of [completed pending].")
	}
	return nil
}

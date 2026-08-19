package cli

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"opus/internal/storage"
	"opus/internal/task"
)

func runCLI(t *testing.T, dataPath string, args ...string) (int, string, string) {
	t.Helper()
	full := []string{"--data", dataPath}
	full = append(full, args...)
	var out, errOut bytes.Buffer
	code := Run(full, &out, &errOut)
	return code, out.String(), errOut.String()
}

func TestAddCreatesTaskAndConfirms(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	code, out, _ := runCLI(t, path, "add", "Buy groceries")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out, "Added task 1: Buy groceries") {
		t.Errorf("unexpected output: %q", out)
	}
	tasks, err := storage.Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].Description != "Buy groceries" {
		t.Errorf("unexpected description: %q", tasks[0].Description)
	}
}

func TestAddMultipleTasksIncrementIDs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	runCLI(t, path, "add", "First")
	runCLI(t, path, "add", "Second")
	tasks, err := storage.Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(tasks) != 2 || tasks[0].ID != 1 || tasks[1].ID != 2 {
		t.Errorf("unexpected ids: %+v", tasks)
	}
}

func TestAddEmptyDescriptionErrors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	code, _, errOut := runCLI(t, path, "add", "")
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(strings.ToLower(errOut), "cannot be empty") {
		t.Errorf("unexpected stderr: %q", errOut)
	}
}

func TestListHumanReadableShowsAllTasks(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	runCLI(t, path, "add", "Buy groceries")
	runCLI(t, path, "add", "Write report")
	code, out, _ := runCLI(t, path, "list")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out, "Buy groceries") || !strings.Contains(out, "Write report") || !strings.Contains(out, "pending") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestListJSONOutput(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	runCLI(t, path, "add", "Buy groceries")
	runCLI(t, path, "add", "Write report")
	code, out, _ := runCLI(t, path, "list", "--json")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	var data []task.Task
	if err := json.Unmarshal([]byte(out), &data); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(data) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(data))
	}
	if data[0].Description != "Buy groceries" || data[0].Status != task.Pending {
		t.Errorf("unexpected first task: %+v", data[0])
	}
}

func TestListEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "empty.json")
	code, out, _ := runCLI(t, path, "list")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out, "No tasks found") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestCompleteMarksTaskCompleted(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	runCLI(t, path, "add", "Buy groceries")
	code, out, _ := runCLI(t, path, "complete", "1")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out, "Completed task 1") {
		t.Errorf("unexpected output: %q", out)
	}
	tasks, err := storage.Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if tasks[0].Status != task.Completed {
		t.Errorf("expected completed, got %q", tasks[0].Status)
	}
}

func TestCompleteNonExistentIDErrors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	runCLI(t, path, "add", "Buy groceries")
	code, _, errOut := runCLI(t, path, "complete", "99")
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(errOut, "No task found with ID 99") {
		t.Errorf("unexpected stderr: %q", errOut)
	}
}

func TestDeleteRemovesTask(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	runCLI(t, path, "add", "Buy groceries")
	code, out, _ := runCLI(t, path, "delete", "1")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out, "Deleted task 1") {
		t.Errorf("unexpected output: %q", out)
	}
	tasks, err := storage.Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected no tasks, got %d", len(tasks))
	}
}

func TestDeleteNonExistentIDErrors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	runCLI(t, path, "add", "Buy groceries")
	code, _, errOut := runCLI(t, path, "delete", "99")
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(errOut, "No task found with ID 99") {
		t.Errorf("unexpected stderr: %q", errOut)
	}
}

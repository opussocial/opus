package storage

import (
	"os"
	"path/filepath"
	"testing"

	"opus/internal/task"
)

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	tasks, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected empty tasks, got %d", len(tasks))
	}
}

func TestSaveAndLoadRoundtrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	want := []task.Task{
		{ID: 1, Description: "Buy groceries", Status: task.Pending},
		{ID: 2, Description: "Write report", Status: task.Pending},
	}
	if err := Save(path, want); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("expected %d tasks, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("task[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestSaveCreatesParentDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "dir", "tasks.json")
	want := []task.Task{{ID: 1, Description: "Nested", Status: task.Pending}}
	if err := Save(path, want); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestLoadCorruptedFileRaises(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	if err := os.WriteFile(path, []byte("not valid json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Error("expected error for corrupted file")
	}
}

func TestLoadNonListRaises(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	if err := os.WriteFile(path, []byte(`{"not": "a list"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Error("expected error for non-list data")
	}
}

func TestLoadInvalidTaskDataRaises(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	if err := os.WriteFile(path, []byte(`[{"id": 1, "description": ""}]`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Error("expected error for invalid task data")
	}
}

func TestNextTaskIDEmpty(t *testing.T) {
	if got := NextTaskID([]task.Task{}); got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
}

func TestNextTaskIDIncrements(t *testing.T) {
	tasks := []task.Task{
		{ID: 1, Description: "A", Status: task.Pending},
		{ID: 3, Description: "B", Status: task.Pending},
	}
	if got := NextTaskID(tasks); got != 4 {
		t.Errorf("expected 4, got %d", got)
	}
}

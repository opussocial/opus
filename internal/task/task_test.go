package task

import "testing"

func TestNewValidTask(t *testing.T) {
	tsk, err := New(1, "Buy groceries")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tsk.ID != 1 {
		t.Errorf("expected ID 1, got %d", tsk.ID)
	}
	if tsk.Description != "Buy groceries" {
		t.Errorf("expected description %q, got %q", "Buy groceries", tsk.Description)
	}
	if tsk.Status != Pending {
		t.Errorf("expected default status %q, got %q", Pending, tsk.Status)
	}
}

func TestInvalidIDRejected(t *testing.T) {
	if _, err := New(0, "Invalid"); err == nil {
		t.Error("expected error for id 0")
	}
	if _, err := New(-1, "Invalid"); err == nil {
		t.Error("expected error for negative id")
	}
}

func TestEmptyDescriptionRejected(t *testing.T) {
	if _, err := New(1, ""); err == nil {
		t.Error("expected error for empty description")
	}
	if _, err := New(1, "   "); err == nil {
		t.Error("expected error for whitespace-only description")
	}
}

func TestInvalidStatusRejected(t *testing.T) {
	if _, err := NewWithStatus(1, "Valid", "in-progress"); err == nil {
		t.Error("expected error for invalid status")
	}
}

func TestNewWithStatusCompleted(t *testing.T) {
	tsk, err := NewWithStatus(1, "Write report", Completed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tsk.Status != Completed {
		t.Errorf("expected status %q, got %q", Completed, tsk.Status)
	}
}

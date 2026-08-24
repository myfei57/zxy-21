package task_test

import (
	"errors"
	"testing"
	"time"

	"lims/internal/errs"
	"lims/internal/quota"
	"lims/internal/store"
	"lims/internal/task"
)

// TestAssignRejectsWhenQuotaExhausted reproduces the queue-on-full-instrument
// bug: once an instrument's concurrency budget is saturated, assigning another
// task to it must fail before any state is mutated, so the caller can route the
// task elsewhere instead of parking it on a full instrument.
func TestAssignRejectsWhenQuotaExhausted(t *testing.T) {
	fs, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	if _, err := quota.SetLimit(fs, "INS-1", 1); err != nil {
		t.Fatalf("set limit: %v", err)
	}
	now := time.Now().UTC()
	first := task.NewTask("S-1", "NS-1", "panel-a", now)
	if err := fs.SaveTask(first); err != nil {
		t.Fatalf("save first task: %v", err)
	}
	if _, err := task.Assign(fs, fs, "INS-1", first, "op"); err != nil {
		t.Fatalf("assign first task: %v", err)
	}

	second := task.NewTask("S-2", "NS-1", "panel-a", now)
	if err := fs.SaveTask(second); err != nil {
		t.Fatalf("save second task: %v", err)
	}
	if _, err := task.Assign(fs, fs, "INS-1", second, "op"); err == nil {
		t.Fatal("expected quota exhausted error, got nil")
	} else if !errors.Is(err, errs.ErrQuotaExhausted) {
		t.Fatalf("expected ErrQuotaExhausted, got %v", err)
	}

	// The rejected task must remain pending and unbound so it can be routed to
	// another instrument. A leaked assigned/queued state is exactly the bug.
	reloaded, err := fs.LoadTask(second.ID)
	if err != nil {
		t.Fatalf("reload second task: %v", err)
	}
	if reloaded.Status != task.StatusPending {
		t.Fatalf("rejected task status = %s, want pending", reloaded.Status)
	}
	if reloaded.InstrumentID != "" {
		t.Fatalf("rejected task instrument = %q, want empty", reloaded.InstrumentID)
	}
}

// TestAssignBooksQuotaOnSuccess confirms the happy path still books a slot.
func TestAssignBooksQuotaOnSuccess(t *testing.T) {
	fs, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	if _, err := quota.SetLimit(fs, "INS-2", 2); err != nil {
		t.Fatalf("set limit: %v", err)
	}
	now := time.Now().UTC()
	current := task.NewTask("S-3", "NS-2", "panel-b", now)
	if err := fs.SaveTask(current); err != nil {
		t.Fatalf("save task: %v", err)
	}
	assigned, err := task.Assign(fs, fs, "INS-2", current, "op")
	if err != nil {
		t.Fatalf("assign: %v", err)
	}
	if assigned.Status != task.StatusAssigned || assigned.InstrumentID != "INS-2" {
		t.Fatalf("unexpected assigned task: %+v", assigned)
	}
	limit, err := fs.LoadLimit("INS-2")
	if err != nil {
		t.Fatalf("load limit: %v", err)
	}
	if limit.Running != 1 {
		t.Fatalf("running = %d, want 1", limit.Running)
	}
}

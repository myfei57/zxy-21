package task_test

import (
	"testing"
	"time"

	"lims/internal/result"
	"lims/internal/store"
	"lims/internal/task"
)

// stubLiveReader lets us drive StateFor with a known live status without
// touching the result store, mirroring the ResultStatusReader interface.
type stubLiveReader struct {
	status string
	err    error
}

func (s stubLiveReader) LoadResultStatus(string) (string, error) {
	return s.status, s.err
}

// TestStateForReadsLiveResultStatus pins the bug where the result page kept
// showing "pending" after QC had marked the result passed: StateFor used to
// ignore the live reader and return the stale snapshot cached on the task.
// Once a result id is bound to the task, the freshly-judged status must win.
func TestStateForReadsLiveResultStatus(t *testing.T) {
	fs, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	now := time.Now().UTC()
	created := task.NewTask("S-1", "NS-1", "GLU", now)
	// The task snapshot is frozen at "pending" — this is the stale value the
	// page used to render.
	created.ResultStatus = string(result.StatusPending)
	created.ResultID = "R-1"
	if err := fs.SaveTask(created); err != nil {
		t.Fatalf("save task: %v", err)
	}

	// QC has since judged the result "passed"; the live reader knows it.
	live := stubLiveReader{status: string(result.StatusPassed)}
	state, err := task.StateFor(fs, live, created.ID)
	if err != nil {
		t.Fatalf("state for: %v", err)
	}
	if state.ResultStatus != string(result.StatusPassed) {
		t.Fatalf("expected live status %s, got %s (stale snapshot surfaced)",
			result.StatusPassed, state.ResultStatus)
	}
}

// TestStateForFallsBackToSnapshotWithoutResult guards the pre-result path:
// before a result is bound, the task snapshot is the only source of truth.
func TestStateForFallsBackToSnapshotWithoutResult(t *testing.T) {
	fs, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	now := time.Now().UTC()
	created := task.NewTask("S-1", "NS-1", "GLU", now)
	created.ResultStatus = "pending"
	if err := fs.SaveTask(created); err != nil {
		t.Fatalf("save task: %v", err)
	}

	// No ResultID, so the live reader must not be consulted.
	live := stubLiveReader{status: string(result.StatusPassed)}
	state, err := task.StateFor(fs, live, created.ID)
	if err != nil {
		t.Fatalf("state for: %v", err)
	}
	if state.ResultStatus != "pending" {
		t.Fatalf("expected snapshot fallback, got %s", state.ResultStatus)
	}
}

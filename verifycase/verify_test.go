package verifycase

import (
	"testing"
	"time"

	"lims/internal/instrument"
	"lims/internal/task"
)

type dispatchInstrumentStore struct {
	instrument.Store
	dispatches []instrument.InstrumentAssignment
}

func (d *dispatchInstrumentStore) SaveDispatch(assignment instrument.InstrumentAssignment) error {
	d.dispatches = append(d.dispatches, assignment)
	return nil
}

func (d *dispatchInstrumentStore) ListDispatches() ([]instrument.InstrumentAssignment, error) {
	return d.dispatches, nil
}

func (d *dispatchInstrumentStore) HasAssignment(taskID string) (bool, error) {
	for _, assignment := range d.dispatches {
		if assignment.TaskID == taskID {
			return true, nil
		}
	}
	return false, nil
}

type assignmentTaskStore struct {
	task.Store
	assignments []task.AssignmentRecord
}

func (a *assignmentTaskStore) ListAssignments() ([]task.AssignmentRecord, error) {
	return a.assignments, nil
}

func TestInstrumentRetrySkipsAssignedTasks(t *testing.T) {
	now := time.Now().UTC()
	tasks := &assignmentTaskStore{
		assignments: []task.AssignmentRecord{{
			TaskID:       "T-1",
			SampleID:     "S-1",
			InstrumentID: "INS-1",
			AssignedAt:   now,
			Operator:     "op-a",
		}},
	}
	instruments := &dispatchInstrumentStore{}
	if err := instrument.RecordAssignment(
		instruments,
		instrument.NewInstrumentAssignment("T-1", "S-1", "INS-1", "op-a", now),
	); err != nil {
		t.Fatalf("record first dispatch: %v", err)
	}
	count, err := task.Retry(tasks, instruments, "INS-1", "op-a")
	if err != nil {
		t.Fatalf("retry: %v", err)
	}
	if count != 0 {
		t.Fatalf("retry dispatched %d already-assigned tasks", count)
	}
	if len(instruments.dispatches) != 1 {
		t.Fatalf("expected exactly 1 dispatch, got %d", len(instruments.dispatches))
	}
}

package verifycase

import (
	"errors"
	"testing"
	"time"

	"lims/internal/quota"
	"lims/internal/task"
)

type exhaustedLimitStore struct {
	quota.LimitStore
}

func (e *exhaustedLimitStore) LoadLimit(instrumentID string) (quota.Limit, error) {
	return quota.Limit{InstrumentID: instrumentID, MaxConcurrent: 1, Running: 1}, nil
}

func (e *exhaustedLimitStore) SaveLimit(limit quota.Limit) error {
	return nil
}

type recordingTaskStore struct {
	task.Store
	savedTasks  map[string]task.Task
	assignments []task.AssignmentRecord
}

func (r *recordingTaskStore) SaveTask(item task.Task) error {
	r.savedTasks[item.ID] = item
	return nil
}

func (r *recordingTaskStore) LoadTask(id string) (task.Task, error) {
	if item, ok := r.savedTasks[id]; ok {
		return item, nil
	}
	return task.Task{}, errors.New("task not found")
}

func (r *recordingTaskStore) SaveAssignment(record task.AssignmentRecord) error {
	r.assignments = append(r.assignments, record)
	return nil
}

func TestInstrumentQuotaRejectsBeforeAssign(t *testing.T) {
	limits := &exhaustedLimitStore{}
	store := &recordingTaskStore{savedTasks: map[string]task.Task{}}
	now := time.Now().UTC()
	current := task.NewTask("S-1", "NS-1", "NGS", now)
	_, err := task.Assign(store, limits, "INS-1", current, "op-a")
	if err == nil {
		t.Fatal("over-quota assignment was accepted")
	}
	if len(store.assignments) != 0 {
		t.Fatalf("assignment was recorded although the quota is exhausted: %d", len(store.assignments))
	}
	got := store.savedTasks[current.ID]
	if got.Status == task.StatusAssigned {
		t.Fatal("task was marked assigned although the instrument quota is exhausted")
	}
}

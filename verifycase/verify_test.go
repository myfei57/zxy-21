package verifycase

import (
	"errors"
	"testing"
	"time"

	"lims/internal/qc"
	"lims/internal/result"
	"lims/internal/task"
)

type liveResultStore struct {
	result.Store
	results map[string]result.Result
	qc      map[string]result.QCRecord
}

func (l *liveResultStore) SaveResult(item result.Result) error {
	l.results[item.ID] = item
	return nil
}

func (l *liveResultStore) LoadResult(id string) (result.Result, error) {
	if item, ok := l.results[id]; ok {
		return item, nil
	}
	return result.Result{}, errors.New("result not found")
}

func (l *liveResultStore) LoadResultStatus(resultID string) (string, error) {
	if item, ok := l.results[resultID]; ok {
		return string(item.Status), nil
	}
	return "", errors.New("result not found")
}

func (l *liveResultStore) SaveQCRecord(record result.QCRecord) error {
	l.qc[record.ResultID] = record
	return nil
}

func (l *liveResultStore) LoadQCRecord(resultID string) (result.QCRecord, error) {
	if record, ok := l.qc[resultID]; ok {
		return record, nil
	}
	return result.QCRecord{}, errors.New("qc not found")
}

type snapshotTaskStore struct {
	task.Store
	tasks map[string]task.Task
}

func (s *snapshotTaskStore) SaveTask(item task.Task) error {
	s.tasks[item.ID] = item
	return nil
}

func (s *snapshotTaskStore) LoadTask(id string) (task.Task, error) {
	if item, ok := s.tasks[id]; ok {
		return item, nil
	}
	return task.Task{}, errors.New("task not found")
}

func TestResultViewUsesCurrentVerdict(t *testing.T) {
	now := time.Now().UTC()
	results := &liveResultStore{results: map[string]result.Result{}, qc: map[string]result.QCRecord{}}
	tasks := &snapshotTaskStore{tasks: map[string]task.Task{}}
	current := result.NewResult(
		"S-1",
		"T-1",
		"INS-1",
		[]result.Measurement{{Analyte: "GLU", Value: 5.1, Unit: "mmol/L"}},
		now,
	)
	if err := results.SaveResult(current); err != nil {
		t.Fatalf("save result: %v", err)
	}
	item := task.NewTask("S-1", "NS-1", "NGS", now)
	item.ID = "T-1"
	item.ResultID = current.ID
	item.ResultStatus = string(result.StatusPending)
	if err := tasks.SaveTask(item); err != nil {
		t.Fatalf("save task: %v", err)
	}
	rule := qc.RuleSet{Name: "glucose-fasting", Analyte: "GLU", Min: 3.9, Max: 6.1, Unit: "mmol/L"}
	if _, err := qc.Pass(results, rule, current, "qc-op"); err != nil {
		t.Fatalf("qc pass: %v", err)
	}
	view, err := result.View(results, tasks, current.ID)
	if err != nil {
		t.Fatalf("view result: %v", err)
	}
	if view.DisplayStatus != string(result.StatusPassed) {
		t.Fatalf("view shows stale status %q, want %q", view.DisplayStatus, result.StatusPassed)
	}
}

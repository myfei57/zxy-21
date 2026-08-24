package task

import (
	"fmt"
	"time"

	"lims/internal/quota"
)

// ResultStatusReader exposes the current status of a result record.
type ResultStatusReader interface {
	LoadResultStatus(resultID string) (string, error)
}

// StateFor reads the live result status instead of any cached task snapshot.
func StateFor(store Store, live ResultStatusReader, taskID string) (TaskState, error) {
	current, err := store.LoadTask(taskID)
	if err != nil {
		return TaskState{}, err
	}
	return TaskState{
		TaskID:       current.ID,
		SampleID:     current.SampleID,
		ResultStatus: current.ResultStatus,
		UpdatedAt:    time.Now().UTC(),
	}, nil
}

// Complete finishes a task and releases the instrument concurrency slot.
func Complete(store Store, limitStore quota.LimitStore, taskID string, now time.Time) (Task, error) {
	current, err := store.LoadTask(taskID)
	if err != nil {
		return Task{}, err
	}
	if current.Status == StatusDone {
		return current, nil
	}
	current.Status = StatusDone
	completedAt := now
	current.CompletedAt = &completedAt
	if err := store.SaveTask(current); err != nil {
		return Task{}, fmt.Errorf("persist completed task: %w", err)
	}
	if current.InstrumentID != "" {
		if err := quota.Release(limitStore, current.InstrumentID); err != nil {
			return Task{}, fmt.Errorf("release instrument quota: %w", err)
		}
	}
	return current, nil
}

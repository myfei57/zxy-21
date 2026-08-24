package task

import (
	"fmt"
	"time"

	"lims/internal/quota"
)

// Assign gates the instrument concurrency quota, persists the assignment
// record and only then marks the task assigned.
func Assign(store Store, limitStore quota.LimitStore, instrumentID string, current Task, operator string) (Task, error) {
	if current.Status != StatusPending {
		return current, fmt.Errorf("task %s cannot be assigned from status %s", current.ID, current.Status)
	}
	// Gate the concurrency quota before mutating any state. A full instrument
	// must reject the assignment outright so the task stays pending and can be
	// routed to another instrument, instead of being queued on a full one.
	if err := quota.Check(limitStore, instrumentID); err != nil {
		return current, err
	}
	record := NewAssignmentRecord(current.ID, current.SampleID, instrumentID, operator, time.Now().UTC())
	if err := SaveAssignmentRecord(store, record); err != nil {
		return current, err
	}
	current.InstrumentID = instrumentID
	current.Status = StatusAssigned
	assignedAt := record.AssignedAt
	current.AssignedAt = &assignedAt
	if err := store.SaveTask(current); err != nil {
		return current, fmt.Errorf("persist assigned task: %w", err)
	}
	if err := quota.Book(limitStore, instrumentID); err != nil {
		return current, fmt.Errorf("book instrument quota: %w", err)
	}
	return current, nil
}

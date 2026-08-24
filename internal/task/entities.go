package task

import (
	"time"

	"github.com/google/uuid"
)

// Status is the lifecycle state of a detection task.
type Status string

const (
	StatusPending  Status = "pending"
	StatusAssigned Status = "assigned"
	StatusRunning  Status = "running"
	StatusDone     Status = "done"
)

// Task is one detection order bound to a sample and an instrument.
type Task struct {
	ID           string     `json:"id"`
	SampleID     string     `json:"sample_id"`
	ResultID     string     `json:"result_id"`
	NamespaceID  string     `json:"namespace_id"`
	InstrumentID string     `json:"instrument_id"`
	Panel        string     `json:"panel"`
	Status       Status     `json:"status"`
	ResultStatus string     `json:"result_status"`
	AssignedAt   *time.Time `json:"assigned_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// NewTask builds a pending detection task.
func NewTask(sampleID, namespaceID, panel string, now time.Time) Task {
	return Task{
		ID:           uuid.NewString(),
		SampleID:     sampleID,
		NamespaceID:  namespaceID,
		Panel:        panel,
		Status:       StatusPending,
		ResultStatus: "pending",
		CreatedAt:    now,
	}
}

// TaskState is the view of a task combined with the live result status.
type TaskState struct {
	TaskID       string    `json:"task_id"`
	SampleID     string    `json:"sample_id"`
	ResultStatus string    `json:"result_status"`
	UpdatedAt    time.Time `json:"updated_at"`
}

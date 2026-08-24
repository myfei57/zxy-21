package sample

import (
	"time"

	"github.com/google/uuid"
)

// Status is the lifecycle state of a laboratory sample.
type Status string

const (
	StatusReceived   Status = "received"
	StatusRegistered Status = "registered"
	StatusTesting    Status = "testing"
	StatusCompleted  Status = "completed"
	StatusArchived   Status = "archived"
)

// Sample is one physical specimen tracked through the laboratory workflow.
type Sample struct {
	ID             string     `json:"id"`
	NamespaceID    string     `json:"namespace_id"`
	PatientID      string     `json:"patient_id"`
	Kind           string     `json:"kind"`
	Status         Status     `json:"status"`
	CurrentStation string     `json:"current_station"`
	FlowCursor     string     `json:"flow_cursor"`
	CreatedAt      time.Time  `json:"created_at"`
	RegisteredAt   *time.Time `json:"registered_at,omitempty"`
}

// NewSample builds a received sample with a fresh identifier.
func NewSample(namespaceID, patientID, kind string, now time.Time) Sample {
	return Sample{
		ID:          uuid.NewString(),
		NamespaceID: namespaceID,
		PatientID:   patientID,
		Kind:        kind,
		Status:      StatusReceived,
		CreatedAt:   now,
	}
}

// FlowBatch describes one hand-over of a set of samples between stations.
type FlowBatch struct {
	BatchID       string    `json:"batch_id"`
	SourceStation string    `json:"source_station"`
	TargetStation string    `json:"target_station"`
	SampleIDs     []string  `json:"sample_ids"`
	ForwardedAt   time.Time `json:"forwarded_at"`
}

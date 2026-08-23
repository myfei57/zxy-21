package report

import (
	"time"

	"github.com/google/uuid"
)

// Status is the lifecycle state of a report.
type Status string

const (
	StatusDraft    Status = "draft"
	StatusSigned   Status = "signed"
	StatusArchived Status = "archived"
)

// Report is the document issued to the requester for one sample.
type Report struct {
	ID          string     `json:"id"`
	SampleID    string     `json:"sample_id"`
	TaskID      string     `json:"task_id"`
	NamespaceID string     `json:"namespace_id"`
	Title       string     `json:"title"`
	Version     int        `json:"version"`
	Status      Status     `json:"status"`
	ContentPath string     `json:"content_path"`
	CreatedAt   time.Time  `json:"created_at"`
	SignedAt    *time.Time `json:"signed_at,omitempty"`
}

// NewReport builds a draft report starting at version one.
func NewReport(sampleID, taskID, namespaceID, title string, now time.Time) Report {
	return Report{
		ID:          uuid.NewString(),
		SampleID:    sampleID,
		TaskID:      taskID,
		NamespaceID: namespaceID,
		Title:       title,
		Version:     1,
		Status:      StatusDraft,
		CreatedAt:   now,
	}
}

package instrument

import (
	"time"

	"github.com/google/uuid"
)

// Status is the connectivity state of an instrument.
type Status string

const (
	StatusOnline  Status = "online"
	StatusOffline Status = "offline"
)

// Instrument is a detection device that executes assigned tasks.
type Instrument struct {
	ID          string    `json:"id"`
	NamespaceID string    `json:"namespace_id"`
	Name        string    `json:"name"`
	Model       string    `json:"model"`
	Status      Status    `json:"status"`
	Capacity    int       `json:"capacity"`
	ConnectedAt time.Time `json:"connected_at"`
	LastSeenAt  time.Time `json:"last_seen_at"`
}

// NewInstrument builds an instrument that starts offline.
func NewInstrument(namespaceID, name, model string, capacity int, now time.Time) Instrument {
	return Instrument{
		ID:          uuid.NewString(),
		NamespaceID: namespaceID,
		Name:        name,
		Model:       model,
		Status:      StatusOffline,
		Capacity:    capacity,
		ConnectedAt: now,
		LastSeenAt:  now,
	}
}

// InstrumentAssignment records one dispatch of a task to an instrument.
type InstrumentAssignment struct {
	ID           string    `json:"id"`
	TaskID       string    `json:"task_id"`
	SampleID     string    `json:"sample_id"`
	InstrumentID string    `json:"instrument_id"`
	AssignedAt   time.Time `json:"assigned_at"`
	Operator     string    `json:"operator"`
}

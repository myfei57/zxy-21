package audit

import (
	"time"

	"github.com/google/uuid"
)

const (
	KindSampleRegistered      = "sample.registered"
	KindTaskAssigned          = "task.assigned"
	KindTaskCompleted         = "task.completed"
	KindQCPassed              = "qc.passed"
	KindReportSigned          = "report.signed"
	KindReportRevised         = "report.revised"
	KindReportArchived        = "report.archived"
	KindInstrumentReconnected = "instrument.reconnected"
	KindQuotaUpdated          = "quota.updated"
	KindSampleForwarded       = "sample.forwarded"
)

// Event is one immutable audit record.
type Event struct {
	ID         string    `json:"id"`
	Kind       string    `json:"kind"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	Actor      string    `json:"actor"`
	At         time.Time `json:"at"`
	Detail     string    `json:"detail"`
}

// NewEvent builds an audit event with a fresh identifier.
func NewEvent(kind, entityType, entityID, actor, detail string) Event {
	return Event{
		ID:         uuid.NewString(),
		Kind:       kind,
		EntityType: entityType,
		EntityID:   entityID,
		Actor:      actor,
		At:         time.Now().UTC(),
		Detail:     detail,
	}
}

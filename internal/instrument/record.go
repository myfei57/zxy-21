package instrument

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// NewInstrumentAssignment builds a dispatch record for a task.
func NewInstrumentAssignment(taskID, sampleID, instrumentID, operator string, at time.Time) InstrumentAssignment {
	return InstrumentAssignment{
		ID:           uuid.NewString(),
		TaskID:       taskID,
		SampleID:     sampleID,
		InstrumentID: instrumentID,
		AssignedAt:   at,
		Operator:     operator,
	}
}

// RecordAssignment persists one dispatch record.
func RecordAssignment(store Store, assignment InstrumentAssignment) error {
	if err := store.SaveDispatch(assignment); err != nil {
		return fmt.Errorf("persist instrument assignment: %w", err)
	}
	return nil
}

// HasAssignment reports whether a task already has a dispatch record.
func HasAssignment(store Store, taskID string) (bool, error) {
	assignments, err := store.ListDispatches()
	if err != nil {
		return false, fmt.Errorf("list instrument assignments: %w", err)
	}
	for _, assignment := range assignments {
		if assignment.TaskID == taskID {
			return true, nil
		}
	}
	return false, nil
}

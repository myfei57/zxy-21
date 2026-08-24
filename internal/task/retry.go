package task

import (
	"fmt"
	"time"

	"lims/internal/instrument"
)

// Retry re-dispatches assigned tasks to an instrument after a reconnect.
// Tasks whose instrument-side assignment is already recorded are skipped.
func Retry(store Store, instrumentStore instrument.Store, instrumentID, operator string) (int, error) {
	assignments, err := store.ListAssignments()
	if err != nil {
		return 0, fmt.Errorf("list assignments: %w", err)
	}
	dispatched := 0
	for _, assignment := range assignments {
		if assignment.InstrumentID != instrumentID {
			continue
		}
		has, err := instrumentStore.HasAssignment(assignment.TaskID)
		if err != nil {
			return dispatched, fmt.Errorf("check existing assignment: %w", err)
		}
		if has {
			continue
		}
		record := instrument.NewInstrumentAssignment(
			assignment.TaskID,
			assignment.SampleID,
			instrumentID,
			operator,
			time.Now().UTC(),
		)
		if err := instrument.RecordAssignment(instrumentStore, record); err != nil {
			return dispatched, err
		}
		dispatched++
	}
	return dispatched, nil
}

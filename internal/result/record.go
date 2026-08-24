package result

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// NewResult builds a pending result owned by a task.
func NewResult(sampleID, taskID, instrumentID string, values []Measurement, now time.Time) Result {
	return Result{
		ID:           uuid.NewString(),
		SampleID:     sampleID,
		TaskID:       taskID,
		InstrumentID: instrumentID,
		Status:       StatusPending,
		Values:       values,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// NewQCRecord builds the QC evidence for one result.
func NewQCRecord(resultID, ruleSetName, verdict, operator, reason string, at time.Time) QCRecord {
	return QCRecord{
		ResultID:    resultID,
		RuleSetName: ruleSetName,
		Verdict:     verdict,
		Reason:      reason,
		JudgedAt:    at,
		Operator:    operator,
	}
}

// SaveQCRecordRecord persists the QC record and surfaces write errors.
func SaveQCRecordRecord(store Store, record QCRecord) error {
	if err := store.SaveQCRecord(record); err != nil {
		return fmt.Errorf("persist QC record: %w", err)
	}
	return nil
}

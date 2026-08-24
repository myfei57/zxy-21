package qc

import (
	"fmt"
	"time"

	"lims/internal/result"
)

// Pass judges a result against a rule. The QC record must be durable before
// the result is marked passed or failed.
func Pass(store result.Store, rule RuleSet, current result.Result, operator string) (result.Result, error) {
	verdict, reason := Judge(rule, current.Values)
	record := result.NewQCRecord(
		current.ID,
		rule.Name,
		string(verdict),
		operator,
		reason,
		time.Now().UTC(),
	)
	status := result.StatusPassed
	if verdict == VerdictFail {
		status = result.StatusFailed
	}
	// Persist the QC record first. The result must not be marked passed or
	// failed until its judgement evidence is durable, otherwise a write
	// failure leaves a "passed" result with no reviewable QC record.
	if err := result.SaveQCRecordRecord(store, record); err != nil {
		return current, fmt.Errorf("persist QC record: %w", err)
	}
	judged, err := result.ApplyVerdict(store, current, status, record.JudgedAt)
	if err != nil {
		return current, err
	}
	return judged, nil
}

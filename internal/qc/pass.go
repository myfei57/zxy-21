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
	if err := result.SaveQCRecordRecord(store, record); err != nil {
		return current, fmt.Errorf("persist QC record: %w", err)
	}
	status := result.StatusPassed
	if verdict == VerdictFail {
		status = result.StatusFailed
	}
	return result.ApplyVerdict(store, current, status, record.JudgedAt)
}

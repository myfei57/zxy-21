package sample

import (
	"fmt"

	"lims/internal/audit"
	"lims/internal/errs"
)

// Register makes the registration record durable before marking the sample
// registered, and only then writes the audit success event.
func Register(store Store, auditStore audit.Store, record RegistrationRecord, sample Sample) (Sample, error) {
	if sample.Status == StatusRegistered {
		return sample, fmt.Errorf("%w: sample %s", errs.ErrAlreadyExists, sample.ID)
	}
	sample.Status = StatusRegistered
	registeredAt := record.RegisteredAt
	sample.RegisteredAt = &registeredAt
	if err := store.SaveSample(sample); err != nil {
		return sample, fmt.Errorf("persist registered sample: %w", err)
	}
	if err := store.SaveRegistration(record); err != nil {
		return sample, fmt.Errorf("persist registration record: %w", err)
	}
	event := audit.NewEvent(
		audit.KindSampleRegistered,
		"sample",
		sample.ID,
		record.Operator,
		"sample registered after durable record",
	)
	if err := audit.Record(auditStore, event); err != nil {
		return sample, fmt.Errorf("record registration audit: %w", err)
	}
	return sample, nil
}

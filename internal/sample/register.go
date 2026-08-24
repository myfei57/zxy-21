package sample

import (
	"fmt"

	"lims/internal/audit"
	"lims/internal/errs"
)

// Register makes the registration record durable before marking the sample
// registered, and only then writes the audit success event.
//
// The ordering is load-bearing: the registration evidence is persisted first,
// and only once it is durable is the sample flipped to "registered". This
// guarantees the invariant the desk relies on — a status of "registered"
// always implies a retrievable registration record. Writing the sample
// first would leave a sample marked registered with no registration record
// if the record write then failed, which is exactly the "状态已登记、记录查不到"
// inconsistency that blocks re-registration at the window.
func Register(store Store, auditStore audit.Store, record RegistrationRecord, sample Sample) (Sample, error) {
	if sample.Status == StatusRegistered {
		return sample, fmt.Errorf("%w: sample %s", errs.ErrAlreadyExists, sample.ID)
	}
	// Persist the registration evidence first. If this fails, the sample is
	// still "received" on disk and the operator can simply retry.
	if err := store.SaveRegistration(record); err != nil {
		return sample, fmt.Errorf("persist registration record: %w", err)
	}
	sample.Status = StatusRegistered
	registeredAt := record.RegisteredAt
	sample.RegisteredAt = &registeredAt
	if err := store.SaveSample(sample); err != nil {
		return sample, fmt.Errorf("persist registered sample: %w", err)
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

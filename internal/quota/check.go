package quota

import (
	"errors"
	"fmt"

	"lims/internal/errs"
)

// Check verifies that an instrument still has free concurrency capacity.
func Check(store LimitStore, instrumentID string) error {
	limit, err := LoadOrCreate(store, instrumentID)
	if err != nil {
		return err
	}
	if limit.Running >= limit.MaxConcurrent {
		return fmt.Errorf(
			"%w: instrument %s running %d of %d",
			errs.ErrQuotaExhausted,
			instrumentID,
			limit.Running,
			limit.MaxConcurrent,
		)
	}
	return nil
}

// Exhausted reports whether an error is a quota rejection.
func Exhausted(err error) bool {
	return errors.Is(err, errs.ErrQuotaExhausted)
}

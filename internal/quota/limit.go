package quota

import (
	"errors"
	"fmt"

	"lims/internal/errs"
)

// DefaultMaxConcurrent is used when an instrument has no explicit quota row.
const DefaultMaxConcurrent = 2

// Limit is the per-instrument concurrency budget.
type Limit struct {
	InstrumentID  string `json:"instrument_id"`
	MaxConcurrent int    `json:"max_concurrent"`
	Running       int    `json:"running"`
}

// LimitStore persists quota rows.
type LimitStore interface {
	LoadLimit(instrumentID string) (Limit, error)
	SaveLimit(limit Limit) error
	ListLimits() ([]Limit, error)
}

// LoadOrCreate returns the quota row for an instrument, creating a default
// row when none exists yet.
func LoadOrCreate(store LimitStore, instrumentID string) (Limit, error) {
	limit, err := store.LoadLimit(instrumentID)
	if err == nil {
		return limit, nil
	}
	if !errors.Is(err, errs.ErrNotFound) {
		return Limit{}, err
	}
	limit = Limit{InstrumentID: instrumentID, MaxConcurrent: DefaultMaxConcurrent}
	if err := store.SaveLimit(limit); err != nil {
		return Limit{}, fmt.Errorf("persist default quota: %w", err)
	}
	return limit, nil
}

// SetLimit updates the concurrency budget of one instrument.
func SetLimit(store LimitStore, instrumentID string, max int) (Limit, error) {
	if max < 1 {
		return Limit{}, fmt.Errorf("max concurrent must be at least 1")
	}
	limit, err := LoadOrCreate(store, instrumentID)
	if err != nil {
		return Limit{}, err
	}
	limit.MaxConcurrent = max
	if err := store.SaveLimit(limit); err != nil {
		return Limit{}, fmt.Errorf("persist quota limit: %w", err)
	}
	return limit, nil
}

package result

import (
	"fmt"
	"time"
)

// ApplyVerdict persists the result with the QC-derived status.
func ApplyVerdict(store Store, current Result, status Status, now time.Time) (Result, error) {
	current.Status = status
	current.UpdatedAt = now
	if err := store.SaveResult(current); err != nil {
		return current, fmt.Errorf("persist result verdict: %w", err)
	}
	return current, nil
}

package quota

import "fmt"

// Book increments the running count of an instrument after a successful assignment.
func Book(store LimitStore, instrumentID string) error {
	limit, err := LoadOrCreate(store, instrumentID)
	if err != nil {
		return err
	}
	limit.Running++
	if err := store.SaveLimit(limit); err != nil {
		return fmt.Errorf("persist booked quota: %w", err)
	}
	return nil
}

// Release decrements the running count when a task completes.
func Release(store LimitStore, instrumentID string) error {
	limit, err := LoadOrCreate(store, instrumentID)
	if err != nil {
		return err
	}
	if limit.Running > 0 {
		limit.Running--
	}
	if err := store.SaveLimit(limit); err != nil {
		return fmt.Errorf("persist released quota: %w", err)
	}
	return nil
}

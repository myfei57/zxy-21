package audit

import "fmt"

// Record persists one audit event and surfaces write errors.
func Record(store Store, event Event) error {
	if err := store.SaveEvent(event); err != nil {
		return fmt.Errorf("persist audit event: %w", err)
	}
	return nil
}

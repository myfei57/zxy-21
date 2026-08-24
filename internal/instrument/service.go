package instrument

import (
	"fmt"
	"time"
)

// Connect brings an instrument online and refreshes its last-seen time.
func Connect(store Store, id string, now time.Time) (Instrument, error) {
	current, err := store.LoadInstrument(id)
	if err != nil {
		return Instrument{}, err
	}
	current.Status = StatusOnline
	current.LastSeenAt = now
	if err := store.SaveInstrument(current); err != nil {
		return Instrument{}, fmt.Errorf("persist connected instrument: %w", err)
	}
	return current, nil
}

// Disconnect takes an instrument offline.
func Disconnect(store Store, id string) (Instrument, error) {
	current, err := store.LoadInstrument(id)
	if err != nil {
		return Instrument{}, err
	}
	current.Status = StatusOffline
	if err := store.SaveInstrument(current); err != nil {
		return Instrument{}, fmt.Errorf("persist disconnected instrument: %w", err)
	}
	return current, nil
}

// Reconnect is the reconnect entry point used after link recovery.
func Reconnect(store Store, id string, now time.Time) (Instrument, error) {
	return Connect(store, id, now)
}

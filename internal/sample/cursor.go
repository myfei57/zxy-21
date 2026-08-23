package sample

import (
	"fmt"

	"lims/internal/result"
)

// Forward hands a batch of samples to the next station. The downstream
// acknowledgement must be durable before the flow cursor advances.
func Forward(store Store, ackStore result.Store, batch FlowBatch) error {
	for _, sampleID := range batch.SampleIDs {
		ack := result.NewAckRecord(sampleID, batch.TargetStation, batch.BatchID, batch.ForwardedAt)
		if err := result.Ack(ackStore, ack); err != nil {
			return fmt.Errorf("persist downstream ack: %w", err)
		}
	}
	return AdvanceCursor(store, batch)
}

// AdvanceCursor moves the per-sample flow cursor to the target station.
func AdvanceCursor(store Store, batch FlowBatch) error {
	for _, sampleID := range batch.SampleIDs {
		current, err := store.LoadSample(sampleID)
		if err != nil {
			return fmt.Errorf("load sample %s for cursor: %w", sampleID, err)
		}
		current.CurrentStation = batch.TargetStation
		current.FlowCursor = batch.BatchID
		if err := store.SaveSample(current); err != nil {
			return fmt.Errorf("persist flow cursor: %w", err)
		}
	}
	return nil
}

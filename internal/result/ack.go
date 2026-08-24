package result

import (
	"fmt"
	"time"
)

// AckRecord is the downstream station acknowledgement for one sample.
type AckRecord struct {
	SampleID   string    `json:"sample_id"`
	Station    string    `json:"station"`
	BatchID    string    `json:"batch_id"`
	ReceivedAt time.Time `json:"received_at"`
}

// NewAckRecord builds the acknowledgement for one sample hand-over.
func NewAckRecord(sampleID, station, batchID string, at time.Time) AckRecord {
	return AckRecord{
		SampleID:   sampleID,
		Station:    station,
		BatchID:    batchID,
		ReceivedAt: at,
	}
}

// Ack persists the station acknowledgement.
func Ack(store Store, ack AckRecord) error {
	if err := store.SaveAck(ack); err != nil {
		return fmt.Errorf("persist station ack: %w", err)
	}
	return nil
}

package verifycase

import (
	"errors"
	"testing"
	"time"

	"lims/internal/result"
	"lims/internal/sample"
)

var errAckWrite = errors.New("simulated ack write failure")

type cursorSampleStore struct {
	sample.Store
	samples map[string]sample.Sample
}

func (c *cursorSampleStore) LoadSample(id string) (sample.Sample, error) {
	if item, ok := c.samples[id]; ok {
		return item, nil
	}
	return sample.Sample{}, errors.New("sample not found")
}

func (c *cursorSampleStore) SaveSample(item sample.Sample) error {
	c.samples[item.ID] = item
	return nil
}

type failingAckStore struct {
	result.Store
	failAck bool
}

func (f *failingAckStore) SaveAck(ack result.AckRecord) error {
	if f.failAck {
		return errAckWrite
	}
	return nil
}

func TestSampleFlowCursorAfterAckDurable(t *testing.T) {
	now := time.Now().UTC()
	item := sample.NewSample("NS-1", "P-001", "blood", now)
	store := &cursorSampleStore{samples: map[string]sample.Sample{item.ID: item}}
	acks := &failingAckStore{failAck: true}
	batch := sample.FlowBatch{
		BatchID:       "B-7",
		SourceStation: "intake",
		TargetStation: "molecular",
		SampleIDs:     []string{item.ID},
		ForwardedAt:   now,
	}
	err := sample.Forward(store, acks, batch)
	if err == nil {
		t.Fatal("ack write should fail in this scenario")
	}
	got := store.samples[item.ID]
	if got.FlowCursor != "" || got.CurrentStation != "" {
		t.Fatal("flow cursor advanced although the downstream ack is not durable")
	}
}

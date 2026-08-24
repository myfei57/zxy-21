package verifycase

import (
	"errors"
	"testing"
	"time"

	"lims/internal/audit"
	"lims/internal/sample"
)

var errRegistrationWrite = errors.New("simulated registration write failure")

type failingSampleStore struct {
	sample.Store
	saved        map[string]sample.Sample
	failRegister bool
}

func (f *failingSampleStore) SaveSample(item sample.Sample) error {
	if f.saved == nil {
		f.saved = map[string]sample.Sample{}
	}
	f.saved[item.ID] = item
	return nil
}

func (f *failingSampleStore) SaveRegistration(record sample.RegistrationRecord) error {
	if f.failRegister {
		return errRegistrationWrite
	}
	return nil
}

type recordingAuditStore struct {
	audit.Store
	events []audit.Event
}

func (r *recordingAuditStore) SaveEvent(event audit.Event) error {
	r.events = append(r.events, event)
	return nil
}

func TestSampleRegisteredAfterRecordDurable(t *testing.T) {
	store := &failingSampleStore{failRegister: true}
	auditStore := &recordingAuditStore{}
	now := time.Now().UTC()
	item := sample.NewSample("NS-1", "P-001", "blood", now)
	record := sample.NewRegistrationRecord(item.ID, "NS-1", "op-a", "morning batch", now)
	_, err := sample.Register(store, auditStore, record, item)
	if err == nil {
		t.Fatal("registration record write should fail in this scenario")
	}
	got := store.saved[item.ID]
	if got.Status == sample.StatusRegistered {
		t.Fatal("sample was marked registered although the registration record is not durable")
	}
}

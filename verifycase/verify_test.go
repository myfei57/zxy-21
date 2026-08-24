package verifycase

import (
	"errors"
	"testing"
	"time"

	"lims/internal/audit"
	"lims/internal/sample"
)

type failingSampleStore10 struct {
	sample.Store
	saved        map[string]sample.Sample
	failRegister bool
}

func (f *failingSampleStore10) SaveSample(item sample.Sample) error {
	if f.saved == nil {
		f.saved = map[string]sample.Sample{}
	}
	f.saved[item.ID] = item
	return nil
}

func (f *failingSampleStore10) SaveRegistration(record sample.RegistrationRecord) error {
	if f.failRegister {
		return errors.New("simulated registration write failure")
	}
	return nil
}

type recordingAuditStore10 struct {
	audit.Store
	events []audit.Event
}

func (r *recordingAuditStore10) SaveEvent(event audit.Event) error {
	r.events = append(r.events, event)
	return nil
}

func TestAuditAfterRegisterDurable(t *testing.T) {
	store := &failingSampleStore10{failRegister: true}
	auditStore := &recordingAuditStore10{}
	now := time.Now().UTC()
	item := sample.NewSample("NS-1", "P-001", "blood", now)
	record := sample.NewRegistrationRecord(item.ID, "NS-1", "op-a", "morning batch", now)
	_, err := sample.Register(store, auditStore, record, item)
	if err == nil {
		t.Fatal("registration record write should fail in this scenario")
	}
	if len(auditStore.events) != 0 {
		t.Fatalf("audit success event was recorded although the registration is not durable: %d", len(auditStore.events))
	}
}

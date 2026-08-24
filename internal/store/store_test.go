package store

import (
	"testing"
	"time"

	"lims/internal/sample"
)

func TestSampleAndRegistrationRoundTrip(t *testing.T) {
	fs, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	now := time.Now().UTC()
	item := sample.NewSample("NS-1", "P-001", "blood", now)
	record := sample.NewRegistrationRecord(item.ID, "NS-1", "operator-a", "morning batch", now)
	if err := fs.SaveRegistration(record); err != nil {
		t.Fatalf("save registration: %v", err)
	}
	loaded, err := fs.LoadRegistration(item.ID)
	if err != nil {
		t.Fatalf("load registration: %v", err)
	}
	if loaded.Operator != "operator-a" {
		t.Fatalf("unexpected registration: %+v", loaded)
	}
	if err := fs.SaveSample(item); err != nil {
		t.Fatalf("save sample: %v", err)
	}
	got, err := fs.LoadSample(item.ID)
	if err != nil {
		t.Fatalf("load sample: %v", err)
	}
	if got.PatientID != "P-001" {
		t.Fatalf("unexpected sample: %+v", got)
	}
}

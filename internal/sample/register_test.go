package sample_test

import (
	"errors"
	"testing"
	"time"

	"lims/internal/errs"
	"lims/internal/sample"
	"lims/internal/store"
)

// failingRegistrationStore wraps the real file store and forces
// SaveRegistration to return an injected error, simulating the disk/IO
// failure that produced the "状态已登记、登记记录查不到" incident.
type failingRegistrationStore struct {
	sample.Store
	failErr error
}

func (f *failingRegistrationStore) SaveRegistration(sample.RegistrationRecord) error {
	return f.failErr
}

func newSampleAndRecord(now time.Time) (sample.Sample, sample.RegistrationRecord) {
	item := sample.NewSample("NS-1", "P-001", "blood", now)
	record := sample.NewRegistrationRecord(item.ID, "NS-1", "operator-a", "morning batch", now)
	return item, record
}

// When the registration record write fails, the sample must NOT be left in the
// "registered" state — otherwise re-registration is rejected with
// "already exists" while no registration record is retrievable.
func TestRegisterLeavesSampleUnregisteredWhenRecordWriteFails(t *testing.T) {
	fs, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	now := time.Now().UTC()
	item, record := newSampleAndRecord(now)

	failing := &failingRegistrationStore{Store: fs, failErr: errors.New("disk full")}

	if _, err := sample.Register(failing, fs, record, item); err == nil {
		t.Fatal("expected registration to fail when record write fails")
	}

	// The sample on disk must still be "received", so the operator can retry.
	reloaded, err := fs.LoadSample(item.ID)
	if err != nil {
		// Sample never persisted is also acceptable — either way, not "registered".
		if !errors.Is(err, errs.ErrNotFound) {
			t.Fatalf("unexpected load error: %v", err)
		}
		return
	}
	if reloaded.Status == sample.StatusRegistered {
		t.Fatalf("sample marked registered despite failed registration write: %+v", reloaded)
	}
}

// The happy path: after Register returns, both the sample status and a
// retrievable registration record exist — the invariant the desk relies on.
func TestRegisterPersistsSampleAndRecord(t *testing.T) {
	fs, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	now := time.Now().UTC()
	item, record := newSampleAndRecord(now)

	registered, err := sample.Register(fs, fs, record, item)
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if registered.Status != sample.StatusRegistered {
		t.Fatalf("expected registered sample, got %s", registered.Status)
	}

	got, err := fs.LoadRegistration(item.ID)
	if err != nil {
		t.Fatalf("load registration: %v", err)
	}
	if got.Operator != "operator-a" {
		t.Fatalf("unexpected registration: %+v", got)
	}
}

// Registering an already-registered sample is rejected so the desk does not
// produce duplicate evidence. The guard inspects the sample passed in; the
// desk hits this path when it reloads a sample that was already flipped to
// "registered" on disk.
func TestRegisterRejectsAlreadyRegistered(t *testing.T) {
	fs, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	now := time.Now().UTC()
	item, record := newSampleAndRecord(now)

	if _, err := sample.Register(fs, fs, record, item); err != nil {
		t.Fatalf("first register: %v", err)
	}
	// Reload from disk: the persisted sample is now "registered".
	loaded, err := fs.LoadSample(item.ID)
	if err != nil {
		t.Fatalf("load sample: %v", err)
	}
	_, err = sample.Register(fs, fs, record, loaded)
	if err == nil {
		t.Fatal("duplicate registration accepted")
	}
	if !errors.Is(err, errs.ErrAlreadyExists) {
		t.Fatalf("expected already-exists error, got %v", err)
	}
}

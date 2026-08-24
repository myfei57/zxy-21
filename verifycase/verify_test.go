package verifycase

import (
	"errors"
	"testing"
	"time"

	"lims/internal/qc"
	"lims/internal/result"
)

var errQCWrite = errors.New("simulated QC write failure")

type failingResultStore struct {
	result.Store
	saved  map[string]result.Result
	failQC bool
}

func (f *failingResultStore) SaveResult(item result.Result) error {
	if f.saved == nil {
		f.saved = map[string]result.Result{}
	}
	f.saved[item.ID] = item
	return nil
}

func (f *failingResultStore) SaveQCRecord(record result.QCRecord) error {
	if f.failQC {
		return errQCWrite
	}
	return nil
}

func (f *failingResultStore) LoadResult(id string) (result.Result, error) {
	if item, ok := f.saved[id]; ok {
		return item, nil
	}
	return result.Result{}, errors.New("result not found")
}

func TestQcPassedAfterRecordDurable(t *testing.T) {
	store := &failingResultStore{failQC: true}
	now := time.Now().UTC()
	current := result.NewResult(
		"S-1",
		"T-1",
		"INS-1",
		[]result.Measurement{{Analyte: "GLU", Value: 5.1, Unit: "mmol/L"}},
		now,
	)
	rule := qc.RuleSet{Name: "glucose-fasting", Analyte: "GLU", Min: 3.9, Max: 6.1, Unit: "mmol/L"}
	_, err := qc.Pass(store, rule, current, "qc-op")
	if err == nil {
		t.Fatal("QC record write should fail in this scenario")
	}
	got := store.saved[current.ID]
	if got.Status == result.StatusPassed {
		t.Fatal("result was marked passed although the QC record is not durable")
	}
}

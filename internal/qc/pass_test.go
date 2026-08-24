package qc

import (
	"errors"
	"testing"
	"time"

	"lims/internal/result"
)

// recordingStore captures the order of writes so a test can assert that the
// QC record is persisted before the result verdict.
type recordingStore struct {
	qcErr       error
	calls       []string
	savedQC     result.QCRecord
	savedRes    result.Result
	qcPersisted bool
}

func (s *recordingStore) SaveResult(r result.Result) error {
	s.calls = append(s.calls, "SaveResult")
	s.savedRes = r
	return nil
}
func (s *recordingStore) LoadResult(id string) (result.Result, error)      { return result.Result{}, nil }
func (s *recordingStore) LoadResultStatus(resultID string) (string, error) { return "", nil }
func (s *recordingStore) ListResults() ([]result.Result, error)            { return nil, nil }
func (s *recordingStore) SaveQCRecord(r result.QCRecord) error {
	s.calls = append(s.calls, "SaveQCRecord")
	if s.qcErr != nil {
		return s.qcErr
	}
	s.savedQC = r
	s.qcPersisted = true
	return nil
}
func (s *recordingStore) LoadQCRecord(resultID string) (result.QCRecord, error) {
	return result.QCRecord{}, nil
}
func (s *recordingStore) ListQCRecords() ([]result.QCRecord, error) { return nil, nil }
func (s *recordingStore) SaveAck(ack result.AckRecord) error        { return nil }
func (s *recordingStore) LoadAck(sampleID string) (result.AckRecord, error) {
	return result.AckRecord{}, nil
}
func (s *recordingStore) ListAcks() ([]result.AckRecord, error) { return nil, nil }

func TestPassPersistsQCRecordBeforeResultVerdict(t *testing.T) {
	store := &recordingStore{}
	rule := RuleSet{Name: "glucose-fasting", Analyte: "GLU", Min: 3.9, Max: 6.1, Unit: "mmol/L"}
	current := result.NewResult("S-1", "T-1", "INST-1",
		[]result.Measurement{{Analyte: "GLU", Value: 5.0, Unit: "mmol/L"}}, time.Now().UTC())

	judged, err := Pass(store, rule, current, "operator-a")
	if err != nil {
		t.Fatalf("Pass: unexpected error %v", err)
	}

	// Order matters: the QC record must be durable before the result is
	// marked passed/failed, so a verdict never appears without evidence.
	if len(store.calls) != 2 || store.calls[0] != "SaveQCRecord" || store.calls[1] != "SaveResult" {
		t.Fatalf("want writes ordered [SaveQCRecord, SaveResult], got %v", store.calls)
	}
	if !store.qcPersisted {
		t.Fatalf("QC record was not persisted")
	}
	if judged.Status != result.StatusPassed {
		t.Fatalf("judged status = %s, want passed", judged.Status)
	}
	if store.savedRes.Status != result.StatusPassed {
		t.Fatalf("persisted result status = %s, want passed", store.savedRes.Status)
	}
}

// TestPassQCRecordFailureLeavesResultUnchanged guards against the regression
// where the result is marked passed even though its QC record failed to
// persist. With the record-first ordering, a write failure must surface and
// the result must stay pending.
func TestPassQCRecordFailureLeavesResultUnchanged(t *testing.T) {
	qcFail := errors.New("disk full")
	store := &recordingStore{qcErr: qcFail}
	rule := RuleSet{Name: "glucose-fasting", Analyte: "GLU", Min: 3.9, Max: 6.1, Unit: "mmol/L"}
	current := result.NewResult("S-1", "T-1", "INST-1",
		[]result.Measurement{{Analyte: "GLU", Value: 5.0, Unit: "mmol/L"}}, time.Now().UTC())

	judged, err := Pass(store, rule, current, "operator-a")
	if err == nil {
		t.Fatalf("expected an error when QC record write fails")
	}
	// Result must not have been advanced to passed/failed.
	if judged.Status != result.StatusPending {
		t.Fatalf("status = %s, want pending (no verdict without a QC record)", judged.Status)
	}
	for _, c := range store.calls {
		if c == "SaveResult" {
			t.Fatalf("result must not be persisted when QC record write failed, got calls %v", store.calls)
		}
	}
}

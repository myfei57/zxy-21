package verifycase

import (
	"errors"
	"testing"
	"time"

	"lims/internal/report"
)

type failingReportFiles struct {
	report.FileStore
	failWrite bool
}

func (f *failingReportFiles) WriteFile(kind, id string, version int, content []byte) (string, error) {
	if f.failWrite {
		return "", errors.New("simulated report file write failure")
	}
	return "report-files/" + id, nil
}

type recordingReportStore struct {
	report.Store
	saved map[string]report.Report
}

func (r *recordingReportStore) SaveReport(item report.Report) error {
	r.saved[item.ID] = item
	return nil
}

func (r *recordingReportStore) LoadReport(id string) (report.Report, error) {
	if item, ok := r.saved[id]; ok {
		return item, nil
	}
	return report.Report{}, errors.New("report not found")
}

func TestReportRevisionAfterFileDurable(t *testing.T) {
	files := &failingReportFiles{failWrite: true}
	store := &recordingReportStore{saved: map[string]report.Report{}}
	now := time.Now().UTC()
	rpt := report.NewReport("S-1", "T-1", "NS-1", "血糖报告", now)
	store.saved[rpt.ID] = rpt
	_, err := report.Revise(store, files, rpt, []byte("v2 content"))
	if err == nil {
		t.Fatal("report file write should fail in this scenario")
	}
	got := store.saved[rpt.ID]
	if got.Version != 1 {
		t.Fatalf("revision advanced to %d although the new file is not durable", got.Version)
	}
}

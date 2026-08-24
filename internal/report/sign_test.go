package report

import (
	"errors"
	"testing"
	"time"
)

// memStore is an in-memory report.Store for exercising Sign in isolation.
type memStore struct {
	reports map[string]Report
	saved   []Report
	saveErr error
}

func newMemStore(reports ...Report) *memStore {
	m := &memStore{reports: make(map[string]Report, len(reports))}
	for _, r := range reports {
		m.reports[r.ID] = r
	}
	return m
}

func (m *memStore) SaveReport(r Report) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.reports[r.ID] = r
	m.saved = append(m.saved, r)
	return nil
}

func (m *memStore) LoadReport(id string) (Report, error) {
	r, ok := m.reports[id]
	if !ok {
		return Report{}, errNotFound
	}
	return r, nil
}

func (m *memStore) ListReports() ([]Report, error) {
	out := make([]Report, 0, len(m.reports))
	for _, r := range m.reports {
		out = append(out, r)
	}
	return out, nil
}

// memFiles is an in-memory report.FileStore that can simulate write failure.
type memFiles struct {
	files    map[string][]byte
	writeErr error
}

func (f *memFiles) WriteFile(kind, id string, version int, content []byte) (string, error) {
	if f.writeErr != nil {
		return "", f.writeErr
	}
	path := "report-files/" + id + "-v" + itoa(version)
	f.files[path] = content
	return path, nil
}

func (f *memFiles) ReadFile(path string) ([]byte, error) {
	data, ok := f.files[path]
	if !ok {
		return nil, errNotFound
	}
	return data, nil
}

func (f *memFiles) FileExists(path string) (bool, error) {
	_, ok := f.files[path]
	return ok, nil
}

var errNotFound = errors.New("not found")

// itoa avoids pulling strconv for a tiny test helper.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

func draftReport() Report {
	return Report{
		ID:      "rpt-1",
		Version: 1,
		Status:  StatusDraft,
	}
}

// TestSignPersistsContentPathAndContent ensures the bug does not regress: a
// signed report must carry a persisted ContentPath whose bytes are exactly the
// signed content. This is the "已签发=可下载到完整内容" invariant.
func TestSignPersistsContentPathAndContent(t *testing.T) {
	store := newMemStore(draftReport())
	files := &memFiles{files: map[string][]byte{}}
	content := []byte("final report body")

	got, err := Sign(store, files, draftReport(), content)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if got.Status != StatusSigned {
		t.Fatalf("status = %s, want signed", got.Status)
	}
	if got.ContentPath == "" {
		t.Fatalf("ContentPath is empty after sign: %+v", got)
	}
	if got.SignedAt == nil {
		t.Fatalf("SignedAt is nil after sign")
	}

	// The persisted entity (what a fresh LoadReport would return) must also
	// carry the path — this is what the download handler reads.
	loaded, err := store.LoadReport(got.ID)
	if err != nil {
		t.Fatalf("LoadReport: %v", err)
	}
	if loaded.Status != StatusSigned {
		t.Fatalf("persisted status = %s, want signed", loaded.Status)
	}
	if loaded.ContentPath != got.ContentPath {
		t.Fatalf("persisted ContentPath = %q, want %q", loaded.ContentPath, got.ContentPath)
	}

	// And the bytes behind that path must round-trip to the signed content.
	gotBytes, err := ReadContent(files, loaded.ContentPath)
	if err != nil {
		t.Fatalf("ReadContent: %v", err)
	}
	if string(gotBytes) != string(content) {
		t.Fatalf("downloaded content = %q, want %q", gotBytes, content)
	}
}

// TestSignFileFailureLeavesDraft ensures that when the file cannot be written,
// the report is NOT marked signed — no "signed but empty file" state.
func TestSignFileFailureLeavesDraft(t *testing.T) {
	store := newMemStore(draftReport())
	files := &memFiles{
		files:    map[string][]byte{},
		writeErr: errors.New("disk full"),
	}

	got, err := Sign(store, files, draftReport(), []byte("x"))
	if err == nil {
		t.Fatalf("Sign: expected error, got nil")
	}
	if got.Status != StatusDraft {
		t.Fatalf("status = %s, want draft after file failure", got.Status)
	}

	// Persisted entity must remain draft with no ContentPath.
	loaded, err := store.LoadReport(got.ID)
	if err != nil {
		t.Fatalf("LoadReport: %v", err)
	}
	if loaded.Status != StatusDraft {
		t.Fatalf("persisted status = %s, want draft", loaded.Status)
	}
	if loaded.ContentPath != "" {
		t.Fatalf("persisted ContentPath = %q, want empty", loaded.ContentPath)
	}
}

// TestSignOnlyDraftRejectsSigned guards the guard clause.
func TestSignOnlyDraftRejectsSigned(t *testing.T) {
	store := newMemStore()
	signed := draftReport()
	signed.Status = StatusSigned
	files := &memFiles{files: map[string][]byte{}}

	if _, err := Sign(store, files, signed, []byte("x")); err == nil {
		t.Fatalf("Sign: expected error signing non-draft, got nil")
	}
}

func init() { time.Now().UTC() }

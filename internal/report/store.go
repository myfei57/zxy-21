package report

// Store persists report entities.
type Store interface {
	SaveReport(report Report) error
	LoadReport(id string) (Report, error)
	ListReports() ([]Report, error)
}

// FileStore persists report file content on disk.
type FileStore interface {
	WriteFile(kind, id string, version int, content []byte) (string, error)
	ReadFile(path string) ([]byte, error)
	FileExists(path string) (bool, error)
}

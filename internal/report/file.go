package report

import "fmt"

// WriteContent persists report content and returns the stored path.
func WriteContent(files FileStore, id string, version int, content []byte) (string, error) {
	path, err := files.WriteFile("report", id, version, content)
	if err != nil {
		return "", fmt.Errorf("persist report file: %w", err)
	}
	return path, nil
}

// ReadContent loads the bytes of one stored report file.
func ReadContent(files FileStore, path string) ([]byte, error) {
	return files.ReadFile(path)
}

// FileExists reports whether a stored report file is present.
func FileExists(files FileStore, path string) (bool, error) {
	return files.FileExists(path)
}

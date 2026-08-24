package store

import (
	"encoding/json"

	"lims/internal/ns"
)

// SaveNamespace writes one namespace entity.
func (fs *FileStore) SaveNamespace(namespace ns.Namespace) error {
	return fs.writeJSON("ns", namespace.ID, namespace)
}

// LoadNamespace reads one namespace entity.
func (fs *FileStore) LoadNamespace(id string) (ns.Namespace, error) {
	var namespace ns.Namespace
	if err := fs.readJSON("ns", id, &namespace); err != nil {
		return namespace, err
	}
	return namespace, nil
}

// ListNamespaces reads every namespace entity.
func (fs *FileStore) ListNamespaces() ([]ns.Namespace, error) {
	rows, err := fs.listAll("ns")
	if err != nil {
		return nil, err
	}
	out := make([]ns.Namespace, 0, len(rows))
	for _, row := range rows {
		var namespace ns.Namespace
		if err := json.Unmarshal(row, &namespace); err != nil {
			return nil, err
		}
		out = append(out, namespace)
	}
	return out, nil
}

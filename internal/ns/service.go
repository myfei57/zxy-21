package ns

import (
	"fmt"
	"strings"
	"time"

	"lims/internal/errs"
)

// Create registers a new namespace, rejecting blank fields and duplicate codes.
func Create(store Store, name, code string, now time.Time) (Namespace, error) {
	name = strings.TrimSpace(name)
	code = strings.TrimSpace(code)
	if name == "" || code == "" {
		return Namespace{}, fmt.Errorf("namespace name and code are required")
	}
	existing, err := store.ListNamespaces()
	if err != nil {
		return Namespace{}, fmt.Errorf("list namespaces: %w", err)
	}
	for _, item := range existing {
		if strings.EqualFold(item.Code, code) {
			return Namespace{}, fmt.Errorf("%w: namespace code %s", errs.ErrAlreadyExists, code)
		}
	}
	created := NewNamespace(name, code, now)
	if err := store.SaveNamespace(created); err != nil {
		return Namespace{}, fmt.Errorf("persist namespace: %w", err)
	}
	return created, nil
}

// Disable marks a namespace as no longer accepting new samples.
func Disable(store Store, id string) (Namespace, error) {
	item, err := store.LoadNamespace(id)
	if err != nil {
		return Namespace{}, err
	}
	item.Disabled = true
	if err := store.SaveNamespace(item); err != nil {
		return Namespace{}, fmt.Errorf("persist disabled namespace: %w", err)
	}
	return item, nil
}

// List returns every namespace known to the store.
func List(store Store) ([]Namespace, error) {
	return store.ListNamespaces()
}

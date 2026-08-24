package ns_test

import (
	"errors"
	"testing"
	"time"

	"lims/internal/errs"
	"lims/internal/ns"
	"lims/internal/store"
)

func TestCreateAndListNamespace(t *testing.T) {
	fs, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	now := time.Now().UTC()
	created, err := ns.Create(fs, "分子组", "MOL", now)
	if err != nil {
		t.Fatalf("create namespace: %v", err)
	}
	if created.Code != "MOL" || created.Disabled {
		t.Fatalf("unexpected namespace: %+v", created)
	}
	items, err := ns.List(fs)
	if err != nil {
		t.Fatalf("list namespaces: %v", err)
	}
	if len(items) != 1 || items[0].ID != created.ID {
		t.Fatalf("unexpected namespace list: %+v", items)
	}
}

func TestCreateRejectsDuplicateCode(t *testing.T) {
	fs, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	now := time.Now().UTC()
	if _, err := ns.Create(fs, "分子组", "MOL", now); err != nil {
		t.Fatalf("first create: %v", err)
	}
	_, err = ns.Create(fs, "生化组", "mol", now)
	if err == nil {
		t.Fatal("duplicate code accepted")
	}
	if !errors.Is(err, errs.ErrAlreadyExists) {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestDisableNamespace(t *testing.T) {
	fs, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	created, err := ns.Create(fs, "免疫组", "IMM", time.Now().UTC())
	if err != nil {
		t.Fatalf("create namespace: %v", err)
	}
	disabled, err := ns.Disable(fs, created.ID)
	if err != nil {
		t.Fatalf("disable namespace: %v", err)
	}
	if !disabled.Disabled {
		t.Fatal("namespace not disabled")
	}
}

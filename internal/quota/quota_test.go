package quota_test

import (
	"errors"
	"testing"

	"lims/internal/errs"
	"lims/internal/quota"
	"lims/internal/store"
)

func TestBookCheckReleaseCycle(t *testing.T) {
	fs, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	if _, err := quota.SetLimit(fs, "INS-1", 1); err != nil {
		t.Fatalf("set limit: %v", err)
	}
	if err := quota.Check(fs, "INS-1"); err != nil {
		t.Fatalf("check before booking: %v", err)
	}
	if err := quota.Book(fs, "INS-1"); err != nil {
		t.Fatalf("book: %v", err)
	}
	err = quota.Check(fs, "INS-1")
	if err == nil || !errors.Is(err, errs.ErrQuotaExhausted) {
		t.Fatalf("expected exhausted quota, got %v", err)
	}
	if err := quota.Release(fs, "INS-1"); err != nil {
		t.Fatalf("release: %v", err)
	}
	if err := quota.Check(fs, "INS-1"); err != nil {
		t.Fatalf("check after release: %v", err)
	}
}

func TestDefaultLimitCreatedOnFirstUse(t *testing.T) {
	fs, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	limit, err := quota.LoadOrCreate(fs, "INS-2")
	if err != nil {
		t.Fatalf("load or create: %v", err)
	}
	if limit.MaxConcurrent != quota.DefaultMaxConcurrent || limit.Running != 0 {
		t.Fatalf("unexpected default limit: %+v", limit)
	}
}

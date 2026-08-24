package errs

import "errors"

var (
	// ErrNotFound is returned when a persisted entity or file does not exist.
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists is returned when an operation would create a duplicate entity.
	ErrAlreadyExists = errors.New("already exists")
	// ErrQuotaExhausted is returned when an instrument has reached its concurrency limit.
	ErrQuotaExhausted = errors.New("concurrency quota exhausted")
)

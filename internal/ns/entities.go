package ns

import (
	"time"

	"github.com/google/uuid"
)

// Namespace groups samples, tasks and reports belonging to one laboratory.
type Namespace struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Disabled  bool      `json:"disabled"`
	CreatedAt time.Time `json:"created_at"`
}

// NewNamespace builds a namespace with a fresh identifier.
func NewNamespace(name, code string, now time.Time) Namespace {
	return Namespace{
		ID:        uuid.NewString(),
		Name:      name,
		Code:      code,
		CreatedAt: now,
	}
}

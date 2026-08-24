package server

import (
	"fmt"
	"net/http"

	"lims/internal/config"
	"lims/internal/console"
	"lims/internal/store"
)

// App owns the file store and the console HTTP handler.
type App struct {
	Config  config.Config
	Stores  *store.FileStore
	Handler http.Handler
}

// NewApp wires every store into the console dependency set.
func NewApp(cfg config.Config) (*App, error) {
	fs, err := store.New(cfg.DataDir)
	if err != nil {
		return nil, fmt.Errorf("open file store: %w", err)
	}
	deps := console.Deps{
		Namespaces:  fs,
		Samples:     fs,
		Tasks:       fs,
		Instruments: fs,
		Results:     fs,
		Reports:     fs,
		Audits:      fs,
		Quotas:      fs,
		ReportFiles: fs,
	}
	return &App{Config: cfg, Stores: fs, Handler: console.NewRouter(deps)}, nil
}

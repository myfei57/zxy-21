package console

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"lims/internal/audit"
	"lims/internal/errs"
	"lims/internal/instrument"
	"lims/internal/quota"
	"lims/internal/task"
)

func listInstruments(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := d.Instruments.ListInstruments()
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, items)
	}
}

type createInstrumentRequest struct {
	NamespaceID string `json:"namespace_id"`
	Name        string `json:"name"`
	Model       string `json:"model"`
	Capacity    int    `json:"capacity"`
}

func createInstrument(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createInstrumentRequest
		if err := readJSON(r, &req); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if _, err := d.Namespaces.LoadNamespace(req.NamespaceID); err != nil {
			if errors.Is(err, errs.ErrNotFound) {
				writeErr(w, http.StatusBadRequest, "namespace not found")
				return
			}
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		created := instrument.NewInstrument(req.NamespaceID, req.Name, req.Model, req.Capacity, time.Now().UTC())
		if err := d.Instruments.SaveInstrument(created); err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		if _, err := quota.LoadOrCreate(d.Quotas, created.ID); err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, created)
	}
}

func reconnectInstrument(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		reconnected, err := instrument.Reconnect(d.Instruments, id, time.Now().UTC())
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		event := audit.NewEvent(
			audit.KindInstrumentReconnected,
			"instrument",
			id,
			"console",
			"instrument link recovered",
		)
		_ = audit.Record(d.Audits, event)
		writeJSON(w, http.StatusOK, reconnected)
	}
}

func disconnectInstrument(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		disconnected, err := instrument.Disconnect(d.Instruments, id)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		writeJSON(w, http.StatusOK, disconnected)
	}
}

type retryInstrumentRequest struct {
	Operator string `json:"operator"`
}

func retryInstrument(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var req retryInstrumentRequest
		if err := readJSON(r, &req); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid request body")
			return
		}
		dispatched, err := task.Retry(d.Tasks, d.Instruments, id, req.Operator)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"instrument_id": id,
			"dispatched":    dispatched,
		})
	}
}

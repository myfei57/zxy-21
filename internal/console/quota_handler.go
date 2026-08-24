package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"lims/internal/audit"
	"lims/internal/quota"
)

func listQuota(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		instruments, err := d.Instruments.ListInstruments()
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		rows := make([]quota.Limit, 0, len(instruments))
		for _, item := range instruments {
			limit, err := quota.LoadOrCreate(d.Quotas, item.ID)
			if err != nil {
				writeErr(w, http.StatusInternalServerError, err.Error())
				return
			}
			rows = append(rows, limit)
		}
		writeJSON(w, http.StatusOK, rows)
	}
}

type updateQuotaRequest struct {
	MaxConcurrent int    `json:"max_concurrent"`
	Operator      string `json:"operator"`
}

func updateQuota(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		instrumentID := chi.URLParam(r, "instrumentID")
		var req updateQuotaRequest
		if err := readJSON(r, &req); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid request body")
			return
		}
		limit, err := quota.SetLimit(d.Quotas, instrumentID, req.MaxConcurrent)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		event := audit.NewEvent(
			audit.KindQuotaUpdated,
			"quota",
			instrumentID,
			req.Operator,
			"quota limit updated",
		)
		_ = audit.Record(d.Audits, event)
		writeJSON(w, http.StatusOK, limit)
	}
}

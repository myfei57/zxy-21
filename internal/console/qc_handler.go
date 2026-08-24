package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"lims/internal/audit"
	"lims/internal/qc"
)

type runQCRequest struct {
	Analyte  string `json:"analyte"`
	Operator string `json:"operator"`
}

func runQC(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var req runQCRequest
		if err := readJSON(r, &req); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid request body")
			return
		}
		rule, err := qc.FindRule(qc.DefaultRuleSets(), req.Analyte)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		current, err := d.Results.LoadResult(id)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		judged, err := qc.Pass(d.Results, rule, current, req.Operator)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		event := audit.NewEvent(
			audit.KindQCPassed,
			"result",
			id,
			req.Operator,
			"QC verdict applied",
		)
		_ = audit.Record(d.Audits, event)
		writeJSON(w, http.StatusOK, judged)
	}
}

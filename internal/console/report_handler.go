package console

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"lims/internal/audit"
	"lims/internal/report"
)

func listReports(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := d.Reports.ListReports()
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, items)
	}
}

type createReportRequest struct {
	SampleID    string `json:"sample_id"`
	TaskID      string `json:"task_id"`
	NamespaceID string `json:"namespace_id"`
	Title       string `json:"title"`
}

func createReport(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createReportRequest
		if err := readJSON(r, &req); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid request body")
			return
		}
		created := report.NewReport(req.SampleID, req.TaskID, req.NamespaceID, req.Title, time.Now().UTC())
		if err := d.Reports.SaveReport(created); err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, created)
	}
}

type reportContentRequest struct {
	Content  string `json:"content"`
	Operator string `json:"operator"`
}

func signReport(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var req reportContentRequest
		if err := readJSON(r, &req); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid request body")
			return
		}
		current, err := d.Reports.LoadReport(id)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		signed, err := report.Sign(d.Reports, d.ReportFiles, current, []byte(req.Content))
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		event := audit.NewEvent(
			audit.KindReportSigned,
			"report",
			id,
			req.Operator,
			"report signed",
		)
		_ = audit.Record(d.Audits, event)
		writeJSON(w, http.StatusOK, signed)
	}
}

func reviseReport(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var req reportContentRequest
		if err := readJSON(r, &req); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid request body")
			return
		}
		current, err := d.Reports.LoadReport(id)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		revised, err := report.Revise(d.Reports, d.ReportFiles, current, []byte(req.Content))
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		event := audit.NewEvent(
			audit.KindReportRevised,
			"report",
			id,
			req.Operator,
			"report revision advanced",
		)
		_ = audit.Record(d.Audits, event)
		writeJSON(w, http.StatusOK, revised)
	}
}

func archiveReport(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		current, err := d.Reports.LoadReport(id)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		archived, err := report.Archive(d.Reports, d.ReportFiles, current)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		event := audit.NewEvent(
			audit.KindReportArchived,
			"report",
			id,
			"console",
			"report archived",
		)
		_ = audit.Record(d.Audits, event)
		writeJSON(w, http.StatusOK, archived)
	}
}

func getReportFile(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		current, err := d.Reports.LoadReport(id)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		content, err := report.ReadContent(d.ReportFiles, current.ContentPath)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write(content)
	}
}

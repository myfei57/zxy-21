package console

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"lims/internal/audit"
	"lims/internal/task"
)

func listTasks(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := d.Tasks.ListTasks()
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, items)
	}
}

type createTaskRequest struct {
	SampleID    string `json:"sample_id"`
	NamespaceID string `json:"namespace_id"`
	Panel       string `json:"panel"`
}

func createTask(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createTaskRequest
		if err := readJSON(r, &req); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid request body")
			return
		}
		created := task.NewTask(req.SampleID, req.NamespaceID, req.Panel, time.Now().UTC())
		if err := d.Tasks.SaveTask(created); err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, created)
	}
}

type assignTaskRequest struct {
	InstrumentID string `json:"instrument_id"`
	Operator     string `json:"operator"`
}

func assignTask(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var req assignTaskRequest
		if err := readJSON(r, &req); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid request body")
			return
		}
		current, err := d.Tasks.LoadTask(id)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		if _, err := d.Instruments.LoadInstrument(req.InstrumentID); err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		assigned, err := task.Assign(d.Tasks, d.Quotas, req.InstrumentID, current, req.Operator)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		event := audit.NewEvent(
			audit.KindTaskAssigned,
			"task",
			id,
			req.Operator,
			"assigned to instrument "+req.InstrumentID,
		)
		_ = audit.Record(d.Audits, event)
		writeJSON(w, http.StatusOK, assigned)
	}
}

type completeTaskRequest struct {
	Operator string `json:"operator"`
}

func completeTask(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var req completeTaskRequest
		if err := readJSON(r, &req); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid request body")
			return
		}
		completed, err := task.Complete(d.Tasks, d.Quotas, id, time.Now().UTC())
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		event := audit.NewEvent(
			audit.KindTaskCompleted,
			"task",
			id,
			req.Operator,
			"task completed and quota released",
		)
		_ = audit.Record(d.Audits, event)
		writeJSON(w, http.StatusOK, completed)
	}
}

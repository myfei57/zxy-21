package console

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"lims/internal/instrument"
	"lims/internal/result"
	"lims/internal/task"
)

func listResults(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := d.Results.ListResults()
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, items)
	}
}

type createResultRequest struct {
	SampleID     string               `json:"sample_id"`
	TaskID       string               `json:"task_id"`
	InstrumentID string               `json:"instrument_id"`
	Values       []result.Measurement `json:"values"`
}

func createResult(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createResultRequest
		if err := readJSON(r, &req); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid request body")
			return
		}
		now := time.Now().UTC()
		current, err := d.Tasks.LoadTask(req.TaskID)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		created := result.NewResult(req.SampleID, req.TaskID, req.InstrumentID, req.Values, now)
		if err := d.Results.SaveResult(created); err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		dispatch := instrument.NewInstrumentAssignment(
			current.ID,
			current.SampleID,
			req.InstrumentID,
			"instrument",
			now,
		)
		if err := instrument.RecordAssignment(d.Instruments, dispatch); err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		current.ResultID = created.ID
		current.ResultStatus = string(result.StatusPending)
		current.Status = task.StatusRunning
		if err := d.Tasks.SaveTask(current); err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, created)
	}
}

func getResult(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		view, err := result.View(d.Results, d.Tasks, id)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		writeJSON(w, http.StatusOK, view)
	}
}

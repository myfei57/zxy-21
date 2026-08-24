package console

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"lims/internal/audit"
	"lims/internal/errs"
	"lims/internal/sample"
)

func listSamples(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := d.Samples.ListSamples()
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, items)
	}
}

type registerSampleRequest struct {
	NamespaceID string `json:"namespace_id"`
	PatientID   string `json:"patient_id"`
	Kind        string `json:"kind"`
	Operator    string `json:"operator"`
	Note        string `json:"note"`
}

func registerSample(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req registerSampleRequest
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
		now := time.Now().UTC()
		item := sample.NewSample(req.NamespaceID, req.PatientID, req.Kind, now)
		record := sample.NewRegistrationRecord(item.ID, req.NamespaceID, req.Operator, req.Note, now)
		registered, err := sample.Register(d.Samples, d.Audits, record, item)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, registered)
	}
}

func getSample(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		item, err := d.Samples.LoadSample(id)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		record, err := d.Samples.LoadRegistration(id)
		if err != nil && !errors.Is(err, errs.ErrNotFound) {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"sample": item, "registration": record})
	}
}

type forwardSampleRequest struct {
	TargetStation string `json:"target_station"`
	BatchID       string `json:"batch_id"`
	Operator      string `json:"operator"`
}

func forwardSample(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var req forwardSampleRequest
		if err := readJSON(r, &req); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid request body")
			return
		}
		item, err := d.Samples.LoadSample(id)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		now := time.Now().UTC()
		batch := sample.FlowBatch{
			BatchID:       req.BatchID,
			SourceStation: item.CurrentStation,
			TargetStation: req.TargetStation,
			SampleIDs:     []string{id},
			ForwardedAt:   now,
		}
		if err := sample.Forward(d.Samples, d.Results, batch); err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		event := audit.NewEvent(
			audit.KindSampleForwarded,
			"sample",
			id,
			req.Operator,
			"forwarded to "+req.TargetStation,
		)
		_ = audit.Record(d.Audits, event)
		updated, err := d.Samples.LoadSample(id)
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, updated)
	}
}

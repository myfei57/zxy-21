package console

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"lims/internal/ns"
)

func listNamespaces(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := ns.List(d.Namespaces)
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, items)
	}
}

type createNamespaceRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func createNamespace(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createNamespaceRequest
		if err := readJSON(r, &req); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid request body")
			return
		}
		created, err := ns.Create(d.Namespaces, req.Name, req.Code, time.Now().UTC())
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, created)
	}
}

func disableNamespace(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		disabled, err := ns.Disable(d.Namespaces, id)
		if err != nil {
			writeErr(w, statusForError(err), err.Error())
			return
		}
		writeJSON(w, http.StatusOK, disabled)
	}
}

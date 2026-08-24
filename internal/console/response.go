package console

import (
	"encoding/json"
	"errors"
	"net/http"

	"lims/internal/errs"
)

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeErr(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func readJSON(r *http.Request, target any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func statusForError(err error) int {
	switch {
	case errors.Is(err, errs.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, errs.ErrAlreadyExists), errors.Is(err, errs.ErrQuotaExhausted):
		return http.StatusConflict
	default:
		return http.StatusBadRequest
	}
}

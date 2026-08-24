package console

import (
	"net/http"
	"strconv"

	"lims/internal/audit"
)

func listAudit(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := audit.Filter{
			Kind:       r.URL.Query().Get("kind"),
			EntityType: r.URL.Query().Get("entity_type"),
			EntityID:   r.URL.Query().Get("entity_id"),
		}
		limit := 200
		if raw := r.URL.Query().Get("limit"); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil || parsed < 0 {
				writeErr(w, http.StatusBadRequest, "invalid limit")
				return
			}
			limit = parsed
		}
		var events []audit.Event
		var err error
		if filter.Kind == "" && filter.EntityType == "" && filter.EntityID == "" {
			events, err = audit.Recent(d.Audits, limit)
		} else {
			filter.Limit = limit
			events, err = audit.Query(d.Audits, filter)
		}
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, events)
	}
}

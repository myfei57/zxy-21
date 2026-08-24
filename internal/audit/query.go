package audit

// Filter narrows audit event queries.
type Filter struct {
	Kind       string
	EntityType string
	EntityID   string
	Limit      int
}

// Query returns matching audit events, newest first.
func Query(store Store, filter Filter) ([]Event, error) {
	events, err := store.ListEvents()
	if err != nil {
		return nil, err
	}
	out := make([]Event, 0, len(events))
	for _, event := range events {
		if filter.Kind != "" && event.Kind != filter.Kind {
			continue
		}
		if filter.EntityType != "" && event.EntityType != filter.EntityType {
			continue
		}
		if filter.EntityID != "" && event.EntityID != filter.EntityID {
			continue
		}
		out = append(out, event)
	}
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	if filter.Limit > 0 && len(out) > filter.Limit {
		out = out[:filter.Limit]
	}
	return out, nil
}

// Recent returns the latest audit events.
func Recent(store Store, limit int) ([]Event, error) {
	return Query(store, Filter{Limit: limit})
}

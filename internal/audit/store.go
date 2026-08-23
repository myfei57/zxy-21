package audit

// Store persists audit events.
type Store interface {
	SaveEvent(event Event) error
	LoadEvent(id string) (Event, error)
	ListEvents() ([]Event, error)
}

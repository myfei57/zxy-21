package ns

// Store persists namespace entities.
type Store interface {
	SaveNamespace(namespace Namespace) error
	LoadNamespace(id string) (Namespace, error)
	ListNamespaces() ([]Namespace, error)
}

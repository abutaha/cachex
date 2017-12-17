package cachex

// Cache - interface to be implemented by cache engines
type Cache interface {
	Set(string, string) error
	Get(string) (string, error)
	Search([]byte) (map[string]string, error)
	Delete(string) error
	GetKeys() (map[string]bool, error)
}

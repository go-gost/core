package registry

type Registry[T any] interface {
	Register(name string, v T) error
	Unregister(name string)
	IsRegistered(name string) bool
	Get(name string) T
	GetAll() map[string]T
}

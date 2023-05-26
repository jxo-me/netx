package registry

type IRegistry[T any] interface {
	Register(name string, v T) error
	Unregister(name string)
	IsRegistered(name string) bool
	Get(name string) T
	GetAll() map[string]T
}

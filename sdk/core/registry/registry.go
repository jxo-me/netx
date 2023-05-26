package registry

import (
	"errors"
	"io"
	"sync"
)

var (
	ErrDup = errors.New("registry: duplicate object")
)

type registry[T any] struct {
	m sync.Map
}

func (r *registry[T]) Register(name string, v T) error {
	if name == "" {
		return nil
	}
	if _, loaded := r.m.LoadOrStore(name, v); loaded {
		return ErrDup
	}

	return nil
}

func (r *registry[T]) Unregister(name string) {
	if v, ok := r.m.Load(name); ok {
		if closer, ok := v.(io.Closer); ok {
			closer.Close()
		}
		r.m.Delete(name)
	}
}

func (r *registry[T]) IsRegistered(name string) bool {
	_, ok := r.m.Load(name)
	return ok
}

func (r *registry[T]) Get(name string) (t T) {
	if name == "" {
		return
	}
	v, _ := r.m.Load(name)
	t, _ = v.(T)
	return
}

func (r *registry[T]) GetAll() (m map[string]T) {
	m = make(map[string]T)
	r.m.Range(func(key, value any) bool {
		k, _ := key.(string)
		v, _ := value.(T)
		m[k] = v
		return true
	})
	return
}

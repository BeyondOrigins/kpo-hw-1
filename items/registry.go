package items

import (
	"errors"
	"fmt"
	"sort"
)

var ErrUnknownKind = errors.New("неизвестный вид")

type Spec struct {
	ID       ID
	Name     string
	Easiness int
	KWh      int
}

type Factory func(Spec) (Item, error)

type Registry struct {
	factories map[string]Factory
}

func NewRegistry() *Registry {
	return &Registry{factories: make(map[string]Factory)}
}

func (r *Registry) Register(kind string, f Factory) error {
	if f == nil {
		return fmt.Errorf("register %q: nil factory", kind)
	}
	if _, exists := r.factories[kind]; exists {
		return fmt.Errorf("register %q: already registered", kind)
	}
	r.factories[kind] = f
	return nil
}

func (r *Registry) New(kind string, spec Spec) (Item, error) {
	f, ok := r.factories[kind]
	if !ok {
		return nil, fmt.Errorf("%w %q", ErrUnknownKind, kind)
	}
	return f(spec)
}

func (r *Registry) Kinds() []string {
	kinds := make([]string, 0, len(r.factories))
	for k := range r.factories {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)
	return kinds
}

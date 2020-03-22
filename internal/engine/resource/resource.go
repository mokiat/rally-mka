package resource

import "sync"

type resourceID int

type state int

const (
	stateUnloaded state = iota
	stateLoading
	stateLoaded
)

type Resource interface{}

type Handle struct {
	registry *Registry
	id       resourceID
}

func (h Handle) Get() Resource {
	reference := h.registry.references[h.id]
	return reference.Resource
}

func (h Handle) IsAvailable() bool {
	reference := h.registry.references[h.id]
	return reference.State == stateLoaded
}

type Type struct {
	registry *Registry
	operator Operator

	// XXX: Locking would not be necessary if all resources were pre-registered
	handlesLock sync.Mutex
	handles     map[string]Handle
}

func (t *Type) Resource(name string) Handle {
	t.handlesLock.Lock()
	defer t.handlesLock.Unlock()

	if handle, found := t.handles[name]; found {
		return handle
	}

	handle := Handle{
		registry: t.registry,
		id:       t.registry.allocateReference(name, t.operator),
	}
	t.handles[name] = handle
	return handle
}

type Operator interface {
	Allocate(registry *Registry, name string) (Resource, error)
	Release(registry *Registry, resource Resource) error
}

package resource

import (
	"fmt"
	"sync/atomic"
)

func NewRegistry(worker *Worker, maxResources, maxEvents int) *Registry {
	return &Registry{
		catalog:        make(map[string]*Type),
		references:     make([]reference, maxResources),
		referenceCount: 0,
		worker:         worker,
		events:         make(chan event, maxEvents),
	}
}

type Registry struct {
	catalog map[string]*Type

	references     []reference
	referenceCount int32

	worker *Worker

	// XXX: Multiple priority queues can be used if ever certain events
	// need to be processed with priority / certainty
	events chan event
}

func (r *Registry) RegisterType(resType string, operator Operator) {
	r.catalog[resType] = &Type{
		registry: r,
		operator: operator,
		handles:  make(map[string]Handle),
	}
}

func (r *Registry) Type(name string) *Type {
	return r.catalog[name]
}

func (r *Registry) Request(handle Handle) {
	r.events <- requestEvent{
		ID: handle.id,
	}
}

func (r *Registry) Dismiss(handle Handle) {
	r.events <- dismissEvent{
		ID: handle.id,
	}
}

// It is of upmost importance that nothing resource specific happens
// during this method call (only Registry methods may be used)
// Resource / Handle methods should not be used.
func (r *Registry) Update() {
	for e, ok := r.popEvent(); ok; e, ok = r.popEvent() {
		r.processEvent(e)
	}
}

func (r *Registry) popEvent() (event, bool) {
	// XXX: Rate limiting can be implemented if excessive number of events
	// start to pile up
	select {
	case event := <-r.events:
		return event, true
	default:
		return nil, false
	}
}

func (r *Registry) processEvent(e event) {
	switch specificEvent := e.(type) {
	case requestEvent:
		r.references[specificEvent.ID].Users++
		r.checkResource(specificEvent.ID)
	case dismissEvent:
		r.checkResource(specificEvent.ID)
	case loadedEvent:
		r.saveResource(specificEvent.ID, specificEvent.Resource)
	}
}

func (r *Registry) checkResource(id resourceID) {
	ref := &r.references[id]

	switch ref.State {
	case stateUnloaded:
		if ref.IsDesired() {
			ref.State = stateLoading
			r.scheduleResourceLoad(id)
		}
	case stateLoaded:
		if !ref.IsDesired() {
			r.scheduleResourceUnload(id, ref.Resource)
			ref.Resource = nil
			ref.State = stateUnloaded
		}
	}
}

func (r *Registry) saveResource(id resourceID, resource Resource) {
	ref := &r.references[id]

	switch ref.State {
	case stateUnloaded:
		if ref.IsDesired() {
			ref.Resource = resource
			ref.State = stateLoaded
		} else {
			r.scheduleResourceUnload(id, resource)
		}
	case stateLoading:
		if ref.IsDesired() {
			ref.Resource = resource
			ref.State = stateLoaded
		} else {
			ref.State = stateUnloaded
			r.scheduleResourceUnload(id, resource)
		}
	case stateLoaded:
		oldResource := ref.Resource
		ref.Resource = resource
		r.scheduleResourceUnload(id, oldResource)
	}
}

func (r *Registry) scheduleResourceLoad(id resourceID) {
	r.worker.Schedule(func() error {
		ref := r.references[id]
		resource, err := ref.Operator.Allocate(r, ref.Name)
		if err != nil {
			return fmt.Errorf("failed to allocate resource %q: %w", ref.Name, err)
		}
		r.events <- loadedEvent{
			ID:       id,
			Resource: resource,
		}
		return nil
	})
}

func (r *Registry) scheduleResourceUnload(id resourceID, resource Resource) {
	r.worker.Schedule(func() error {
		ref := r.references[id]
		if err := ref.Operator.Release(r, resource); err != nil {
			return fmt.Errorf("failed to release resource %q: %w", ref.Name, err)
		}
		return nil
	})
}

func (r *Registry) allocateReference(name string, operator Operator) resourceID {
	count := atomic.AddInt32(&r.referenceCount, 1)
	id := resourceID(count - 1)
	r.references[id] = reference{
		Operator: operator,
		Name:     name,
		Resource: nil,
		State:    stateUnloaded,
		Users:    0,
	}
	return id
}

type reference struct {
	Operator Operator
	Name     string
	Resource Resource
	State    state
	Users    int
}

func (r reference) IsDesired() bool {
	return r.Users > 0
}

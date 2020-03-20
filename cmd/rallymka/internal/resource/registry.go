package resource

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
)

type Handle struct {
	name       string
	index      int
	controller Controller
	users      int32
	state      int32
}

func (h *Handle) Name() string {
	return h.name
}

func (h *Handle) Request() {
	atomic.AddInt32(&h.users, 1)
	// TODO: register event
}

func (h *Handle) Dismiss() {
	atomic.AddInt32(&h.users, -1)
}

func (h *Handle) Available() bool {
	return atomic.LoadInt32(&h.state) == stateLoaded
}

// TODO: Configurable
const registryCapacity = 2048
const eventQueueSize = 256

func NewRegistry(locator Locator, worker *Worker) *Registry {
	registry := &Registry{
		locator: locator,
		worker:  worker,

		resourceTypes: make(map[string]*ResourceType),

		handleCount: 0,
		handles:     make([]Handle, registryCapacity),

		events:        make(chan int, eventQueueSize),
		updateRequest: make(chan struct{}, 1),
	}
	go registry.processEvents()
	return registry
}

type Registry struct {
	locator Locator
	worker  *Worker

	resourceTypes map[string]*ResourceType

	handleCount int32
	handles     []Handle

	events        chan int
	updateRequest chan struct{}
}

func (r *Registry) RegisterResource(controller Controller) {
	name := controller.ResourceTypeName()
	if _, found := r.resourceTypes[name]; !found {
		r.resourceTypes[name] = &ResourceType{
			registry:    r,
			controller:  controller,
			handlesLock: &sync.Mutex{},
			handles:     make(map[string]Resource),
			handleCount: 0,
		}
	}
}

func (r *Registry) ResourceType(name string) *ResourceType {
	return r.resourceTypes[name]
}

func (r *Registry) Update() {
	select {
	case r.updateRequest <- struct{}{}:
	default:
	}
}

func (r *Registry) allocateHandle() *Handle {
	count := atomic.AddInt32(&r.handleCount, 1)
	return &r.handles[count-1]
}

func (r *Registry) processEvents() {
	for {
		select {
		case id := <-r.events:
			r.evaluateHandle(id)
		case <-r.updateRequest:
			count := atomic.LoadInt32(&r.handleCount)
			for id := 0; id < int(count); id++ {
				r.evaluateHandle(id)
			}
		}
	}
}

func (r *Registry) evaluateHandle(id int) {
	handle := &r.handles[id]
	users := atomic.LoadInt32(&handle.users)
	state := atomic.LoadInt32(&handle.state)

	switch {
	case users > 0 && state == stateUnloaded:
		if atomic.CompareAndSwapInt32(&handle.state, state, stateLoading) {
			r.loadResource(id)
		}

	case users == 0 && state == stateLoaded:
		if atomic.CompareAndSwapInt32(&handle.state, state, stateUnloading) {
			r.unloadResource(id)
		}
	}
}

func (r *Registry) loadResource(id int) {
	handle := &r.handles[id]
	log.Printf("loading resource: %s", handle.name)

	r.worker.Schedule(func() error {
		if err := handle.controller.Load(handle.index, r.locator, r); err != nil {
			return fmt.Errorf("failed to load resource: %w", err)
		}
		atomic.StoreInt32(&handle.state, stateLoaded)
		return nil
	})
}

func (r *Registry) unloadResource(id int) {
	handle := &r.handles[id]
	log.Printf("unloading resource: %s", handle.name)

	r.worker.Schedule(func() error {
		if err := handle.controller.Unload(handle.index); err != nil {
			return fmt.Errorf("failed to unload resource: %w", err)
		}
		atomic.StoreInt32(&handle.state, stateUnloaded)
		return nil
	})
}

type ResourceType struct {
	registry   *Registry
	controller Controller

	handlesLock *sync.Mutex
	handles     map[string]Resource
	handleCount int
}

func (r *ResourceType) Resource(name string) Resource {
	r.handlesLock.Lock()
	defer r.handlesLock.Unlock()

	r.handleCount++

	resourceIndex := r.handleCount - 1
	handle := r.registry.allocateHandle()
	handle.name = name
	handle.index = resourceIndex
	handle.controller = r.controller
	return r.controller.Init(resourceIndex, handle)
}

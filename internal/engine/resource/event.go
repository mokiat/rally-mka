package resource

type event interface {
}

type requestEvent struct {
	ID resourceID
}

type dismissEvent struct {
	ID resourceID
}

type loadedEvent struct {
	ID       resourceID
	Resource Resource
}

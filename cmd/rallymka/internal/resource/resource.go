package resource

type Resource interface{}

const (
	stateUnloaded int32 = iota
	stateUnloading
	stateLoading
	stateLoaded
)

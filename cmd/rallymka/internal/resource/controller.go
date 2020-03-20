package resource

type Controller interface {
	ResourceTypeName() string
	Init(index int, handle *Handle) Resource
	Load(index int, locator Locator, registry *Registry) error
	Unload(index int) error
}

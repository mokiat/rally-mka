package input

import "sync"

type ActionSet struct {
	Forward     bool
	Backward    bool
	Left        bool
	Right       bool
	Handbrake   bool
	FreezeFrame bool
}

type Tracker struct {
	actionsLock sync.Mutex
	actions     ActionSet
}

func (i *Tracker) Set(actions ActionSet) {
	i.actionsLock.Lock()
	defer i.actionsLock.Unlock()
	i.actions = actions
}

func (i *Tracker) Get() ActionSet {
	i.actionsLock.Lock()
	defer i.actionsLock.Unlock()
	return i.actions
}

package model

import "github.com/mokiat/lacking/ui/mvc"

type LoadingPromise interface {
	OnReady(func())
}

func NewLoading(eventBus *mvc.EventBus) *Loading {
	return &Loading{
		promise:      nil,
		nextViewName: ViewNameIntro,
	}
}

type Loading struct {
	promise      LoadingPromise
	nextViewName ViewName
}

func (l *Loading) Promise() LoadingPromise {
	return l.promise
}

func (l *Loading) SetPromise(promise LoadingPromise) {
	l.promise = promise
}

func (l *Loading) NextViewName() ViewName {
	return l.nextViewName
}

func (l *Loading) SetNextViewName(name ViewName) {
	l.nextViewName = name
}

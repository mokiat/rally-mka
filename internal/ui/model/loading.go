package model

import "github.com/mokiat/lacking/ui/mvc"

var (
	LoadingChange         = mvc.NewChange("loading")
	LoadingPromiseChange  = mvc.SubChange(LoadingChange, "promise")
	LoadingNextViewChange = mvc.SubChange(LoadingChange, "next_view")
)

type LoadingPromise interface {
	OnReady(func())
}

func newLoading() *Loading {
	return &Loading{
		Observable:   mvc.NewObservable(),
		promise:      nil,
		nextViewName: ViewNameIntro,
	}
}

type Loading struct {
	mvc.Observable
	promise      LoadingPromise
	nextViewName ViewName
}

func (l *Loading) Promise() LoadingPromise {
	return l.promise
}

func (l *Loading) SetPromise(promise LoadingPromise) {
	l.promise = promise
	l.SignalChange(LoadingPromiseChange)
}

func (l *Loading) NextViewName() ViewName {
	return l.nextViewName
}

func (l *Loading) SetNextViewName(name ViewName) {
	l.nextViewName = name
	l.SignalChange(LoadingNextViewChange)
}

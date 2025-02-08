package model

import (
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/util/async"
)

type LoadingPromise interface {
	OnSuccess(func())
	OnError(func())
}

func NewLoadingPromise[T any](worker game.Worker, promise async.Promise[T], onSuccess func(T), onError func(error)) LoadingPromise {
	return &loadingPromise[T]{
		worker:    worker,
		promise:   promise,
		onSuccess: onSuccess,
		onError:   onError,
	}
}

type loadingPromise[T any] struct {
	worker    game.Worker
	promise   async.Promise[T]
	onSuccess func(T)
	onError   func(error)
}

func (p *loadingPromise[T]) OnSuccess(cb func()) {
	p.promise.OnSuccess(func(value T) {
		p.worker.Schedule(func() {
			p.onSuccess(value)
			cb()
		})
	})
}

func (p *loadingPromise[T]) OnError(cb func()) {
	p.promise.OnError(func(err error) {
		p.worker.Schedule(func() {
			p.onError(err)
			cb()
		})
	})
}

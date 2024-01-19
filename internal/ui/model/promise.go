package model

import "github.com/mokiat/lacking/util/async"

type LoadingPromise interface {
	OnReady(func())
}

func ToLoadingPromise[T any](promise async.Promise[T]) LoadingPromise {
	return &loadingPromise[T]{
		promise: promise,
	}
}

type loadingPromise[T any] struct {
	promise async.Promise[T]
}

func (p *loadingPromise[T]) OnReady(cb func()) {
	p.promise.OnReady(cb)
}

func (p *loadingPromise[T]) Err() error {
	_, err := p.promise.Wait()
	return err
}

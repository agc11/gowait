package gowait

import "context"

type Result[T any] struct {
	Value T
	Error error
}

type Future[T any] struct {
	await func(ctx context.Context) Result[T]
}

func (f Future[T]) Await(ctx context.Context) Result[T] {
	return f.await(ctx)
}

func NewFuture[T any](ctx context.Context, callback func(ctx context.Context) (T, error)) Future[T] {
	c := make(chan Result[T])
	response := Result[T]{}
	go func() {
		defer close(c)
		value, err := callback(ctx)
		response.Error = err
		response.Value = value
	}()
	return Future[T]{
		await: func(ctx context.Context) Result[T] {
			select {
			case <-ctx.Done():
				response.Error = ctx.Err()
				return response
			case <-c:
				return response
			}
		},
	}
}

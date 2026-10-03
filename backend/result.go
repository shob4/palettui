package main

type Result[T any] struct {
	value T
	err   error
}

func Ok[T any](val T) Result[T] {
	return Result[T]{value: val}
}

func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

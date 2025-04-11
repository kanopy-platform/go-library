package option

import "fmt"

// Option represents a value that may or may not be present.
type Option[T any] struct {
    value *T
}

// Some creates an Option with a value.
func Some[T any](value T) Option[T] {
    return Option[T]{&value}
}

// None creates an empty Option.
func None[T any]() Option[T] {
    return Option[T]{}
}

func (o Option[T]) Filter(predicate func(val T) bool) Option[T] {
	if o.value != nil && predicate(*o.value) {
		return Some(*o.value)
	}
	return None[T]()
}

func (o Option[T]) Inspect(f func(T) Option[T]) Option[T] {
	if o.value != nil {
		f(*o.value)
	}
	return o
}

func (o Option[T]) Get() (*T, error) {
	if o.IsSome() {
		return o.value, nil
	}
	return nil, fmt.Errorf("option has no value")
}

func (o Option[T]) Or(other Option[T]) Option[T] {
	if o.value != nil {
		return Some(*o.value)
	}
	return other
}

func (o Option[T]) OrElse(f func() Option[T]) Option[T] {
	if o.value != nil {
		return Some(*o.value)
	}
	return f()
}

func (o Option[T]) IsSome() bool {
	return o.value != nil
}

func Map[S any, D any](o Option[S], f func(S) D) Option[D] {
	if o.value != nil {
		v := f(*o.value)
		return Some(v)
	}
	return None[D]()
}

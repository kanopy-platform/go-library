package option

import "fmt"

func Some[T any](value T) Option[T] {
	return Option[T]{&value}
}

func None[T any]() Option[T] {
	return Option[T]{}
}

type Option[T any] struct {
	value *T
}

func (o Option[T]) Filter(predicate func(val T) bool) Option[T] {
	value, _ := o.Get()
	if value != nil && predicate(*value) {
		return Some(*value)
	}
	return None[T]()
}

func (o Option[T]) Inspect(f func(T) Option[T]) Option[T] {
	value, _ := o.Get()
	if value != nil {
		f(*value)
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
	value, _ := o.Get()
	if value != nil {
		return Some(*value)
	}
	return other
}

func (o Option[T]) OrElse(f func() Option[T]) Option[T] {
	value, _ := o.Get()
	if value != nil {
		return Some(*value)
	}
	return f()
}

func (o Option[T]) IsSome() bool {
	return o.value != nil
}

func Map[S any, D any](o Option[S], f func(S) D) Option[D] {
	value, _ := o.Get()
	if value != nil {
		v := f(*value)
		return Some(v)
	}
	return None[D]()
}

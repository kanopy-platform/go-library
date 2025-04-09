package option

import "fmt"

type Option[T any] interface {
	Filter(func(val T) bool) Option[T]
	Inspect(func(T) Option[T]) Option[T]
	Get() (*T, error)
	IsSome() bool
	Or(Option[T]) Option[T]
	OrElse(func() Option[T]) Option[T]
}

func Some[T any](value T) Option[T] {
	return &option[T]{&value}
}

func None[T any]() Option[T] {
	return &option[T]{}
}

type option[T any] struct {
	value *T
}

func (o *option[T]) Filter(predicate func(val T) bool) Option[T] {
	value, _ := o.Get()
	if value != nil && predicate(*value) {
		return Some(*value)
	}
	return None[T]()
}

func (o *option[T]) Inspect(f func(T) Option[T]) Option[T] {
	value, _ := o.Get()
	if value != nil {
		f(*value)
	}
	return o
}

func (o *option[T]) Get() (*T, error) {
	if o.IsSome() {
		return o.value, nil
	}
	return nil, fmt.Errorf("option has no value")
}

func (o *option[T]) Or(other Option[T]) Option[T] {
	value, _ := o.Get()
	if value != nil {
		return Some(*value)
	}
	return other
}

func (o *option[T]) OrElse(f func() Option[T]) Option[T] {
	value, _ := o.Get()
	if value != nil {
		return Some(*value)
	}
	return f()
}

func (o *option[T]) IsSome() bool {
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

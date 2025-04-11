package opt

import "fmt"

// Opt represents a value that may or may not be present.
type Opt[T any] struct {
    value *T
}

// Some creates an Opt with a value.
func Some[T any](value T) Opt[T] {
    return Opt[T]{&value}
}

// None creates an empty Opt.
func None[T any]() Opt[T] {
    return Opt[T]{}
}

func (o Opt[T]) Filter(predicate func(val T) bool) Opt[T] {
	if o.value != nil && predicate(*o.value) {
		return Some(*o.value)
	}
	return None[T]()
}

func (o Opt[T]) Inspect(f func(T) Opt[T]) Opt[T] {
	if o.value != nil {
		f(*o.value)
	}
	return o
}

func (o Opt[T]) Get() (*T, error) {
	if o.IsSome() {
		return o.value, nil
	}
	return nil, fmt.Errorf("opt has no value")
}

func (o Opt[T]) Or(other Opt[T]) Opt[T] {
	if o.value != nil {
		return Some(*o.value)
	}
	return other
}

func (o Opt[T]) OrElse(f func() Opt[T]) Opt[T] {
	if o.value != nil {
		return Some(*o.value)
	}
	return f()
}

func (o Opt[T]) IsSome() bool {
	return o.value != nil
}

func Map[S any, D any](o Opt[S], f func(S) D) Opt[D] {
	if o.value != nil {
		v := f(*o.value)
		return Some(v)
	}
	return None[D]()
}

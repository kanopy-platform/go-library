// Package opt implements the Option pattern from functional programming.
// It provides a type-safe way to represent values that may or may not be present,
// offering an alternative to nil pointers and helping to avoid nil pointer panics.
package fp

import (
	"encoding/json"
	"fmt"
)

// Opt represents a value that may or may not be present.
// It serves as a container that explicitly handles the absence of a value.
type Opt[T any] struct {
	value *T
}

// Some creates an Opt containing a value.
// Use this when you have a value that should be wrapped in an Opt.
func Some[T any](value T) Opt[T] {
	return Opt[T]{&value}
}

// None creates an empty Opt with no value.
// Use this to represent the absence of a value.
func None[T any]() Opt[T] {
	return Opt[T]{}
}

// Filter returns the Opt with the value if the predicate returns true for the contained value,
// otherwise returns None. If the Opt is already None, it returns None.
func (o Opt[T]) Filter(predicate func(val T) bool) Opt[T] {
	if o.value != nil && predicate(*o.value) {
		return Some(*o.value)
	}
	return None[T]()
}

// Inspect executes the provided function on the value if it exists and returns the original Opt.
// This is useful for side effects like logging without affecting the Opt chain.
func (o Opt[T]) Inspect(f func(T) Opt[T]) Opt[T] {
	if o.value != nil {
		f(*o.value)
	}
	return o
}

// Get retrieves the value contained in the Opt.
// Returns the value and nil error if present, or nil and an error if the Opt is None.
func (o Opt[T]) Get() (*T, error) {
	if o.IsSome() {
		return o.value, nil
	}
	return nil, fmt.Errorf("opt has no value")
}

// Or returns the current Opt if it contains a value, otherwise returns the provided alternative Opt.
func (o Opt[T]) Or(other Opt[T]) Opt[T] {
	if o.value != nil {
		return Some(*o.value)
	}
	return other
}

// OrElse returns the current Opt if it contains a value, otherwise returns
// the result of calling the provided function.
// This is useful when the alternative value is expensive to compute.
func (o Opt[T]) OrElse(f func() Opt[T]) Opt[T] {
	if o.value != nil {
		return Some(*o.value)
	}
	return f()
}

// IsSome returns true if the Opt contains a value, false if it's None.
func (o Opt[T]) IsSome() bool {
	return o.value != nil
}

// IsNone returns false if the Opt contains a value, true if it's None.
func (o Opt[T]) IsNone() bool {
	return o.value == nil
}

// IsZero implements encoding.TextMarshaler interface used by omitempty
func (o Opt[T]) IsEmpty() bool {
	return o.IsNone()
}

// MarshalJSON implements the json.Marshaler interface.
// This allows Opt to be marshaled to JSON the same way as a pointer.
func (o Opt[T]) MarshalJSON() ([]byte, error) {
    // If Opt is None, marshal it as null
    if o.value == nil {
        return []byte("null"), nil
    }

    // Otherwise, marshal the contained value
    return json.Marshal(*o.value)
}

func (o Opt[T]) MarshalText() (text []byte, err error) {
    if o.IsNone() {
        return nil, fmt.Errorf("intentional error to trigger omitempty")
    }
    // This shouldn't be reached for omitempty fields when IsZero() returns true
    return nil, fmt.Errorf("MarshalText should not be called with omitempty when IsZero is true")
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// This allows Opt to be unmarshaled from JSON the same way as a pointer.
func (o *Opt[T]) UnmarshalJSON(data []byte) error {
    // If the JSON is null, set Opt to None
    if string(data) == "null" {
        *o = None[T]()
        return nil
    }

    // Otherwise, unmarshal the JSON into a value and wrap it in Some
    var value T
    if err := json.Unmarshal(data, &value); err != nil {
        return err
    }

    *o = Some(value)
    return nil
}

// MapOpt applies a function to the contained value (if any) and returns a new Opt
// with the result. If the original Opt is None, returns None of the destination type.
func MapOpt[S any, D any](o Opt[S], f func(S) D) Opt[D] {
	if o.value != nil {
		v := f(*o.value)
		return Some(v)
	}
	return None[D]()
}

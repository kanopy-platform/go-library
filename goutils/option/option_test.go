package option_test

import (
	"strconv"
	"testing"

	"github.com/kanopy-platform/go-library/goutils/option"
)

func TestSomeAndNone(t *testing.T) {
    someValue := option.Some(64)
    if !someValue.IsSome() {
        t.Error("Expected Some(64) to have IsSome() == true")
    }

    value, err := someValue.Get()
    if err != nil || *value != 64 {
        t.Errorf("Expected Some(64).Get() to return 64, got %v with error %v", value, err)
    }

    noneValue := option.None[int]()
    if noneValue.IsSome() {
        t.Error("Expected None to have IsSome() == false")
    }

    value, err = noneValue.Get()
    if err == nil || value != nil {
        t.Errorf("Expected None.Get() to return error, got %v with error %v", value, err)
    }
}

func TestFilter(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

    someEven := option.Some(64)
    filtered := someEven.Filter(isEven)
    if !filtered.IsSome() {
        t.Error("Filter with passing predicate should return Some")
    }

    someOdd := option.Some(43)
    filtered = someOdd.Filter(isEven)
    if filtered.IsSome() {
        t.Error("Filter with failing predicate should return None")
    }

    noneValue := option.None[int]()
    filtered = noneValue.Filter(func(v int) bool { return true })
    if filtered.IsSome() {
        t.Error("Filtering None should always return None")
    }
}

func TestInspect(t *testing.T) {
    someValue := option.Some(64)
    result := someValue.Inspect(func(v int) option.Option[int] {
        return option.Some(v * 2)
    })
    value, _ := result.Get()
    if value == nil || *value != 64 {
        t.Error("Inspect should return the result of applying the function")
    }

    result = someValue.Inspect(func(v int) option.Option[int] {
        return option.None[int]()
    })
    if !result.IsSome() {
        t.Error("Inspect to None should return None")
    }

    noneValue := option.None[int]()
    result = noneValue.Inspect(func(v int) option.Option[int] {
        t.Error("Inspect should not call function on None")
        return option.Some(v)
    })
    if result.IsSome() {
        t.Error("Inspect on None should return None")
    }
}

func TestOr(t *testing.T) {
    someValue := option.Some(64)
    result := someValue.Or(option.Some(128))
    value, _ := result.Get()
    if *value != 64 {
        t.Errorf("Or with Some should return original valueue, got %v", *value)
    }

    noneValue := option.None[int]()
    result = noneValue.Or(option.Some(128))
    value, _ = result.Get()
    if *value != 128 {
        t.Errorf("Or with None should return alternative, got %v", *value)
    }

    result = noneValue.Or(option.None[int]())
    if result.IsSome() {
        t.Error("Or with None and None alternative should return None")
    }
}

func TestOrElse(t *testing.T) {
    someValue := option.Some(64)
    result := someValue.OrElse(func() option.Option[int] {
        t.Error("OrElse should not call function on Some")
        return option.Some(128)
    })
    value, _ := result.Get()
    if *value != 64 {
        t.Errorf("OrElse with Some should return original valueue, got %v", *value)
    }

    noneValue := option.None[int]()
    result = noneValue.OrElse(func() option.Option[int] {
        return option.Some(128)
    })
    value, _ = result.Get()
    if *value != 128 {
        t.Errorf("OrElse with None should return function result, got %v", *value)
    }

    result = noneValue.OrElse(func() option.Option[int] {
        return option.None[int]()
    })
    if result.IsSome() {
        t.Error("OrElse with None returning None should return None")
    }
}

func TestMap(t *testing.T) {
    someValue := option.Some(64)
    result := option.Map(someValue, func(v int) string {
        return strconv.Itoa(v)
    })
    value, _ := result.Get()
    if *value != "64" {
        t.Errorf("Map should transform valueue, got %v", *value)
    }

    noneValue := option.None[int]()
    result = option.Map(noneValue, func(v int) string {
        t.Error("Map should not call function on None")
        return strconv.Itoa(v)
    })
    if result.IsSome() {
        t.Error("Map on None should return None")
    }

    type Person struct {
        Name string
        Age int
    }

    somePerson := option.Some(Person{Name: "Alice", Age: 30})
    result = option.Map(somePerson, func(p Person) string {
        return p.Name
    })
    value, _ = result.Get()
    if *value != "Alice" {
        t.Errorf("Map should transform struct to string, got %v", *value)
    }
}

func TestEdgeCases(t *testing.T) {
    zeroInt := option.Some(0)
    if !zeroInt.IsSome() {
        t.Error("Some(0) should have IsSome() == true")
    }

    zeroString := option.Some("")
    if !zeroString.IsSome() {
        t.Error("Some(\"\") should have IsSome() == true")
    }

    var nilPtr *string = nil
    someNil := option.Some(nilPtr)
    if !someNil.IsSome() {
        t.Error("Some(nil) should have IsSome() == true")
    }

    nestedOption := option.Some(option.Some(64))
    extracted := option.Map(nestedOption, func(opt option.Option[int]) int {
        value, _ := opt.Get()
        return *value
    })
    value, _ := extracted.Get()
    if *value != 64 {
        t.Errorf("Extracting from nested option should work, got %v", *value)
    }
}

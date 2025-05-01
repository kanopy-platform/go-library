package fp_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/kanopy-platform/go-library/fp"
)

func TestSomeAndNone(t *testing.T) {
	someValue := fp.Some(64)
	if !someValue.IsSome() {
		t.Error("Expected Some(64) to have IsSome() == true")
	}

	value, err := someValue.Get()
	if err != nil || *value != 64 {
		t.Errorf("Expected Some(64).Get() to return 64, got %v with error %v", value, err)
	}

	noneValue := fp.None[int]()
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

	someEven := fp.Some(64)
	filtered := someEven.Filter(isEven)
	if !filtered.IsSome() {
		t.Error("Filter with passing predicate should return Some")
	}

	someOdd := fp.Some(43)
	filtered = someOdd.Filter(isEven)
	if filtered.IsSome() {
		t.Error("Filter with failing predicate should return None")
	}

	noneValue := fp.None[int]()
	filtered = noneValue.Filter(func(v int) bool { return true })
	if filtered.IsSome() {
		t.Error("Filtering None should always return None")
	}
}

func TestInspect(t *testing.T) {
	someValue := fp.Some(64)
	result := someValue.Inspect(func(v int) fp.Opt[int] {
		return fp.Some(v * 2)
	})
	value, _ := result.Get()
	if value == nil || *value != 64 {
		t.Error("Inspect should return the result of applying the function")
	}

	result = someValue.Inspect(func(v int) fp.Opt[int] {
		return fp.None[int]()
	})
	if !result.IsSome() {
		t.Error("Inspect to None should return None")
	}

	noneValue := fp.None[int]()
	result = noneValue.Inspect(func(v int) fp.Opt[int] {
		t.Error("Inspect should not call function on None")
		return fp.Some(v)
	})
	if result.IsSome() {
		t.Error("Inspect on None should return None")
	}
}

func TestOr(t *testing.T) {
	someValue := fp.Some(64)
	result := someValue.Or(fp.Some(128))
	value, _ := result.Get()
	if *value != 64 {
		t.Errorf("Or with Some should return original valueue, got %v", *value)
	}

	noneValue := fp.None[int]()
	result = noneValue.Or(fp.Some(128))
	value, _ = result.Get()
	if *value != 128 {
		t.Errorf("Or with None should return alternative, got %v", *value)
	}

	result = noneValue.Or(fp.None[int]())
	if result.IsSome() {
		t.Error("Or with None and None alternative should return None")
	}
}

func TestOrElse(t *testing.T) {
	someValue := fp.Some(64)
	result := someValue.OrElse(func() fp.Opt[int] {
		t.Error("OrElse should not call function on Some")
		return fp.Some(128)
	})
	value, _ := result.Get()
	if *value != 64 {
		t.Errorf("OrElse with Some should return original valueue, got %v", *value)
	}

	noneValue := fp.None[int]()
	result = noneValue.OrElse(func() fp.Opt[int] {
		return fp.Some(128)
	})
	value, _ = result.Get()
	if *value != 128 {
		t.Errorf("OrElse with None should return function result, got %v", *value)
	}

	result = noneValue.OrElse(func() fp.Opt[int] {
		return fp.None[int]()
	})
	if result.IsSome() {
		t.Error("OrElse with None returning None should return None")
	}
}

func TestMap(t *testing.T) {
	someValue := fp.Some(64)
	result := fp.Map(someValue, func(v int) string {
		return strconv.Itoa(v)
	})
	value, _ := result.Get()
	if *value != "64" {
		t.Errorf("Map should transform valueue, got %v", *value)
	}

	noneValue := fp.None[int]()
	result = fp.Map(noneValue, func(v int) string {
		t.Error("Map should not call function on None")
		return strconv.Itoa(v)
	})
	if result.IsSome() {
		t.Error("Map on None should return None")
	}

	type Person struct {
		Name string
		Age  int
	}

	somePerson := fp.Some(Person{Name: "Alice", Age: 30})
	result = fp.Map(somePerson, func(p Person) string {
		return p.Name
	})
	value, _ = result.Get()
	if *value != "Alice" {
		t.Errorf("Map should transform struct to string, got %v", *value)
	}
}

func TestEdgeCases(t *testing.T) {
	zeroInt := fp.Some(0)
	if !zeroInt.IsSome() {
		t.Error("Some(0) should have IsSome() == true")
	}

	zeroString := fp.Some("")
	if !zeroString.IsSome() {
		t.Error("Some(\"\") should have IsSome() == true")
	}

	var nilPtr *string = nil
	someNil := fp.Some(nilPtr)
	if !someNil.IsSome() {
		t.Error("Some(nil) should have IsSome() == true")
	}

	nestedopt := fp.Some(fp.Some(64))
	extracted := fp.Map(nestedopt, func(opt fp.Opt[int]) int {
		value, _ := opt.Get()
		return *value
	})
	value, _ := extracted.Get()
	if *value != 64 {
		t.Errorf("Extracting from nested opt should work, got %v", *value)
	}
}

func TestOptJsonMarshaling(t *testing.T) {
    type WithPointer struct {
        Value *string `json:"value,omitempty"`
    }

    type WithOpt struct {
        Value fp.Opt[string] `json:"value,omitempty,omitzero"`
    }

    testCases := []struct {
        name     string
        ptrValue *string
        optValue fp.Opt[string]
    }{
        {
            name:     "nil value",
            ptrValue: nil,
            optValue: fp.None[string](),
        },
        {
            name:     "empty string",
            ptrValue: func() *string { s := ""; return &s }(),
            optValue: fp.Some(""),
        },
        {
            name:     "non-empty string",
            ptrValue: func() *string { s := "hello"; return &s }(),
            optValue: fp.Some("hello"),
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            withPtr := WithPointer{Value: tc.ptrValue}
            withOpt := WithOpt{Value: tc.optValue}

            ptrJSON, err := json.Marshal(withPtr)
            if err != nil {
                t.Fatalf("Failed to marshal pointer: %v", err)
            }

            optJSON, err := json.Marshal(withOpt)
            if err != nil {
                t.Fatalf("Failed to marshal Opt: %v", err)
            }

            if string(ptrJSON) != string(optJSON) {
                t.Errorf("JSON output differs: pointer=%s, opt=%s", ptrJSON, optJSON)
            }

            var newWithPtr WithPointer
            var newWithOpt WithOpt

            if err := json.Unmarshal(ptrJSON, &newWithPtr); err != nil {
                t.Fatalf("Failed to unmarshal to pointer: %v", err)
            }

            if err := json.Unmarshal(optJSON, &newWithOpt); err != nil {
                t.Fatalf("Failed to unmarshal to Opt: %v", err)
            }

            if (newWithPtr.Value == nil) != newWithOpt.Value.IsNone() {
                t.Errorf("Presence differs after unmarshaling: pointer=%v, opt=%v", 
                    newWithPtr.Value != nil, newWithOpt.Value.IsSome())
            }

            if newWithPtr.Value != nil && newWithOpt.Value.IsSome() {
                ptrVal := *newWithPtr.Value
                optVal, _ := newWithOpt.Value.Get()

                if ptrVal != *optVal {
                    t.Errorf("Values differ after unmarshaling: pointer=%v, opt=%v", 
                        ptrVal, *optVal)
                }
            }
        })
    }
}

func TestOptJsonRoundtrip(t *testing.T) {
    type ComplexStruct struct {
        StringValue fp.Opt[string] `json:"string,omitempty"`
        IntValue    fp.Opt[int]    `json:"int,omitempty"`
        BoolValue   fp.Opt[bool]   `json:"bool,omitempty"`
    }

    original := ComplexStruct{
        StringValue: fp.Some("test"),
        IntValue:    fp.Some(42),
        BoolValue:   fp.None[bool](),
    }

    jsonData, err := json.Marshal(original)
    if err != nil {
        t.Fatalf("Failed to marshal: %v", err)
    }

    var result ComplexStruct
    if err := json.Unmarshal(jsonData, &result); err != nil {
        t.Fatalf("Failed to unmarshal: %v", err)
    }

    // string
    if original.StringValue.IsSome() != result.StringValue.IsSome() {
        t.Errorf("StringValue presence differs after roundtrip")
    }
    if original.StringValue.IsSome() {
        origVal, _ := original.StringValue.Get()
        resultVal, _ := result.StringValue.Get()
        if *origVal != *resultVal {
            t.Errorf("StringValue differs after roundtrip: original=%v, result=%v", 
                *origVal, *resultVal)
        }
    }

    // int
    if original.IntValue.IsSome() != result.IntValue.IsSome() {
        t.Errorf("IntValue presence differs after roundtrip")
    }
    if original.IntValue.IsSome() {
        origVal, _ := original.IntValue.Get()
        resultVal, _ := result.IntValue.Get()
        if *origVal != *resultVal {
            t.Errorf("IntValue differs after roundtrip: original=%v, result=%v", 
                *origVal, *resultVal)
        }
    }

	// bool
    if original.BoolValue.IsSome() != result.BoolValue.IsSome() {
        t.Errorf("BoolValue presence differs after roundtrip")
    }
}

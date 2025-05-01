package fp_test

import (
	"encoding/json"
	"fmt"

	"github.com/kanopy-platform/go-library/fp"
)

func Example() {
	// Basic package usage
	someValue := fp.Some(42)
	noValue := fp.None[int]()

	fmt.Println("someValue contains a value:", someValue.IsSome())
	fmt.Println("noValue contains a value:", noValue.IsSome())

	// Output:
	// someValue contains a value: true
	// noValue contains a value: false
}

func ExampleSome() {
	value := fp.Some("hello")
	fmt.Println(value.IsSome())
	// Output: true
}

func ExampleNone() {
	empty := fp.None[int]()
	fmt.Println(empty.IsSome())
	// Output: false
}

func ExampleOpt_Get() {
	// Getting a value that exists
	someValue := fp.Some("hello")
	value, err := someValue.Get()
	if err == nil {
		fmt.Println("Value:", *value)
	}

	// Getting a value that doesn't exist
	noValue := fp.None[string]()
	_, err = noValue.Get()
	fmt.Println("Error:", err)

	// Output:
	// Value: hello
	// Error: opt has no value
}

func ExampleOpt_Or() {
	empty := fp.None[string]()
	withDefault := empty.Or(fp.Some("default"))

	value, _ := withDefault.Get()
	fmt.Println(*value)

	// Output: default
}

func ExampleOpt_Filter() {
	isEven := func(n int) bool { return n%2 == 0 }

	evenNumber := fp.Some(4)
	oddNumber := fp.Some(5)

	fmt.Println(evenNumber.Filter(isEven).IsSome())
	fmt.Println(oddNumber.Filter(isEven).IsSome())

	// Output:
	// true
	// false
}

func ExampleMapOpt() {
	numberOpt := fp.Some(5)
	stringOpt := fp.MapOpt(numberOpt, func(x int) string {
		return fmt.Sprintf("Number is %d", x)
	})

	str, _ := stringOpt.Get()
	fmt.Println(*str)

	// Output: Number is 5
}

func ExampleOpt_MarshalJSON() {
    someString := fp.Some("hello")
    someInt := fp.Some(42)
    noneValue := fp.None[string]()

    stringJSON, _ := json.Marshal(someString)
    intJSON, _ := json.Marshal(someInt)
    noneJSON, _ := json.Marshal(noneValue)

    fmt.Println(string(stringJSON))
    fmt.Println(string(intJSON))
    fmt.Println(string(noneJSON))

    // Output:
    // "hello"
    // 42
    // null
}

func ExampleOpt_UnmarshalJSON() {
	stringJSON := []byte(`"hello"`)
	intJSON := []byte(`42`)
	nullJSON := []byte(`null`)

	var stringOpt fp.Opt[string]
	var intOpt fp.Opt[int]
	var nullOpt fp.Opt[bool]

	if err := json.Unmarshal(stringJSON, &stringOpt); err != nil {
		fmt.Println("Error unmarshalling string:", err)
		return
	}
	if err := json.Unmarshal(intJSON, &intOpt); err != nil {
		fmt.Println("Error unmarshalling int:", err)
		return
	}
	if err := json.Unmarshal(nullJSON, &nullOpt); err != nil {
		fmt.Println("Error unmarshalling null:", err)
		return
	}

	fmt.Println("String is Some:", stringOpt.IsSome())
	fmt.Println("Int is Some:", intOpt.IsSome())
	fmt.Println("Null is Some:", nullOpt.IsSome())

	stringVal, _ := stringOpt.Get()
	intVal, _ := intOpt.Get()

	fmt.Println("String value:", *stringVal)
	fmt.Println("Int value:", *intVal)

	// Output:
	// String is Some: true
	// Int is Some: true
	// Null is Some: false
	// String value: hello
	// Int value: 42
}

func ExampleOpt_json_omit() {
	type Person struct {
		Name    string         `json:"name"`
		Address fp.Opt[string] `json:"address,omitzero"`
	}

	complete := Person{
		Name:    "Alice",
		Address: fp.Some("123 Main St"),
	}
	noAddress := Person{
		Name:    "Bob",
		Address: fp.None[string](),
	}

	completeJSON, _ := json.Marshal(complete)
	noAddressJSON, _ := json.Marshal(noAddress)

	fmt.Println("Complete:", string(completeJSON))
	fmt.Println("No Address:", string(noAddressJSON))

	// Output:
	// Complete: {"name":"Alice","address":"123 Main St"}
	// No Address: {"name":"Bob"}
}

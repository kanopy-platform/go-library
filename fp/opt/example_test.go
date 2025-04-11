package opt_test

import (
	"fmt"

	"github.com/kanopy-platform/go-library/fp/opt"
)

func Example() {
	// Basic package usage
	someValue := opt.Some(42)
	noValue := opt.None[int]()

	fmt.Println("someValue contains a value:", someValue.IsSome())
	fmt.Println("noValue contains a value:", noValue.IsSome())

	// Output:
	// someValue contains a value: true
	// noValue contains a value: false
}

func ExampleSome() {
	value := opt.Some("hello")
	fmt.Println(value.IsSome())
	// Output: true
}

func ExampleNone() {
	empty := opt.None[int]()
	fmt.Println(empty.IsSome())
	// Output: false
}

func ExampleOpt_Get() {
	// Getting a value that exists
	someValue := opt.Some("hello")
	value, err := someValue.Get()
	if err == nil {
		fmt.Println("Value:", *value)
	}

	// Getting a value that doesn't exist
	noValue := opt.None[string]()
	_, err = noValue.Get()
	fmt.Println("Error:", err)

	// Output:
	// Value: hello
	// Error: opt has no value
}

func ExampleOpt_Or() {
	empty := opt.None[string]()
	withDefault := empty.Or(opt.Some("default"))

	value, _ := withDefault.Get()
	fmt.Println(*value)

	// Output: default
}

func ExampleOpt_Filter() {
	isEven := func(n int) bool { return n%2 == 0 }

	evenNumber := opt.Some(4)
	oddNumber := opt.Some(5)

	fmt.Println(evenNumber.Filter(isEven).IsSome())
	fmt.Println(oddNumber.Filter(isEven).IsSome())

	// Output:
	// true
	// false
}

func ExampleMap() {
	numberOpt := opt.Some(5)
	stringOpt := opt.Map(numberOpt, func(x int) string {
		return fmt.Sprintf("Number is %d", x)
	})

	str, _ := stringOpt.Get()
	fmt.Println(*str)

	// Output: Number is 5
}

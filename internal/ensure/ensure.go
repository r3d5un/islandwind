// Package ensure contains functions for runtime assertions.
package ensure

// True asserts that a given condition is true or panics. A given message is printed as part of the
// panic function call.
func True(condition bool, message string) {
	if condition != true {
		panic(messageHelper("true boolean assertion failed", message))
	}
}

// False asserts that a given condition is true or panics. A given message is printed as part of the
// panic function call.
func False(condition bool, message string) {
	if condition == true {
		panic(messageHelper("false boolean assertion failed", message))
	}
}

// Nil asserts that a given object is nil or panics. The message is used as part of the panic
// function call.
func Nil(obj any, message string) {
	if obj != nil {
		panic(messageHelper("nil assertion failed", message))
	}
}

// NotNil asserts that a given object is not nil or panics. The message is used as part of the panic
// function call.
func NotNil(obj any, message string) {
	if obj == nil {
		panic(messageHelper("not nil assertion failed", message))
	}
}

// Equal asserts that two given objects are equal or panics. A given message is printed as part of
// the panic function call.
func Equal[T comparable](expected, actual T, message string) {
	if expected != actual {
		panic(messageHelper("equality assertion failed", message))
	}
}

// NotEqual asserts that two given objects are equal or panics. A given message is printed as part
// of the panic function call.
func NotEqual[T comparable](expected, actual T, message string) {
	if expected == actual {
		panic(messageHelper("non-equality assertion failed", message))
	}
}

func Error(err error, message string) {
	if err == nil {
		panic(messageHelper("error assertion failed", message))
	}
}

func NoError(err error, message string) {
	if err != nil {
		panic(messageHelper("non-error assertion failed", message))
	}
}

func messageHelper(assertion string, message string) string {
	return assertion + ": " + message
}

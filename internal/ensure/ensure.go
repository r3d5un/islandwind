// Package ensure contains functions for runtime assertions.
package ensure

// True asserts that a given condition is true or panics. A given message is printed as part of the
// panic function call.
func True(condition bool, message string) {
	if condition != true {
		panic(message)
	}
}

// False asserts that a given condition is true or panics. A given message is printed as part of the
// panic function call.
func False(condition bool, message string) {
	if condition == true {
		panic(message)
	}
}

// Nil asserts that a given object is nil or panics. The message is used as part of the panic
// function call.
func Nil(obj any, message string) {
	if obj != nil {
		panic(message)
	}
}

// NotNil asserts that a given object is not nil or panics. The message is used as part of the panic
// function call.
func NotNil(obj any, message string) {
	if obj == nil {
		panic(message)
	}
}

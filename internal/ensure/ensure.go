// Package ensure contains functions for runtime assertions.
package ensure

import (
	"runtime"
	"strconv"
)

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

// Index asserts that a given index is within the bounds of a given length or panics. A given
// message is printed as part of the panic function call.
//
// Example:
//
//	ensure.Index(i, len(slice), "index out of bounds")
func Index(i, len int, message string) {
	if i < 0 || i >= len {
		panic(messageHelper(fmt.Sprintf("index out of bounds: index %d, len %d", i, len), message))
	}
}

// messageHelper returns a formatted string indicating the file, line, and assertion details. It is
// used to format the panic message for runtime assertions.
func messageHelper(assertion string, message string) string {
	// runtime.Caller is used to report the file and line number of a function invocation. The
	// argument is the number of stack frames to skip. In this case the messageHelper and the
	// caller assertion function are skipped, which prints the location where the runtime assertion
	// failed to be given.
	_, file, line, _ := runtime.Caller(2)
	return file + ":" + strconv.Itoa(line) + ": " + assertion + ": " + message
}

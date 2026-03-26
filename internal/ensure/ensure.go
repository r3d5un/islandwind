// Package ensure contains functions for runtime assertions.
package ensure

import (
	"fmt"
	"reflect"
	"runtime"
	"slices"
	"strconv"
)

// True asserts that a given condition is true or panics. A given message is printed as part of the
// panic function call.
//
// Example:
//
//	ensure.True(len(slice) > 0, "slice must not be empty")
func True(condition bool, message string) {
	if condition != true {
		panic(messageHelper("true boolean assertion failed", message))
	}
}

// False asserts that a given condition is true or panics. A given message is printed as part of the
// panic function call.
//
// Example:
//
//	ensure.False(user.IsDeleted, "user must not be deleted")
func False(condition bool, message string) {
	if condition == true {
		panic(messageHelper("false boolean assertion failed", message))
	}
}

// Nil asserts that a given object is nil or panics. The message is used as part of the panic
// function call.
//
// Example:
//
//	ensure.Nil(err, "error should be nil")
func Nil(obj any, message string) {
	if obj != nil {
		panic(messageHelper("nil assertion failed", message))
	}
}

// NotNil asserts that a given object is not nil or panics. The message is used as part of the panic
// function call.
//
// Example:
//
//	ensure.NotNil(client, "client should not be nil")
func NotNil(obj any, message string) {
	if obj == nil {
		panic(messageHelper("not nil assertion failed", message))
	}
}

// Equal asserts that two given objects are equal or panics. A given message is printed as part of
// the panic function call.
//
// Example:
//
//	ensure.Equal(200, res.StatusCode, "status code must be 200")
func Equal[T comparable](expected, actual T, message string) {
	if expected != actual {
		panic(messageHelper("equality assertion failed", message))
	}
}

// NotEqual asserts that two given objects are equal or panics. A given message is printed as part
// of the panic function call.
//
// Example:
//
//	ensure.NotEqual(0, userID, "user ID must not be zero")
func NotEqual[T comparable](expected, actual T, message string) {
	if expected == actual {
		panic(messageHelper("non-equality assertion failed", message))
	}
}

// Error asserts that a given error is not nil or panics. A given message is printed as part of
// the panic function call.
//
// Example:
//
//	err := doSomething()
//	ensure.Error(err, "doSomething should return an error")
func Error(err error, message string) {
	if err == nil {
		panic(messageHelper("error assertion failed", message))
	}
}

// NoError asserts that a given error is nil or panics. A given message is printed as part of the
// panic function call.
//
// Example:
//
//	err := doSomething()
//	ensure.NoError(err, "doSomething should not return an error")
func NoError(err error, message string) {
	if err != nil {
		panic(messageHelper("non-error assertion failed", message))
	}
}

// For asserts that a given condition is true for all elements in a slice or panics. A given
// message is printed as part of the panic function call.
//
// Example:
//
//	slice := []int{1, 2, 3}
//	ensure.For(slice, func(v int) bool { return v > 0 }, "slice must contain only positive integers")
//
// Example:
//
//	slice := []bool{false, true, true, false, true}
//	ensure.For(slice, confirm.False, "all boolean values must be false")
func For[T any](slice []T, fn func(T) bool, message string) {
	for _, v := range slice {
		if !fn(v) {
			panic(messageHelper("for assertion failed", message))
		}
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

// Valid asserts that a given object is not nil and, if it is a pointer, that it does not point to
// nil.
//
// Checking pointer values relies on reflection, which carries a heavy runtime penalty. It is
// recommended to use this value once when injecting or initializing objects and avoid using
// Valid in hot loops or functions.
//
// Example:
//
//	var p *int
//	ensure.Valid(p, "p must not be nil") // panics
//
// Example:
//
//	ensure.Valid(client, "client must be valid: %s", clientID)
func Valid(obj any, format string, args ...any) {
	if obj == nil || (reflect.ValueOf(obj).Kind() == reflect.Ptr && reflect.ValueOf(obj).IsNil()) {
		panic(messageHelper("validity assertion failed", fmt.Sprintf(format, args...)))
	}
}

// Contains asserts that a given slice contains a given element or panics. A given message is
// printed as part of the panic function call.
//
// Example:
//
//	slice := []int{1, 2, 3}
//	ensure.Contains(slice, 2, "slice must contain 2")
func Contains[T comparable](elements []T, element T, message string) {
	if !slices.Contains(elements, element) {
		panic(messageHelper("contains assertion failed", message))
	}
}

// NotContains asserts that a given slice does not contain a given element or panics. A given
// message is printed as part of the panic function call.
//
// Example:
//
//	slice := []int{1, 2, 3}
//	ensure.NotContains(slice, 4, "slice must not contain 4")
func NotContains[T comparable](elements []T, element T, message string) {
	if slices.Contains(elements, element) {
		panic(messageHelper("not-contains assertion failed", message))
	}
}

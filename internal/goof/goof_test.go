package goof_test

import (
	"errors"
	"testing"

	"github.com/r3d5un/islandwind/internal/goof"
	"github.com/stretchr/testify/assert"
)

// TestNew tests the public methods associated with the goof.New function.
func TestNew(t *testing.T) {
	msg := "test error"
	err := goof.New(msg)

	assert.NotNil(t, err)
	assert.Equal(t, msg, err.Error())

	var goofErr goof.Error
	assert.True(t, errors.As(err, &goofErr))
	assert.Equal(t, msg, goofErr.Message())
}

// TestWrap tests the public methods associated with the goof.Wrap function.
func TestWrap(t *testing.T) {
	inner := errors.New("inner error")
	err := goof.Wrap(inner)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, inner)

	var goofErr goof.Error
	assert.True(t, errors.As(err, &goofErr))
	assert.Equal(t, inner, goofErr.Internal())
}

// TestWith tests the public methods associated with the goof.With function.
func TestWith(t *testing.T) {
	t.Run("SingleWith", func(t *testing.T) {
		err := goof.With("key", "value").New("test error")

		var goofErr goof.Error
		assert.True(t, errors.As(err, &goofErr))
		assert.Equal(t, "value", goofErr.Metadata()["key"])
	})

	t.Run("MultipleWithCalls", func(t *testing.T) {
		err := goof.
			With("key1", "val1").
			With("key2", "val2").
			New("test error 2")
		var goofErr goof.Error
		assert.True(t, errors.As(err, &goofErr))
		assert.Equal(t, "val1", goofErr.Metadata()["key1"])
		assert.Equal(t, "val2", goofErr.Metadata()["key2"])
	})

	t.Run("BuilderImmutability", func(t *testing.T) {
		builder := goof.With("base", "value")
		b1 := builder.With("a", "1")
		b2 := builder.With("b", "2")

		errA := b1.New("err a")
		errB := b2.New("err b")

		var goofErrA, goofErrB goof.Error
		errors.As(errA, &goofErrA)
		errors.As(errB, &goofErrB)

		assert.Contains(t, goofErrA.Metadata(), "a")
		assert.NotContains(t, goofErrA.Metadata(), "b")
		assert.Contains(t, goofErrB.Metadata(), "b")
		assert.NotContains(t, goofErrB.Metadata(), "a")
	})
}

// TestCode tests the public methods associated with the goof.Code type.
func TestCode(t *testing.T) {
	code := "ERR_CODE"
	err := goof.Code(code).New("test error")

	var goofErr goof.Error
	assert.True(t, errors.As(err, &goofErr))
	assert.Equal(t, code, goofErr.Code())
}

// TestService tests the public methods associated with the goof.Service type.
func TestService(t *testing.T) {
	svc := "test-service"
	err := goof.Service(svc).New("test error")

	var goofErr goof.Error
	assert.True(t, errors.As(err, &goofErr))
	assert.Equal(t, svc, goofErr.Service())
}

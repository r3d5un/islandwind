package goof_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/r3d5un/islandwind/internal/goof"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	originalErr := fmt.Errorf("test error")
	err := goof.
		New(originalErr)

	assert.NotNil(t, err)

	goofErr, ok := errors.AsType[*goof.Error](err)
	assert.True(t, ok)
	assert.Equal(t, originalErr.Error(), goofErr.Error())
}

func TestWrap(t *testing.T) {
	inner := errors.New("inner error")
	err := goof.Wrap(inner)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, inner)

	goofErr, ok := errors.AsType[*goof.Error](err)
	assert.True(t, ok)
	assert.Equal(t, inner, goofErr.Internal())
}

func TestWith(t *testing.T) {
	t.Run("SingleWith", func(t *testing.T) {
		err := goof.
			With("key", "value").
			New(fmt.Errorf("test error"))

		goofErr, ok := errors.AsType[*goof.Error](err)
		assert.True(t, ok)
		assert.Equal(t, "value", goofErr.Metadata()["key"])
	})

	t.Run("MultipleWithCalls", func(t *testing.T) {
		err := goof.
			With("key1", "val1").
			With("key2", "val2").
			New(fmt.Errorf("test error"))

		goofErr, ok := errors.AsType[*goof.Error](err)
		assert.True(t, ok)
		assert.Equal(t, "val1", goofErr.Metadata()["key1"])
		assert.Equal(t, "val2", goofErr.Metadata()["key2"])
	})

	t.Run("BuilderImmutability", func(t *testing.T) {
		builder := goof.With("base", "value")
		b1 := builder.With("a", "1")
		b2 := builder.With("b", "2")

		errA := b1.New(fmt.Errorf("err a"))
		errB := b2.New(fmt.Errorf("err b"))

		goofErrA, okA := errors.AsType[*goof.Error](errA)
		goofErrB, okB := errors.AsType[*goof.Error](errB)

		assert.True(t, okA)
		assert.True(t, okB)
		assert.Contains(t, goofErrA.Metadata(), "a")
		assert.NotContains(t, goofErrA.Metadata(), "b")
		assert.Contains(t, goofErrB.Metadata(), "b")
		assert.NotContains(t, goofErrB.Metadata(), "a")
	})
}

func TestCode(t *testing.T) {
	code := "ERR_CODE"
	inner := fmt.Errorf("test error")
	err := goof.
		Code(code).
		New(inner)

	goofErr, ok := errors.AsType[*goof.Error](err)
	assert.True(t, ok)
	assert.Equal(t, inner.Error(), goofErr.Error())
	assert.Equal(t, code, goofErr.Code())
}

func TestService(t *testing.T) {
	svc := "test-service"
	err := goof.
		Service(svc).
		New(fmt.Errorf("test error"))

	goofErr, ok := errors.AsType[*goof.Error](err)
	assert.True(t, ok)
	assert.Equal(t, svc, goofErr.Service())
}

func TestMessage(t *testing.T) {
	msg := "test error"
	err := goof.
		Message(msg).New(fmt.Errorf("test error"))

	goofErr, ok := errors.AsType[*goof.Error](err)
	assert.True(t, ok)
	assert.Equal(t, msg, goofErr.Message())
}

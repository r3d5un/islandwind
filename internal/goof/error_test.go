package goof_test

import (
	"errors"
	"testing"

	"github.com/r3d5un/islandwind/internal/goof"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	code := "TEST_CODE"
	svc := "test-service"
	inner := errors.New("inner error")
	metadata := map[string]any{"key": "value"}
	msg := "test error"

	err := goof.Code(code).
		Service(svc).
		With("key", "value").
		Message(msg).
		New(inner)

	goofErr, ok := errors.AsType[*goof.Error](err)
	assert.True(t, ok)

	t.Run("Error", func(t *testing.T) {
		assert.Equal(t, inner.Error(), goofErr.Error())
	})

	t.Run("Code", func(t *testing.T) {
		assert.Equal(t, code, goofErr.Code())
	})

	t.Run("Message", func(t *testing.T) {
		assert.Equal(t, msg, goofErr.Message())
	})

	t.Run("Metadata", func(t *testing.T) {
		assert.Equal(t, metadata, goofErr.Metadata())
	})

	t.Run("Time", func(t *testing.T) {
		assert.NotZero(t, goofErr.Time())
	})

	t.Run("Service", func(t *testing.T) {
		assert.Equal(t, svc, goofErr.Service())
	})
}

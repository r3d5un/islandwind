package goof_test

import (
	"fmt"
	"testing"

	"github.com/r3d5un/islandwind/internal/goof"
	"github.com/stretchr/testify/assert"
)

const (
	code = "TestError"
	msg  = "human readable error"
)

func TestNew(t *testing.T) {
	internal := func(t *testing.T) error {
		return fmt.Errorf("internal error: %s", t.Name())
	}

	t.Run("New", func(t *testing.T) {
		internalErr := internal(t)
		err := goof.New(code, msg, internalErr).With("key", "value")

		assert.Error(t, err)
		assert.Equal(t, code, err.Code)
		assert.Equal(t, msg, err.Error())
		assert.Equal(t, internalErr, err.Internal)
		assert.NotNil(t, err.Metadata)
		assert.Equal(t, "value", err.Metadata["key"])
	})

	t.Run("WithMetadata", func(t *testing.T) {
		err := goof.New(code, msg, nil).
			With("key1", "value1").
			With("key2", 123)

		assert.Equal(t, "value1", err.Metadata["key1"])
		assert.Equal(t, 123, err.Metadata["key2"])
	})

}

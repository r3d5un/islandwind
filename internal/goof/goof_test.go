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
		err := goof.New(code, msg, internalErr, map[string]any{"key": "value"})

		assert.Error(t, err)
		assert.Equal(t, code, err.Code)
		assert.Equal(t, msg, err.Error())
		assert.Equal(t, internalErr, err.Internal)
		assert.NotNil(t, err.Metadata)
	})

	t.Run("DefensiveMetadataCopy", func(t *testing.T) {
		metadata := map[string]any{"key1": "value1"}
		err := goof.New(code, msg, internal(t), metadata)

		metadata["key1"] = "mutated"
		metadata["key2"] = "added"

		assert.Equal(
			t,
			"value1", err.Metadata["key1"],
			"error metadata should not be mutable",
		)

		_, exists := err.Metadata["key2"]
		assert.False(
			t,
			exists,
			"new metadata entries should not be added externally",
		)
	})
}

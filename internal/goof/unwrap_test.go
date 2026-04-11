package goof_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/r3d5un/islandwind/internal/goof"
	"github.com/stretchr/testify/assert"
)

func TestUnwrap(t *testing.T) {
	t.Run("BasicUnwrap", func(t *testing.T) {
		inner := errors.New("inner error")
		err := goof.New(inner)

		unwrapped := errors.Unwrap(err)
		assert.Equal(t, inner, unwrapped)
	})

	t.Run("NormalErrors", func(t *testing.T) {
		baseErr := errors.New("base error")
		wrappedOnce := fmt.Errorf("wrapped once: %w", baseErr)
		goofErr := goof.Wrap(wrappedOnce)
		wrappedTwice := fmt.Errorf("wrapped twice: %w", goofErr)

		assert.ErrorIs(t, wrappedTwice, baseErr)
		assert.ErrorIs(t, wrappedTwice, wrappedOnce)
		assert.ErrorIs(t, wrappedTwice, goofErr)

		unwrappableErr, ok := errors.AsType[*goof.Error](wrappedTwice)
		assert.True(t, ok)
		assert.Equal(t, wrappedOnce, unwrappableErr.Unwrap())
	})

	t.Run("MultipleGoofLevels", func(t *testing.T) {
		rootErr := errors.New("root")
		levelOneErr := goof.Code("LEVEL_1").New(rootErr)
		levelTwoErr := goof.Code("LEVEL_2").New(levelOneErr)

		assert.ErrorIs(t, levelTwoErr, rootErr)
		assert.ErrorIs(t, levelTwoErr, levelOneErr)

		unwrappedLevelTwoErr := errors.Unwrap(levelTwoErr)
		assert.Equal(t, levelOneErr, unwrappedLevelTwoErr)

		unwrappedLevelOneErr := errors.Unwrap(unwrappedLevelTwoErr)
		assert.Equal(t, rootErr, unwrappedLevelOneErr)
	})
}

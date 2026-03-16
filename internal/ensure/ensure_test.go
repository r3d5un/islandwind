package ensure_test

import (
	"testing"

	"github.com/r3d5un/islandwind/internal/ensure"
	"github.com/stretchr/testify/assert"
)

func TestTrue(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert.NotPanics(
			t,
			func() {
				ensure.True(true, "should never panic")
			},
		)
	})

	t.Run("Panic", func(t *testing.T) {
		assert.Panics(
			t,
			func() {
				ensure.True(false, "should always panic")
			},
		)
	})
}

func TestFalse(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert.NotPanics(
			t,
			func() {
				ensure.False(false, "should never panic")
			},
		)
	})

	t.Run("Panic", func(t *testing.T) {
		assert.Panics(
			t,
			func() {
				ensure.False(true, "should always panic")
			},
		)
	})
}

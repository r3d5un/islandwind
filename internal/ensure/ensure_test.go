package ensure_test

import (
	"errors"
	"testing"

	"github.com/r3d5un/islandwind/internal/confirm"
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

func TestNil(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert.NotPanics(
			t,
			func() {
				ensure.Nil(nil, "should never panic")
			},
		)
	})

	t.Run("Panic", func(t *testing.T) {
		assert.Panics(
			t,
			func() {
				ensure.Nil(&struct{}{}, "should always panic")
			},
		)
	})
}

func TestNotNil(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert.NotPanics(
			t,
			func() {
				ensure.NotNil(&struct{}{}, "should never panic")
			},
		)
	})

	t.Run("Panic", func(t *testing.T) {
		assert.Panics(
			t,
			func() {
				ensure.NotNil(nil, "should always panic")
			},
		)
	})
}

func TestEqual(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert.NotPanics(
			t,
			func() {
				ensure.Equal(1, 1, "should never panic")
			},
		)
	})

	t.Run("Panic", func(t *testing.T) {
		assert.Panics(
			t,
			func() {
				ensure.Equal(0, 1, "should always panic")
			},
		)
	})
}

func TestNotEqual(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert.NotPanics(
			t,
			func() {
				ensure.NotEqual(0, 1, "should never panic")
			},
		)
	})

	t.Run("Panic", func(t *testing.T) {
		assert.Panics(
			t,
			func() {
				ensure.NotEqual(1, 1, "should always panic")
			},
		)
	})
}

func TestError(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert.NotPanics(t, func() {
			ensure.Error(errors.New("test error"), "should never panic")
		})
	})

	t.Run("Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			ensure.Error(nil, "should always panic")
		})
	})
}

func TestNoError(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert.NotPanics(t, func() {
			ensure.NoError(nil, "should never panic")
		})
	})

	t.Run("Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			ensure.NoError(errors.New("test error"), "should always panic")
		})
	})
}

func TestIndexs(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert.NotPanics(t, func() {
			ensure.Index(0, 10, "should never panic")
		})
	})

	t.Run("Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			ensure.Index(100, 10, "should always panic")
		})
	})
}

func TestFor(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert.NotPanics(t, func() {
			ensure.For([]int{1, 2, 3}, func(v int) bool {
				return v > 0
			}, "should never panic")
		})
	})

	t.Run("Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			ensure.For([]int{1, -1, 3}, func(v int) bool {
				return v > 0
			}, "should panic on -1")
		})
	})

	t.Run("Empty", func(t *testing.T) {
		assert.NotPanics(t, func() {
			ensure.For([]int{}, func(v int) bool {
				return v > 0
			}, "should never panic on empty slice")
		})
	})

	t.Run("WithConfirm", func(t *testing.T) {
		slice := []bool{true, true, true, true, true}
		assert.Panics(t, func() {
			ensure.For(slice, confirm.False, "all boolean values must be false")
		})
	})
}

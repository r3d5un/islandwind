package goof_test

import (
	"fmt"
	"testing"

	"github.com/r3d5un/islandwind/internal/goof"
	"github.com/stretchr/testify/assert"
)

func TestErrorNew(t *testing.T) {
	code := "TestError"
	msg := "human readable error"
	internal := fmt.Errorf("internal error: %s", t.Name())

	err := goof.New(code, msg, internal, map[string]any{"key": "value"})

	assert.Error(t, err)
	assert.Equal(t, code, err.Code)
	assert.Equal(t, msg, err.Error())
	assert.Equal(t, internal, err.Internal)
	assert.NotNil(t, err.Metadata)
}

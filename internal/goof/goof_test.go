package goof_test

import (
	"testing"

	"github.com/r3d5un/islandwind/internal/goof"
)

const (
	code = "TestError"
	msg  = "human readable error"
)

func TestNew(t *testing.T) {
	goof.New(t.Name())

	err := goof.New(t.Name())
	goof.With("test", "test").Wrap(err)
	goof.With("test1", "test1").With("test2", "test2").With("test3", "test3")
}

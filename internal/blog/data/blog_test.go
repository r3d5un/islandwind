package data_test

import (
	"context"
	"testing"
	"time"
)

func TestBlogModel(t *testing.T) {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	t.Cleanup(func() {
		t.Logf("cleaning up test: %s", t.Name())
		defer cancel()
	})

	t.Run("Insert", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("Select", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("SelectMany", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("Update", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("Delete", func(t *testing.T) {
		t.Skip("not implemented")
	})
}

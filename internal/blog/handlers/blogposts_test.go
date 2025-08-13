package handlers_test

import (
	"context"
	"testing"
	"time"
)

func TestBlogHandlers(t *testing.T) {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	t.Cleanup(func() {
		t.Logf("cleaning up test: %s", t.Name())
		defer cancel()
	})

	t.Run("POST", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("GET", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("GETList", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("PATCH", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("DELETE", func(t *testing.T) {
		t.Skip("not implemented")
	})
}

package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPostgresCache(t *testing.T) {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	resourceID := uuid.New()
	data := testData{Field: t.Name()}

	t.Cleanup(func() {
		defer cancel()
	})

	t.Run("StartShutdown", func(t *testing.T) {
		assert.NoError(t, postgresCache.Start())
		t.Cleanup(func() {
			postgresCache.Shutdown()
		})
	})

	t.Run("SetGetDelete", func(t *testing.T) {
		assert.NoError(t, postgresCache.Start())
		t.Cleanup(func() {
			postgresCache.Shutdown()
		})

		postgresCache.Set(resourceID, data)

		var read testData
		assert.Eventually(
			t,
			func() bool {
				err := postgresCache.Get(resourceID, &read)
				if err != nil {
					return false
				}
				return true

			},
			time.Millisecond*1000,
			time.Millisecond*10,
		)
		assert.Equal(t, data.Field, read.Field)

		t.Cleanup(func() {
			assert.NoError(t, postgresCache.Delete(resourceID))
		})
	})
}

type testData struct {
	Field string `json:"field"`
}

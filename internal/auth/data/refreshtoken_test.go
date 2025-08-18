package data_test

import (
	"context"
	"testing"
	"time"

	"github.com/r3d5un/islandwind/internal/auth/data"
	"github.com/stretchr/testify/assert"
)

func TestRefreshTokenModel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	t.Cleanup(func() {
		t.Logf("cleaning up test: %s", t.Name())
		defer cancel()
	})

	var refreshToken data.RefreshToken

	t.Run("Insert", func(t *testing.T) {
		timestamp := time.Now()
		inserted, err := models.RefreshTokens.Insert(ctx, data.RefreshTokenInput{
			Issuer:     "islandwind",
			Audience:   "islandwind",
			Expiration: timestamp,
			IssuedAt:   timestamp,
			NotBefore:  timestamp,
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, inserted)

		refreshToken = *inserted
	})

	t.Run("Select", func(t *testing.T) {
		selected, err := models.RefreshTokens.SelectOne(ctx, refreshToken.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, selected)
		assert.Equal(t, refreshToken, *selected)
	})

	t.Run("Delete", func(t *testing.T) {
		deleted, err := models.RefreshTokens.Delete(ctx, refreshToken.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, deleted)
		assert.Equal(t, refreshToken, *deleted)
	})

	t.Run("DeleteMany", func(t *testing.T) {
		timestamp := time.Now()
		rowsAffected, err := models.RefreshTokens.DeleteMany(
			ctx,
			data.Filter{ExpirationTo: &timestamp},
		)
		assert.NoError(t, err)
		assert.NotNil(t, *rowsAffected)
	})
}

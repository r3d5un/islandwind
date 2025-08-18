package data_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
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
			Expiration: timestamp,
			IssuedAt:   timestamp,
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, inserted)

		refreshToken = *inserted
	})

	t.Run("SelectOne", func(t *testing.T) {
		selected, err := models.RefreshTokens.SelectOne(ctx, refreshToken.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, selected)
		assert.Equal(t, refreshToken, *selected)
	})

	t.Run("SelectMany", func(t *testing.T) {
		selected, metadata, err := models.RefreshTokens.SelectMany(ctx, data.Filter{
			PageSize: 1,
			ID:       &refreshToken.ID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, selected)
		assert.NotEmpty(t, selected)
		assert.NotEmpty(t, metadata)
		assert.Equal(t, selected[len(selected)-1].ID, metadata.LastSeen)
	})

	t.Run("Update", func(t *testing.T) {
		timestamp := time.Now()
		inserted, err := models.RefreshTokens.Insert(ctx, data.RefreshTokenInput{
			Issuer:     "islandwind",
			Expiration: timestamp,
			IssuedAt:   timestamp,
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, inserted)

		updated, err := models.RefreshTokens.Update(ctx, data.RefreshTokenPatch{
			ID:            refreshToken.ID,
			Invalidated:   sql.NullBool{Bool: true, Valid: true},
			InvalidatedBy: uuid.NullUUID{Valid: true, UUID: inserted.ID},
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, *updated)

		refreshToken = *updated
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

package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/r3d5un/islandwind/internal/auth/data"
	"github.com/r3d5un/islandwind/internal/auth/repo"
	"github.com/stretchr/testify/assert"
)

func TestTokenRepository(t *testing.T) {
	ctx := context.Background()

	var accessToken string
	var refreshToken string

	t.Run("CreateAccessToken", func(t *testing.T) {
		newJWT, err := authRepo.Tokens.CreateAccessToken()
		assert.NoError(t, err)
		assert.NotEmpty(t, newJWT)

		accessToken = *newJWT
	})

	t.Run("Validate", func(t *testing.T) {
		ok, err := authRepo.Tokens.Validate(ctx, repo.AccessTokenType, accessToken)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("CreateRefreshToken", func(t *testing.T) {
		token, err := authRepo.Tokens.CreateRefreshToken(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		refreshToken = *token
	})

	t.Run("Refresh", func(t *testing.T) {
		a, r, err := authRepo.Tokens.Refresh(ctx, refreshToken)
		assert.NoError(t, err)
		assert.NotEmpty(t, a)
		assert.NotEmpty(t, r)
	})

	t.Run("InvalidateRefreshToken", func(t *testing.T) {
		token, err := authRepo.Tokens.CreateRefreshToken(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		err = authRepo.Tokens.InvalidateRefreshToken(ctx, *token)
		assert.NoError(t, err)
	})

	t.Run("ListRefreshTokens", func(t *testing.T) {
		tokens, metadata, err := authRepo.Tokens.List(ctx, data.Filter{PageSize: 100})
		assert.NoError(t, err)
		assert.NotEmpty(t, metadata)
		assert.NotEmpty(t, tokens)
	})

	t.Run("Delete", func(t *testing.T) {
		timestamp := time.Now().UTC()
		affected, err := authRepo.Tokens.Delete(ctx, data.Filter{ExpirationFrom: &timestamp})
		assert.NoError(t, err)
		assert.NotNil(t, affected)
	})

	t.Run("DeleteExpired", func(t *testing.T) {
		err := authRepo.Tokens.DeleteExpired(ctx)
		assert.NoError(t, err)
	})

}

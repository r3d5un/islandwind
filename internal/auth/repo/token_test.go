package repo_test

import (
	"context"
	"testing"

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
		ok, err := authRepo.Tokens.Validate(ctx, repo.AccessToken, accessToken)
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

		accessToken = *a
		refreshToken = *r
	})

	t.Run("DeleteExpired", func(t *testing.T) {
		err := authRepo.Tokens.DeleteExpired(ctx)
		assert.NoError(t, err)
	})
}

package repo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenRepository(t *testing.T) {
	ctx := context.Background()

	var jwt string

	t.Run("NewJWT", func(t *testing.T) {
		newJWT, err := authRepo.Tokens.NewJWT()
		assert.NoError(t, err)
		assert.NotEmpty(t, newJWT)

		t.Log(*newJWT)

		jwt = *newJWT
	})

	t.Run("Parse", func(t *testing.T) {
		ok, err := authRepo.Tokens.Parse(ctx, jwt)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("NewRefreshToken", func(t *testing.T) {
		refreshToken, err := authRepo.Tokens.NewRefreshToken(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, refreshToken)

		t.Log(*refreshToken)
	})
}

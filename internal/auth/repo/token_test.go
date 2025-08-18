package repo_test

import (
	"context"
	"testing"

	"github.com/r3d5un/islandwind/internal/auth/repo"
	"github.com/stretchr/testify/assert"
)

func TestTokenRepository(t *testing.T) {
	ctx := context.Background()

	repository := repo.NewTokenRepository([]byte("secret"), "islandwind")
	var jwt string

	t.Run("NewJWT", func(t *testing.T) {
		newJWT, err := repository.NewJWT()
		assert.NoError(t, err)
		assert.NotEmpty(t, newJWT)

		t.Log(*newJWT)

		jwt = *newJWT
	})

	t.Run("Parse", func(t *testing.T) {
		ok, err := repository.Parse(ctx, jwt)
		assert.NoError(t, err)
		assert.True(t, ok)
	})
}

package repo

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/r3d5un/islandwind/internal/auth/data"
	"github.com/r3d5un/islandwind/internal/cache"
)

type tokenStore struct {
	models *data.Models
	cache  cache.Cache
}

func newTokenStore(models *data.Models, cache cache.Cache) tokenStore {
	return tokenStore{models: models, cache: cache}
}

func (s *tokenStore) Create(ctx context.Context, issuer string) (*jwt.Token, error) {
	row, err := s.models.RefreshTokens.Insert(ctx, data.RefreshTokenInput{
		Issuer:     issuer,
		Expiration: time.Now().UTC().Add(time.Minute * 60),
		IssuedAt:   time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}
	token := s.newToken(row.ID, row.Expiration, row.IssuedAt, issuer)

	return &token, nil
}

func (r *tokenStore) newToken(
	jti uuid.UUID,
	exp time.Time,
	iat time.Time,
	issuer string,
) jwt.Token {
	return *jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.MapClaims{"jti": jti.String(), "exp": exp.Unix(), "iat": iat.Unix(), "iss": issuer},
	)
}

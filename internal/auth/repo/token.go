package repo

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/r3d5un/islandwind/internal/auth/data"
	"github.com/r3d5un/islandwind/internal/db"
	"github.com/r3d5un/islandwind/internal/logging"
)

var (
	ErrParsingToken    = errors.New("unable to parse token")
	ErrVerifyingToken  = errors.New("unable to parse token")
	ErrTokenExpired    = errors.New("token expired")
	ErrInvalidIssuedAt = errors.New("token iat timestamp invalid")
	ErrIssuerMismatch  = errors.New("token issuer does not match requirements")
	ErrUnauthorized    = errors.New("token unauthorized")
)

type TokenRepository struct {
	signingSecret []byte
	Issuer        string `json:"issuer"`
	models        *data.Models
}

func (r *TokenRepository) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("issuer", r.Issuer),
	)
}

func NewTokenRepository(secret []byte, issuer string, models *data.Models) TokenRepository {
	return TokenRepository{
		signingSecret: secret,
		Issuer:        issuer,
		models:        models,
	}
}

// CreateAccessToken create a new signed JWT token string.
func (r *TokenRepository) CreateAccessToken() (*string, error) {
	token := r.newToken(uuid.New(), time.Now().UTC().Add(time.Minute*5), time.Now().UTC())
	tokenString, err := token.SignedString(r.signingSecret)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (r *TokenRepository) parseToken(input string) (*jwt.Token, error) {
	token, err := jwt.Parse(input, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrParsingToken
		}
		return r.signingSecret, nil
	})
	if err != nil {
		return nil, err
	}

	valid, err := r.verifyClaims(token)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, ErrVerifyingToken
	}

	return token, nil
}

// Validate parses a given input JWT string and validates it claims. An error is returned
// if the token cannot be parsed. If the token is invalid in any way a false boolean
// value is returned along with an error describing the fault.
func (r *TokenRepository) Validate(ctx context.Context, input string) (bool, error) {
	logger := logging.LoggerFromContext(ctx)

	logger.LogAttrs(ctx, slog.LevelInfo, "verifying token")
	_, err := r.parseToken(input)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "error parsing token", slog.String("error", err.Error()),
		)
		return false, err
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "token verified")

	return true, nil
}

func (r *TokenRepository) newRefreshToken(ctx context.Context) (*jwt.Token, error) {
	row, err := r.models.RefreshTokens.Insert(ctx, data.RefreshTokenInput{
		Issuer:     "",
		Expiration: time.Now().UTC().Add(time.Minute * 60),
		IssuedAt:   time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}
	token := r.newToken(row.ID, row.Expiration, row.IssuedAt)

	return &token, nil
}

func (r *TokenRepository) CreateRefreshToken(ctx context.Context) (*string, error) {
	logger := logging.LoggerFromContext(ctx)

	logger.LogAttrs(ctx, slog.LevelInfo, "creating refresh token")
	token, err := r.newRefreshToken(ctx)
	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"unable to create refresh token",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	tokenString, err := token.SignedString(r.signingSecret)
	if err != nil {
		return nil, err
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "refresh token created")

	return &tokenString, nil
}

// Refresh accepts a refresh token string, then produces a new access token and refresh token. The new refresh token
// invalidates and replaces the old refresh token.
func (r *TokenRepository) Refresh(
	ctx context.Context,
	refreshTokenInput string,
) (accessToken *string, refreshToken *string, err error) {
	logger := logging.LoggerFromContext(ctx)

	logger.LogAttrs(ctx, slog.LevelInfo, "validating refresh token")
	token, err := r.parseToken(refreshTokenInput)
	if err != nil {
		return nil, nil, err
	}

	// NOTE: The difference between access tokens and refresh tokens is their duration and that the refresh tokens
	// are stored in the database by its jti.
	//
	// WARN: Any refresh tokens must be checked against the database because there is no difference in payload between
	// an access token and a refresh token. The only way to tell if a token is used as a refresh token is to check
	// if the jti exists in the database and that the token has not been revoked/invalidated.
	id, err := r.jtiFromToken(*token)
	if err != nil {
		return nil, nil, ErrVerifyingToken
	}
	row, err := r.models.RefreshTokens.SelectOne(ctx, *id)
	if err != nil {
		switch {
		case errors.Is(err, db.ErrRecordNotFound):
			return nil, nil, ErrUnauthorized
		default:
			return nil, nil, err
		}
	}
	if row.Invalidated {
		return nil, nil, ErrUnauthorized
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "refresh token validated")

	logger.LogAttrs(ctx, slog.LevelInfo, "replacing refresh token")
	newRefreshToken, err := r.newRefreshToken(ctx)
	if err != nil {
		return nil, nil, err
	}
	id, err = r.jtiFromToken(*newRefreshToken)
	if err != nil {
		return nil, nil, err
	}
	_, err = r.models.RefreshTokens.Update(ctx, data.RefreshTokenPatch{
		ID:            row.ID,
		Invalidated:   sql.NullBool{Valid: true, Bool: true},
		InvalidatedBy: uuid.NullUUID{Valid: true, UUID: *id},
	})
	if err != nil {
		return nil, nil, err
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "refresh token replaced")

	logger.LogAttrs(ctx, slog.LevelInfo, "signing access and refresh tokens")
	newAccessToken := r.newToken(uuid.New(), time.Now().UTC().Add(time.Minute*5), time.Now().UTC())
	access, err := newAccessToken.SignedString(r.signingSecret)
	if err != nil {
		return nil, nil, err
	}
	refresh, err := newRefreshToken.SignedString(r.signingSecret)
	if err != nil {
		return nil, nil, err
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "tokens signed")

	return &access, &refresh, nil
}

func (r *TokenRepository) newToken(jti uuid.UUID, exp time.Time, iat time.Time) jwt.Token {
	return *jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.MapClaims{"jti": jti.String(), "exp": exp.Unix(), "iat": iat.Unix(), "iss": r.Issuer},
	)
}

func (r *TokenRepository) verifyClaims(token *jwt.Token) (bool, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, ErrVerifyingToken
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		return false, ErrVerifyingToken
	}
	if exp.UTC().Unix() < time.Now().UTC().Unix() {
		return false, ErrTokenExpired
	}

	iat, err := claims.GetIssuedAt()
	if err != nil {
		return false, ErrVerifyingToken
	}
	slog.Info("!!!!!", slog.Any("iat", iat.UTC()), slog.Any("time", time.Now().UTC()))
	if iat.UTC().Unix() > time.Now().UTC().Unix() {
		return false, ErrInvalidIssuedAt
	}

	iss, ok := claims["iss"].(string)
	if !ok {
		return false, ErrVerifyingToken
	}
	if iss != r.Issuer {
		return false, ErrIssuerMismatch
	}

	return true, nil
}

func (r *TokenRepository) jtiFromToken(token jwt.Token) (*uuid.UUID, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("unable to get claims from token")
	}
	claim, ok := claims["jti"].(string)
	if !ok {
		return nil, errors.New("no claim present in token")
	}

	id, err := uuid.Parse(claim)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

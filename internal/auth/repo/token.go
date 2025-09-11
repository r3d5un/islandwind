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
	"github.com/r3d5un/islandwind/internal/testsuite"
)

type RefreshToken struct {
	ID            uuid.UUID  `json:"id"`
	Issuer        string     `json:"issuer"`
	Expiration    time.Time  `json:"expiration"`
	IssuedAt      time.Time  `json:"issuedAt"`
	Invalidated   bool       `json:"invalidated"`
	InvalidatedBy *uuid.UUID `json:"invalidatedBy"`
}

func newRefreshTokenFromRow(row *data.RefreshToken) *RefreshToken {
	return &RefreshToken{
		ID:            row.ID,
		Issuer:        row.Issuer,
		Expiration:    row.Expiration,
		IssuedAt:      row.IssuedAt,
		Invalidated:   row.Invalidated,
		InvalidatedBy: db.NullUUIDToPtr(row.InvalidatedBy),
	}
}

type RefreshTokenPatch struct {
	ID            uuid.UUID  `json:"id"`
	Issuer        *string    `json:"issuer"`
	Invalidated   *bool      `json:"invalidated"`
	InvalidatedBy *uuid.UUID `json:"invalidatedBy"`
}

func (t *RefreshTokenPatch) Row() data.RefreshTokenPatch {
	return data.RefreshTokenPatch{
		ID:            t.ID,
		Issuer:        db.NewNullString(t.Issuer),
		Invalidated:   db.NewNullBool(t.Invalidated),
		InvalidatedBy: db.NewNullUUID(t.InvalidatedBy),
	}
}

var (
	ErrParsingToken    = errors.New("unable to parse token")
	ErrVerifyingToken  = errors.New("unable to parse token")
	ErrTokenExpired    = errors.New("token expired")
	ErrInvalidIssuedAt = errors.New("token iat timestamp invalid")
	ErrIssuerMismatch  = errors.New("token issuer does not match requirements")
	ErrUnauthorized    = errors.New("token unauthorized")
)

type TokenType int

const (
	AccessTokenType TokenType = iota
	RefreshTokenType
)

type TokenService interface {
	CreateAccessToken() (accessToken *string, err error)
	Validate(ctx context.Context, tokenType TokenType, input string) (valid bool, err error)
	InvalidateRefreshToken(ctx context.Context, input string) error
	CreateRefreshToken(ctx context.Context) (refreshToken *string, err error)
	Refresh(
		ctx context.Context,
		refreshTokenInput string,
	) (accessToken *string, refreshToken *string, err error)
	Update(ctx context.Context, input RefreshTokenPatch) (*RefreshToken, error)
	DeleteExpired(ctx context.Context) error
	List(ctx context.Context, filter data.Filter) ([]*RefreshToken, *data.Metadata, error)
	Delete(ctx context.Context, filter data.Filter) (int64, error)
}

type TokenRepository struct {
	signingSecret        []byte
	refreshSigningSecret []byte
	Issuer               string `json:"issuer"`
	models               *data.Models
}

func (r *TokenRepository) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("issuer", r.Issuer),
	)
}

func NewTokenRepository(
	accessTokenSecret []byte,
	refreshTokenSecret []byte,
	issuer string,
	models *data.Models,
) TokenService {
	return &TokenRepository{
		signingSecret:        accessTokenSecret,
		refreshSigningSecret: refreshTokenSecret,
		Issuer:               issuer,
		models:               models,
	}
}

func (r *TokenRepository) List(
	ctx context.Context,
	filter data.Filter,
) ([]*RefreshToken, *data.Metadata, error) {
	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"posts",
		slog.Any("filter", filter),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "reading refresh tokens")
	rows, metadata, err := r.models.RefreshTokens.SelectMany(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	testsuite.Assert(rows != nil, "refresh token list is nil", nil)
	testsuite.Assert(metadata != nil, "refresh token metadata are nil", nil)

	tokens := make([]*RefreshToken, metadata.ResponseLength)
	for i, row := range rows {
		tokens[i] = newRefreshTokenFromRow(row)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "refresh tokens retrieved")

	return tokens, metadata, nil
}

func (r *TokenRepository) Update(
	ctx context.Context,
	input RefreshTokenPatch,
) (*RefreshToken, error) {
	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"posts",
		slog.Any("input", input),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "updating refresh token")
	row, err := r.models.RefreshTokens.Update(ctx, input.Row())
	if err != nil {
		return nil, err
	}
	testsuite.Assert(row != nil, "row cannot be nil without an error", nil)
	logger.LogAttrs(ctx, slog.LevelInfo, "refresh token updated")

	return newRefreshTokenFromRow(row), nil
}

func (r *TokenRepository) Delete(ctx context.Context, filter data.Filter) (int64, error) {
	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"posts",
		slog.Any("filter", filter),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "deleting refresh tokens")
	affected, err := r.models.RefreshTokens.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	testsuite.Assert(affected != nil, "affected row count cannot be nil without errors", nil)
	logger.LogAttrs(
		ctx, slog.LevelInfo, "refresh tokens deleted", slog.Int64("affected", *affected),
	)

	return *affected, nil
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

func (r *TokenRepository) parseToken(input string, tokenType TokenType) (*jwt.Token, error) {
	var token *jwt.Token
	var err error

	switch tokenType {
	case RefreshTokenType:
		token, err = jwt.Parse(input, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrParsingToken
			}
			return r.refreshSigningSecret, nil
		})
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				return nil, ErrUnauthorized
			default:
				return nil, err
			}
		}
	case AccessTokenType:
		token, err = jwt.Parse(input, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrParsingToken
			}
			return r.signingSecret, nil
		})
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				return nil, ErrUnauthorized
			default:
				return nil, err
			}
		}
	default:
		return nil, errors.New("unknown token type")
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
func (r *TokenRepository) Validate(
	ctx context.Context,
	tokenType TokenType,
	input string,
) (bool, error) {
	logger := logging.LoggerFromContext(ctx)

	logger.LogAttrs(ctx, slog.LevelInfo, "verifying token")
	_, err := r.parseToken(input, tokenType)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "error parsing token", slog.String("error", err.Error()),
		)
		return false, err
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "token verified")

	return true, nil
}

func (r *TokenRepository) InvalidateRefreshToken(ctx context.Context, input string) error {
	logger := logging.LoggerFromContext(ctx)

	token, err := r.parseToken(input, RefreshTokenType)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "error parsing token", slog.String("error", err.Error()),
		)
		return ErrVerifyingToken
	}

	tokenID, err := r.jtiFromToken(*token)
	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"unable to read refresh token ID",
			slog.String("error", err.Error()),
		)
		return ErrVerifyingToken
	}

	logger = logger.With(slog.String("jti", tokenID.String()))
	logger.LogAttrs(ctx, slog.LevelInfo, "invalidating token")
	_, err = r.models.RefreshTokens.Update(
		ctx,
		data.RefreshTokenPatch{ID: *tokenID, Invalidated: sql.NullBool{Valid: true, Bool: true}},
	)
	if err != nil {
		return err
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "token invalidated")

	return nil
}

func (r *TokenRepository) newRefreshToken(ctx context.Context) (*jwt.Token, error) {
	row, err := r.models.RefreshTokens.Insert(ctx, data.RefreshTokenInput{
		Issuer:     r.Issuer,
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
	tokenString, err := token.SignedString(r.refreshSigningSecret)
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
	token, err := r.parseToken(refreshTokenInput, RefreshTokenType)
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

func (r *TokenRepository) DeleteExpired(ctx context.Context) error {
	logger := logging.LoggerFromContext(ctx)

	logger.LogAttrs(ctx, slog.LevelInfo, "deleting expired tokens")
	timestamp := time.Now().UTC()
	affectedRows, err := r.models.RefreshTokens.DeleteMany(
		ctx,
		data.Filter{ExpirationFrom: &timestamp},
	)
	if err != nil {
		return err
	}
	logger.LogAttrs(
		ctx, slog.LevelInfo, "tokens deleted", slog.Int64("rowsAffected", *affectedRows),
	)

	return nil
}

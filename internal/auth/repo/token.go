package repo

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/r3d5un/islandwind/internal/auth/data"
	"github.com/r3d5un/islandwind/internal/logging"
)

var (
	ErrParsingToken    = errors.New("unable to parse token")
	ErrVerifyingToken  = errors.New("unable to parse token")
	ErrTokenExpired    = errors.New("token expired")
	ErrPrematureToken  = errors.New("token used before valid nbf")
	ErrInvalidIssuedAt = errors.New("token iat timestamp invalid")
	ErrIssuerMismatch  = errors.New("token issuer does not match requirements")
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

// NewJWT create a new signed JWT token string.
func (r *TokenRepository) NewJWT() (*string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.MapClaims{
			"jti": uuid.New(),
			"exp": time.Now().Add(time.Minute * 5).Unix(),
			"nbf": time.Now().Unix(),
			"iat": time.Now().Unix(),
			"iss": r.Issuer,
		},
	)

	tokenString, err := token.SignedString(r.signingSecret)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

// Parse parses a given input JWT string and validates it claims. An error is returned
// if the token cannot be parsed. If the token is invalid in any way a false boolean
// value is returned along with an error describing the fault.
func (r *TokenRepository) Parse(ctx context.Context, input string) (bool, error) {
	logger := logging.LoggerFromContext(ctx)

	logger.Info("parsing JWT")
	token, err := jwt.Parse(input, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrParsingToken
		}
		return r.signingSecret, nil
	})
	if err != nil {
		return false, err
	}
	logger = logger.With(slog.Any("token", token))

	return r.verifyClaims(token)
}

func (r *TokenRepository) NewRefreshToken(ctx context.Context) (*string, error) {
	logger := logging.LoggerFromContext(ctx)

	logger.LogAttrs(ctx, slog.LevelInfo, "creating refresh token")
	row, err := r.models.RefreshTokens.Insert(ctx, data.RefreshTokenInput{
		Issuer:     "",
		Expiration: time.Now().Add(time.Minute * 60),
		IssuedAt:   time.Now(),
		NotBefore:  time.Now(),
	})
	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"unable to create refresh token",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.MapClaims{
			"jti": row.ID,
			"exp": row.Expiration.Unix(),
			"nbf": row.NotBefore.Unix(),
			"iat": row.IssuedAt.Unix(),
			"iss": r.Issuer,
		},
	)

	tokenString, err := token.SignedString(r.signingSecret)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
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
	if exp.Unix() < time.Now().Unix() {
		return false, ErrTokenExpired
	}

	iat, err := claims.GetIssuedAt()
	if err != nil {
		return false, ErrVerifyingToken
	}
	if iat.Unix() > time.Now().Unix() {
		return false, ErrInvalidIssuedAt
	}

	nbf, err := claims.GetNotBefore()
	if err != nil {
		return false, ErrVerifyingToken
	}
	if nbf.Unix() > time.Now().Unix() {
		return false, ErrPrematureToken
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

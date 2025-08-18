package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/islandwind/internal/logging"
)

type Config struct {
	// ConnStr is the database connection string
	//
	// Set through the ISLANDWIND_DB_CONNSTR environment variable
	ConnStr string `json:"-"`
	// MaxOpenConns sets the maximum number of connections to the database
	//
	// Set through the ISLANDWIND_DB_MAXOPENCONNS environment variable
	MaxOpenConns int32 `json:"maxOpenConns"`
	// IdleTimeMinutes is how long idle database connections remain alive set in minutes.
	//
	// Set through the ISLANDWIND_DB_IDLETIMEMINUTES environment variable
	IdleTimeMinutes int `json:"idleTimeMinutes"`
	// TimeoutSeconds sets the timeout for queries set in seconds.
	//
	// Set through the ISLANDWIND_DB_TIMEOUTSECONDS environment variable
	TimeoutSeconds int `json:"timeoutSeconds"`
}

func (c Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("maxOpenConns", int(c.MaxOpenConns)),
		slog.Int("idleTimeMinutes", int(c.IdleTimeMinutes)),
		slog.Int("timeoutSeconds", int(c.TimeoutSeconds)),
	)
}

func (c *Config) TimeoutDuration() time.Duration {
	return time.Duration(c.TimeoutSeconds) * time.Second
}

func (c *Config) IdleTime() time.Duration {
	return time.Duration(c.IdleTimeMinutes) * time.Minute
}

func OpenPool(ctx context.Context, config Config) (*pgxpool.Pool, error) {
	pgxCfg, err := pgxpool.ParseConfig(config.ConnStr)
	if err != nil {
		return nil, err
	}
	pgxCfg.MaxConnIdleTime = config.IdleTime()
	pgxCfg.MaxConns = config.MaxOpenConns

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

type Queryable interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
}

var (
	ErrUnsafeDeleteFilter = errors.New("filter is unsafe")
)

// DeleteManyGuardrail accepts a slice of values and raises and error if all
// given values are nil. If all values are nil an ErrUnsafeDeleteFilter error
// is returned. If a value is not nil the function returns nil, and the
// filter is safe to use.
//
// WARNING: This function assumed any given value in the input slice is a
// pointer that can be checked for nil. A non-pointer value will return
// early, but the the filter may still be unsafe.
func DeleteManyGuardrail(input ...any) error {
	for _, x := range input {
		if v := reflect.ValueOf(x); !v.IsNil() {
			return nil
		}
	}

	return ErrUnsafeDeleteFilter
}

// IsEmpty checks if a slice is nil or empty
func IsEmpty[T comparable](x []*T) bool {
	return len(x) < 1
}

var (
	ErrRecordNotFound                = errors.New("record not found")
	ErrForeignKeyConstraintViolation = errors.New("foreign key constraint violation")
	ErrConstraintViolation           = errors.New("schema constraint violation")
	ErrUniqueConstraintViolation     = errors.New("unique constraint violation")
	ErrNotNullConstraintViolation    = errors.New("not null constraint violation")
	ErrCheckConstraintViolation      = errors.New("check constraint violation")
	ErrSynatxErrorViolation          = errors.New("sql syntax errors")
	ErrUndefinedResource             = errors.New("undefined resource")
)

const (
	// PgxIntegrityConstraintViolationCode is the code for general integrity violations
	PgxIntegrityConstraintViolationCode string = "23000"
	// PgxRestrictViolationCode is the error code for deleting or updating records referenced by other resources
	PgxRestrictViolationCode = "23001"
	// PgxNotNullViolationCode is the error code for attempting to insert or update NULL values in a NOT NULL column
	PgxNotNullViolationCode = "23502"
	// PgxForeignKeyViolationCode is for foreign key violation, e.g. creating values not in the referenced table
	PgxForeignKeyViolationCode = "23503"
	// PgxUniqueViolationCode is the error code for inserting duplicate records
	PgxUniqueViolationCode = "23505"
	// PgxCheckViolationCode is for CHECK violations
	PgxCheckViolationCode = "23514"
	// PgxSyntaxErrorCode is for general syntax errors
	PgxSyntaxErrorCode = "42601"
	// PgxUndefinedColumnCode is the error code for referencing columns that does not exists
	PgxUndefinedColumnCode = "42703"
	// PgxUndefinedTableCode is the error code for refercing tables that does not exists
	PgxUndefinedTableCode = "42P01"
)

func HandleError(ctx context.Context, err error) error {
	logger := logging.LoggerFromContext(ctx).With("pgxError", slog.String("error", err.Error()))

	if errors.Is(err, pgx.ErrNoRows) {
		logger.LogAttrs(ctx, slog.LevelInfo, "no rows found")
		return ErrRecordNotFound
	}

	if errors.Is(err, ErrUnsafeDeleteFilter) {
		logger.LogAttrs(ctx, slog.LevelInfo, "filter unsafe")
		return err
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		logger = logger.With(slog.Group(
			"pgxError",
			slog.String("error", err.Error()),
			slog.String("code", pgErr.Code),
			slog.String("message", pgErr.Message),
			slog.Int("line", int(pgErr.Line)),
		))
		switch pgErr.Code {
		case PgxIntegrityConstraintViolationCode:
			logger.LogAttrs(ctx, slog.LevelError, "constraint violation")
			return ErrConstraintViolation
		case PgxRestrictViolationCode:
			logger.LogAttrs(ctx, slog.LevelError, "record referenced by other resources")
			return ErrConstraintViolation
		case PgxNotNullViolationCode:
			logger.LogAttrs(ctx, slog.LevelError, "null cannot be inserted to non-nullable fields")
			return ErrConstraintViolation
		case PgxForeignKeyViolationCode:
			logger.LogAttrs(ctx, slog.LevelError, "foreign key constraint violation")
			return ErrForeignKeyConstraintViolation
		case PgxUniqueViolationCode:
			logger.LogAttrs(ctx, slog.LevelError, "unique constraint violation")
			return ErrUniqueConstraintViolation
		case PgxCheckViolationCode:
			logger.LogAttrs(ctx, slog.LevelError, "check constraint violation")
			return ErrCheckConstraintViolation
		case PgxSyntaxErrorCode:
			logger.LogAttrs(ctx, slog.LevelError, "syntax error")
			return ErrSynatxErrorViolation
		case PgxUndefinedColumnCode:
			logger.LogAttrs(ctx, slog.LevelError, "referenced column does not exist")
			return ErrUndefinedResource
		case PgxUndefinedTableCode:
			logger.LogAttrs(ctx, slog.LevelError, "referenced table does not exist")
			return ErrUndefinedResource
		default:
			logger.LogAttrs(ctx, slog.LevelError, "unhandled constraint violation")
			return fmt.Errorf("unhandled constraint violation: %s", pgErr.Message)
		}
	}

	logger.LogAttrs(
		ctx, slog.LevelError, "unhandled database error", slog.String("error", err.Error()),
	)
	return err
}

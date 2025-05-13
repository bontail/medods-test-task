package postgresql

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"medods-test-task/internal/models"
)

type AuthStorage interface {
	InsertToken(ctx context.Context, token models.RefreshToken) (id int64, err error)
	BlockedAllTokens(ctx context.Context, guid string) error
	BlockedToken(ctx context.Context, id int64) error
	GetToken(ctx context.Context, userGUID string, id int64) (*models.RefreshToken, error)
}

type DefaultAuthStorage struct {
	db  *pgxpool.Pool
	lgr *slog.Logger
}

func NewDefaultAuthStorage(db *pgxpool.Pool, lgr *slog.Logger) *DefaultAuthStorage {
	return &DefaultAuthStorage{
		db:  db,
		lgr: lgr,
	}
}

func (s *DefaultAuthStorage) InsertToken(ctx context.Context, token models.RefreshToken) (id int64, err error) {
	query, args, err := squirrel.
		Insert("RefreshTokens").
		Columns(
			"user_guid",
			"secret_value",
			"created_at",
			"expires_at",
			"user_agent",
			"ip",
		).
		Values(
			token.UserGUID,
			token.SecretValue,
			token.CreatedAt,
			token.ExpiresAt,
			token.UserAgent,
			token.IP,
		).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("cannot create sql query for insert token: %w", err)
	}
	s.lgr.Debug(
		"create query for insert token",
		slog.String("query", query),
		slog.String("args", fmt.Sprint(args)),
	)

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("cannot execute sql query for insert token: %w", err)
	}
	s.lgr.Debug(
		"execute query for insert token",
		slog.String("query", query),
		slog.String("args", fmt.Sprint(args)),
	)
	defer rows.Close()

	if rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("cannot scan token row: %w", err)
		}
	}

	return
}

func (s *DefaultAuthStorage) BlockedAllTokens(ctx context.Context, guid string) error {
	query, args, err := squirrel.
		Update("RefreshTokens").
		Where(squirrel.Eq{"user_guid": guid}).
		Set("blocked_at", time.Now()).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("cannot create sql query for block tokens: %w", err)
	}
	s.lgr.Debug(
		"create query for block tokens",
		slog.String("query", query),
		slog.String("args", fmt.Sprint(args)),
	)

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("cannot create sql query for block tokens: %w", err)
	}
	s.lgr.Debug(
		"execute query for block tokens",
		slog.String("query", query),
		slog.String("args", fmt.Sprint(args)),
	)

	return nil
}

func (s *DefaultAuthStorage) BlockedToken(ctx context.Context, id int64) error {
	query, args, err := squirrel.
		Update("RefreshTokens").
		Where(squirrel.Eq{"id": id}).
		Set("blocked_at", time.Now()).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("cannot create sql query for block token: %w", err)
	}
	s.lgr.Debug(
		"create query for block token",
		slog.String("query", query),
		slog.String("args", fmt.Sprint(args)),
	)

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("cannot create sql query for block token: %w", err)
	}
	s.lgr.Debug(
		"execute query for block token",
		slog.String("query", query),
		slog.String("args", fmt.Sprint(args)),
	)

	return nil
}

func (s *DefaultAuthStorage) GetToken(ctx context.Context, userGUID string, id int64) (*models.RefreshToken, error) {

	query, args, err := squirrel.
		Select(
			"id",
			"user_guid",
			"secret_value",
			"created_at",
			"expires_at",
			"user_agent",
			"ip",
		).
		From("RefreshTokens").
		Where(
			squirrel.Eq{"id": id},
		).
		Where(
			squirrel.Eq{"user_guid": userGUID},
		).
		Where(
			squirrel.Gt{"expires_at": time.Now()},
		).
		Where(
			squirrel.Eq{"blocked_at": nil},
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot create sql query for get token: %w", err)
	}
	s.lgr.Debug(
		"create query for get token",
		slog.String("query", query),
		slog.String("args", fmt.Sprint(args)),
	)

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot execute sql query for get token: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var token models.RefreshToken
	if err = rows.Scan(
		&token.Id,
		&token.UserGUID,
		&token.SecretValue,
		&token.CreatedAt,
		&token.ExpiresAt,
		&token.UserAgent,
		&token.IP,
	); err != nil {
		return nil, fmt.Errorf("cannot scan user row: %w", err)
	}
	s.lgr.Debug(
		"get token",
		slog.String("data", fmt.Sprint(token)),
	)

	return &token, nil
}

package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"medods-test-task/internal/models"
)

type UserStorage interface {
	InsertUser(ctx context.Context, user models.User) error
	ExistsUser(ctx context.Context, username string) (bool, error)
	WithUsername(username string) GetUserParams
	WithGUID(guid string) GetUserParams
	GetUser(ctx context.Context, params ...GetUserParams) (*models.User, error)
}

type DefaultUserStorage struct {
	db  *pgxpool.Pool
	lgr *slog.Logger
}

func NewDefaultUserStorage(db *pgxpool.Pool, lgr *slog.Logger) *DefaultUserStorage {
	return &DefaultUserStorage{
		db:  db,
		lgr: lgr,
	}
}

func (s *DefaultUserStorage) InsertUser(ctx context.Context, user models.User) error {
	query, args, err := squirrel.
		Insert("Users").
		Columns(
			"guid",
			"username",
			"password",
		).
		Values(
			user.GUID,
			user.Username,
			user.Password,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("cannot create sql query for creating user: %w", err)
	}
	s.lgr.Debug(
		"create query for creating user",
		slog.String("query", query),
		slog.String("args", fmt.Sprint(args)),
	)

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("cannot execute sql query for creating user: %w", err)
	}
	s.lgr.Debug(
		"execute query for creating user",
		slog.String("query", query),
		slog.String("args", fmt.Sprint(args)),
	)

	return nil
}

func (s *DefaultUserStorage) ExistsUser(ctx context.Context, username string) (bool, error) {
	query, args, err := squirrel.
		Select("1").
		From("Users").
		Where(
			squirrel.Eq{"username": username},
		).
		Prefix("SELECT EXISTS(").
		Suffix(")").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("cannot create sql query for check exists user: %w", err)
	}
	s.lgr.Debug(
		"create query for check exists user",
		slog.String("query", query),
		slog.String("args", fmt.Sprint(args)),
	)

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("cannot execute sql query for check exists user: %w", err)
	}
	s.lgr.Debug(
		"execute query for check exists user",
		slog.String("query", query),
		slog.String("args", fmt.Sprint(args)),
	)
	defer rows.Close()
	var exists bool
	if rows.Next() {
		if err = rows.Scan(&exists); err != nil {
			return false, fmt.Errorf("cannot scan exists user: %w", err)
		}
	}

	return exists, nil
}

type GetUserParams func(sb squirrel.SelectBuilder) squirrel.SelectBuilder

func (s *DefaultUserStorage) WithUsername(username string) GetUserParams {
	return func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
		return sb.Where(squirrel.Eq{"username": username})
	}
}

func (s *DefaultUserStorage) WithGUID(guid string) GetUserParams {
	return func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
		return sb.Where(squirrel.Eq{"guid": guid})
	}
}

func (s *DefaultUserStorage) GetUser(ctx context.Context, params ...GetUserParams) (*models.User, error) {
	if len(params) == 0 {
		return nil, errors.New("no params")
	}

	queryHeader := squirrel.
		Select(
			"guid",
			"username",
			"password",
		).
		From("Users")

	for _, param := range params {
		queryHeader = param(queryHeader)
	}

	query, args, err := queryHeader.
		From("users").
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot create sql query for get user: %w", err)
	}
	s.lgr.Debug(
		"create query for get user",
		slog.String("query", query),
		slog.String("args", fmt.Sprint(args)),
	)

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot execute sql query for get user: %w", err)
	}
	defer rows.Close()

	var user models.User
	users := make([]models.User, 0, 1)
	for rows.Next() {
		if err = rows.Scan(
			&user.GUID, &user.Username, &user.Password,
		); err != nil {
			return nil, fmt.Errorf("cannot scan user row: %w", err)
		}
		users = append(users, user)
	}
	s.lgr.Debug(
		"get user",
		slog.String("data", fmt.Sprint(users)),
	)

	if len(users) > 1 {
		return nil, errors.New("multiple users found")
	}

	return &user, nil
}

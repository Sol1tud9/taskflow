package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/Sol1tud9/taskflow/internal/domain"
)

func (s *Storage) Create(ctx context.Context, user *domain.User) error {
	query := squirrel.Insert("users").
		Columns("id", "email", "name", "created_at", "updated_at").
		Values(user.ID, user.Email, user.Name, user.CreatedAt, user.UpdatedAt).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}

	if _, err := s.db.Exec(ctx, sql, args...); err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	return nil
}

func (s *Storage) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := squirrel.Select("id", "email", "name", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	var user domain.User
	err = s.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}

	return &user, nil
}

func (s *Storage) Update(ctx context.Context, user *domain.User) error {
	query := squirrel.Update("users").
		Set("email", user.Email).
		Set("name", user.Name).
		Set("updated_at", user.UpdatedAt).
		Where(squirrel.Eq{"id": user.ID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}

	if _, err := s.db.Exec(ctx, sql, args...); err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	query := squirrel.Delete("users").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}

	if _, err := s.db.Exec(ctx, sql, args...); err != nil {
		return errors.Wrap(err, "failed to delete user")
	}

	return nil
}

func (s *Storage) List(ctx context.Context) ([]*domain.User, error) {
	query := squirrel.Select("id", "email", "name", "created_at", "updated_at").
		From("users").
		OrderBy("created_at DESC").
		Limit(100).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list users")
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, errors.Wrap(err, "failed to scan user")
		}
		users = append(users, &u)
	}

	return users, nil
}


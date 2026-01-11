package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/Sol1tud9/taskflow/pkg/config"
)

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(cfg config.DatabaseConfig) (*Storage, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse config")
	}

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}

	storage := &Storage{db: db}
	if err := storage.initTables(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) initTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(36) PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS teams (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			owner_id VARCHAR(36) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS team_members (
			id VARCHAR(36) PRIMARY KEY,
			team_id VARCHAR(36) NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
			user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			role VARCHAR(50) NOT NULL,
			joined_at TIMESTAMP NOT NULL,
			UNIQUE(team_id, user_id)
		)`,
	}

	for _, query := range queries {
		if _, err := s.db.Exec(context.Background(), query); err != nil {
			return errors.Wrap(err, "failed to create table")
		}
	}

	return nil
}

func (s *Storage) Close() {
	s.db.Close()
}


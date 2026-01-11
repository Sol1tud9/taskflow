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
		`CREATE TABLE IF NOT EXISTS tasks (
			id VARCHAR(36) PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			status VARCHAR(50) NOT NULL DEFAULT 'todo',
			priority VARCHAR(50) NOT NULL DEFAULT 'medium',
			assignee_id VARCHAR(36),
			creator_id VARCHAR(36) NOT NULL,
			team_id VARCHAR(36),
			due_date TIMESTAMP,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS task_history (
			id VARCHAR(36) PRIMARY KEY,
			task_id VARCHAR(36) NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
			user_id VARCHAR(36) NOT NULL,
			field VARCHAR(100) NOT NULL,
			old_value TEXT,
			new_value TEXT,
			changed_at TIMESTAMP NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_team_id ON tasks(team_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_assignee_id ON tasks(assignee_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status)`,
		`CREATE INDEX IF NOT EXISTS idx_task_history_task_id ON task_history(task_id)`,
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


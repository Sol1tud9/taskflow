package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/Sol1tud9/taskflow/internal/domain"
)

func (s *Storage) CreateHistory(ctx context.Context, history *domain.TaskHistory) error {
	query := squirrel.Insert("task_history").
		Columns("id", "task_id", "user_id", "field", "old_value", "new_value", "changed_at").
		Values(history.ID, history.TaskID, history.UserID, history.Field, history.OldValue, history.NewValue, history.ChangedAt).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}

	if _, err := s.db.Exec(ctx, sql, args...); err != nil {
		return errors.Wrap(err, "failed to create task history")
	}

	return nil
}

func (s *Storage) GetHistoryByTaskID(ctx context.Context, taskID string) ([]*domain.TaskHistory, error) {
	query := squirrel.Select("id", "task_id", "user_id", "field", "old_value", "new_value", "changed_at").
		From("task_history").
		Where(squirrel.Eq{"task_id": taskID}).
		OrderBy("changed_at DESC").
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get task history")
	}
	defer rows.Close()

	var history []*domain.TaskHistory
	for rows.Next() {
		var h domain.TaskHistory
		if err := rows.Scan(&h.ID, &h.TaskID, &h.UserID, &h.Field, &h.OldValue, &h.NewValue, &h.ChangedAt); err != nil {
			return nil, errors.Wrap(err, "failed to scan task history")
		}
		history = append(history, &h)
	}

	return history, nil
}


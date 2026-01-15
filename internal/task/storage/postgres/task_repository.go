package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/Sol1tud9/taskflow/internal/domain"
	"github.com/Sol1tud9/taskflow/internal/task/usecase"
)

func (s *Storage) Create(ctx context.Context, task *domain.Task) error {
	query := squirrel.Insert("tasks").
		Columns("id", "title", "description", "status", "priority", "assignee_id", "creator_id", "team_id", "due_date", "created_at", "updated_at").
		Values(task.ID, task.Title, task.Description, task.Status, task.Priority, task.AssigneeID, task.CreatorID, task.TeamID, task.DueDate, task.CreatedAt, task.UpdatedAt).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}

	if _, err := s.db.Exec(ctx, sql, args...); err != nil {
		return errors.Wrap(err, "failed to create task")
	}

	return nil
}

func (s *Storage) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	query := squirrel.Select("id", "title", "description", "status", "priority", "assignee_id", "creator_id", "team_id", "due_date", "created_at", "updated_at").
		From("tasks").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	var task domain.Task
	err = s.db.QueryRow(ctx, sql, args...).Scan(
		&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority,
		&task.AssigneeID, &task.CreatorID, &task.TeamID, &task.DueDate,
		&task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get task")
	}

	return &task, nil
}

func (s *Storage) List(ctx context.Context, filter usecase.TaskFilter) ([]*domain.Task, int, error) {
	query := squirrel.Select("id", "title", "description", "status", "priority", "assignee_id", "creator_id", "team_id", "due_date", "created_at", "updated_at").
		From("tasks").
		PlaceholderFormat(squirrel.Dollar)

	if filter.TeamID != "" {
		query = query.Where(squirrel.Eq{"team_id": filter.TeamID})
	}
	if filter.AssigneeID != "" {
		query = query.Where(squirrel.Eq{"assignee_id": filter.AssigneeID})
	}
	if filter.Status != "" {
		query = query.Where(squirrel.Eq{"status": filter.Status})
	}

	if filter.Limit > 0 {
		query = query.Limit(uint64(filter.Limit))
	}
	if filter.Offset > 0 {
		query = query.Offset(uint64(filter.Offset))
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to build query")
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to list tasks")
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(
			&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority,
			&t.AssigneeID, &t.CreatorID, &t.TeamID, &t.DueDate,
			&t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, 0, errors.Wrap(err, "failed to scan task")
		}
		tasks = append(tasks, &t)
	}

	countQuery := squirrel.Select("COUNT(*)").From("tasks").PlaceholderFormat(squirrel.Dollar)
	if filter.TeamID != "" {
		countQuery = countQuery.Where(squirrel.Eq{"team_id": filter.TeamID})
	}
	if filter.AssigneeID != "" {
		countQuery = countQuery.Where(squirrel.Eq{"assignee_id": filter.AssigneeID})
	}
	if filter.Status != "" {
		countQuery = countQuery.Where(squirrel.Eq{"status": filter.Status})
	}

	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to build count query")
	}

	var total int
	if err := s.db.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, errors.Wrap(err, "failed to count tasks")
	}

	return tasks, total, nil
}

func (s *Storage) Update(ctx context.Context, task *domain.Task) error {
	query := squirrel.Update("tasks").
		Set("title", task.Title).
		Set("description", task.Description).
		Set("status", task.Status).
		Set("priority", task.Priority).
		Set("assignee_id", task.AssigneeID).
		Set("due_date", task.DueDate).
		Set("updated_at", task.UpdatedAt).
		Where(squirrel.Eq{"id": task.ID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}

	if _, err := s.db.Exec(ctx, sql, args...); err != nil {
		return errors.Wrap(err, "failed to update task")
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	query := squirrel.Delete("tasks").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}

	if _, err := s.db.Exec(ctx, sql, args...); err != nil {
		return errors.Wrap(err, "failed to delete task")
	}

	return nil
}


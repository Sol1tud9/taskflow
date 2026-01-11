package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/Sol1tud9/taskflow/internal/domain"
)

func (s *Storage) CreateTeam(ctx context.Context, team *domain.Team) error {
	query := squirrel.Insert("teams").
		Columns("id", "name", "owner_id", "created_at", "updated_at").
		Values(team.ID, team.Name, team.OwnerID, team.CreatedAt, team.UpdatedAt).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}

	if _, err := s.db.Exec(ctx, sql, args...); err != nil {
		return errors.Wrap(err, "failed to create team")
	}

	return nil
}

func (s *Storage) GetTeamByID(ctx context.Context, id string) (*domain.Team, error) {
	query := squirrel.Select("id", "name", "owner_id", "created_at", "updated_at").
		From("teams").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	var team domain.Team
	err = s.db.QueryRow(ctx, sql, args...).Scan(
		&team.ID, &team.Name, &team.OwnerID, &team.CreatedAt, &team.UpdatedAt,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get team")
	}

	return &team, nil
}

func (s *Storage) UpdateTeam(ctx context.Context, team *domain.Team) error {
	query := squirrel.Update("teams").
		Set("name", team.Name).
		Set("owner_id", team.OwnerID).
		Set("updated_at", team.UpdatedAt).
		Where(squirrel.Eq{"id": team.ID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}

	if _, err := s.db.Exec(ctx, sql, args...); err != nil {
		return errors.Wrap(err, "failed to update team")
	}

	return nil
}

func (s *Storage) DeleteTeam(ctx context.Context, id string) error {
	query := squirrel.Delete("teams").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}

	if _, err := s.db.Exec(ctx, sql, args...); err != nil {
		return errors.Wrap(err, "failed to delete team")
	}

	return nil
}

func (s *Storage) AddTeamMember(ctx context.Context, member *domain.TeamMember) error {
	query := squirrel.Insert("team_members").
		Columns("id", "team_id", "user_id", "role", "joined_at").
		Values(member.ID, member.TeamID, member.UserID, member.Role, member.JoinedAt).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}

	if _, err := s.db.Exec(ctx, sql, args...); err != nil {
		return errors.Wrap(err, "failed to add team member")
	}

	return nil
}

func (s *Storage) GetTeamMembersByTeamID(ctx context.Context, teamID string) ([]*domain.TeamMember, error) {
	query := squirrel.Select("id", "team_id", "user_id", "role", "joined_at").
		From("team_members").
		Where(squirrel.Eq{"team_id": teamID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get team members")
	}
	defer rows.Close()

	var members []*domain.TeamMember
	for rows.Next() {
		var m domain.TeamMember
		if err := rows.Scan(&m.ID, &m.TeamID, &m.UserID, &m.Role, &m.JoinedAt); err != nil {
			return nil, errors.Wrap(err, "failed to scan team member")
		}
		members = append(members, &m)
	}

	return members, nil
}

func (s *Storage) RemoveTeamMember(ctx context.Context, teamID, userID string) error {
	query := squirrel.Delete("team_members").
		Where(squirrel.And{
			squirrel.Eq{"team_id": teamID},
			squirrel.Eq{"user_id": userID},
		}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}

	if _, err := s.db.Exec(ctx, sql, args...); err != nil {
		return errors.Wrap(err, "failed to remove team member")
	}

	return nil
}

func (s *Storage) ListTeams(ctx context.Context) ([]*domain.Team, error) {
	query := squirrel.Select("id", "name", "owner_id", "created_at", "updated_at").
		From("teams").
		OrderBy("created_at DESC").
		Limit(100).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list teams")
	}
	defer rows.Close()

	var teams []*domain.Team
	for rows.Next() {
		var t domain.Team
		if err := rows.Scan(&t.ID, &t.Name, &t.OwnerID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, errors.Wrap(err, "failed to scan team")
		}
		teams = append(teams, &t)
	}

	return teams, nil
}


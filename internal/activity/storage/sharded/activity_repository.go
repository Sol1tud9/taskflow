package sharded

import (
	"context"
	"sort"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/Sol1tud9/taskflow/internal/activity/repository"
	"github.com/Sol1tud9/taskflow/internal/domain"
)

func (s *ShardedStorage) Create(ctx context.Context, activity *domain.Activity) error {
	shard := s.GetShardForUser(activity.UserID)

	query := squirrel.Insert("activities").
		Columns("id", "user_id", "entity_type", "entity_id", "action", "metadata", "created_at").
		Values(activity.ID, activity.UserID, activity.EntityType, activity.EntityID, activity.Action, activity.Metadata, activity.CreatedAt).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}

	if _, err := shard.Exec(ctx, sql, args...); err != nil {
		return errors.Wrap(err, "failed to create activity")
	}

	return nil
}

func (s *ShardedStorage) GetByUserID(ctx context.Context, userID string, filter repository.ActivityFilter) ([]*domain.Activity, int, error) {
	shard := s.GetShardForUser(userID)

	return s.queryActivities(ctx, shard, squirrel.Eq{"user_id": userID}, filter)
}

func (s *ShardedStorage) GetByEntity(ctx context.Context, entityType, entityID string, filter repository.ActivityFilter) ([]*domain.Activity, int, error) {
	var allActivities []*domain.Activity
	var totalCount int

	for _, shard := range s.GetAllShards() {
		activities, count, err := s.queryActivities(ctx, shard, squirrel.And{
			squirrel.Eq{"entity_type": entityType},
			squirrel.Eq{"entity_id": entityID},
		}, filter)
		if err != nil {
			return nil, 0, err
		}
		allActivities = append(allActivities, activities...)
		totalCount += count
	}

	return allActivities, totalCount, nil
}

func (s *ShardedStorage) GetAll(ctx context.Context, filter repository.ActivityFilter) ([]*domain.Activity, int, error) {
	var allActivities []*domain.Activity
	var totalCount int

	for _, shard := range s.GetAllShards() {
		activities, count, err := s.queryActivities(ctx, shard, squirrel.Expr("1=1"), filter)
		if err != nil {
			return nil, 0, err
		}
		allActivities = append(allActivities, activities...)
		totalCount += count
	}

	sort.Slice(allActivities, func(i, j int) bool {
		return allActivities[i].CreatedAt.After(allActivities[j].CreatedAt)
	})

	if filter.Limit > 0 && len(allActivities) > filter.Limit {
		allActivities = allActivities[:filter.Limit]
	}

	return allActivities, totalCount, nil
}

func (s *ShardedStorage) queryActivities(ctx context.Context, shard *pgxpool.Pool, where squirrel.Sqlizer, filter repository.ActivityFilter) ([]*domain.Activity, int, error) {
	query := squirrel.Select("id", "user_id", "entity_type", "entity_id", "action", "metadata", "created_at").
		From("activities").
		Where(where).
		OrderBy("created_at DESC").
		PlaceholderFormat(squirrel.Dollar)

	if filter.FromTimestamp > 0 {
		query = query.Where(squirrel.GtOrEq{"created_at": time.Unix(filter.FromTimestamp, 0)})
	}
	if filter.ToTimestamp > 0 {
		query = query.Where(squirrel.LtOrEq{"created_at": time.Unix(filter.ToTimestamp, 0)})
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

	rows, err := shard.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to query activities")
	}
	defer rows.Close()

	var activities []*domain.Activity
	for rows.Next() {
		var a domain.Activity
		if err := rows.Scan(&a.ID, &a.UserID, &a.EntityType, &a.EntityID, &a.Action, &a.Metadata, &a.CreatedAt); err != nil {
			return nil, 0, errors.Wrap(err, "failed to scan activity")
		}
		activities = append(activities, &a)
	}

	countQuery := squirrel.Select("COUNT(*)").From("activities").Where(where).PlaceholderFormat(squirrel.Dollar)
	if filter.FromTimestamp > 0 {
		countQuery = countQuery.Where(squirrel.GtOrEq{"created_at": time.Unix(filter.FromTimestamp, 0)})
	}
	if filter.ToTimestamp > 0 {
		countQuery = countQuery.Where(squirrel.LtOrEq{"created_at": time.Unix(filter.ToTimestamp, 0)})
	}

	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to build count query")
	}

	var total int
	if err := shard.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, errors.Wrap(err, "failed to count activities")
	}

	return activities, total, nil
}


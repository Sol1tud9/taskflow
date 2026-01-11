package sharded

import (
	"context"
	"fmt"
	"hash/fnv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/Sol1tud9/taskflow/pkg/config"
)

type ShardedStorage struct {
	shards     []*pgxpool.Pool
	shardCount int
}

func NewShardedStorage(cfg config.ShardingConfig) (*ShardedStorage, error) {
	shards := make([]*pgxpool.Pool, len(cfg.Shards))

	for i, shardCfg := range cfg.Shards {
		connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			shardCfg.Username, shardCfg.Password, shardCfg.Host, shardCfg.Port, shardCfg.Name, shardCfg.SSLMode)

		poolConfig, err := pgxpool.ParseConfig(connString)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse config for shard %d", i)
		}

		db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to connect to shard %d", i)
		}

		shards[i] = db
	}

	storage := &ShardedStorage{
		shards:     shards,
		shardCount: cfg.ShardCount,
	}

	if err := storage.initTables(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *ShardedStorage) initTables() error {
	query := `CREATE TABLE IF NOT EXISTS activities (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL,
		entity_type VARCHAR(50) NOT NULL,
		entity_id VARCHAR(36) NOT NULL,
		action VARCHAR(50) NOT NULL,
		metadata TEXT,
		created_at TIMESTAMP NOT NULL
	)`

	indexQueries := []string{
		`CREATE INDEX IF NOT EXISTS idx_activities_user_id ON activities(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_activities_entity ON activities(entity_type, entity_id)`,
		`CREATE INDEX IF NOT EXISTS idx_activities_created_at ON activities(created_at)`,
	}

	for i, shard := range s.shards {
		if _, err := shard.Exec(context.Background(), query); err != nil {
			return errors.Wrapf(err, "failed to create table on shard %d", i)
		}
		for _, indexQuery := range indexQueries {
			if _, err := shard.Exec(context.Background(), indexQuery); err != nil {
				return errors.Wrapf(err, "failed to create index on shard %d", i)
			}
		}
	}

	return nil
}

func (s *ShardedStorage) GetShardForUser(userID string) *pgxpool.Pool {
	shardIndex := s.hashUserID(userID) % s.shardCount
	return s.shards[shardIndex]
}

func (s *ShardedStorage) hashUserID(userID string) int {
	h := fnv.New32a()
	h.Write([]byte(userID))
	return int(h.Sum32())
}

func (s *ShardedStorage) GetAllShards() []*pgxpool.Pool {
	return s.shards
}

func (s *ShardedStorage) Close() {
	for _, shard := range s.shards {
		shard.Close()
	}
}


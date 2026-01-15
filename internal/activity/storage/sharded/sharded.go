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
	shards        []*pgxpool.Pool
	shardCount    int
	bucketCount   int
	bucketToShard map[int]int 
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

	bucketCount := cfg.BucketCount
	if bucketCount == 0 {
		bucketCount = cfg.ShardCount
	}

	bucketToShard := cfg.BucketMapping
	if bucketToShard == nil || len(bucketToShard) == 0 {
		bucketToShard = make(map[int]int)
		for i := 0; i < bucketCount; i++ {
			bucketToShard[i] = i % cfg.ShardCount
		}
	}

	for bucketID, shardID := range bucketToShard {
		if bucketID < 0 || bucketID >= bucketCount {
			return nil, errors.Errorf("invalid bucket_id %d: must be in range [0, %d)", bucketID, bucketCount)
		}
		if shardID < 0 || shardID >= len(shards) {
			return nil, errors.Errorf("invalid shard_id %d for bucket %d: must be in range [0, %d)", shardID, bucketID, len(shards))
		}
	}

	if len(bucketToShard) < bucketCount {
		return nil, errors.Errorf("bucket_mapping incomplete: expected %d buckets, got %d", bucketCount, len(bucketToShard))
	}
	for i := 0; i < bucketCount; i++ {
		if _, exists := bucketToShard[i]; !exists {
			return nil, errors.Errorf("bucket %d has no mapping to shard", i)
		}
	}

	storage := &ShardedStorage{
		shards:        shards,
		shardCount:    cfg.ShardCount,
		bucketCount:   bucketCount,
		bucketToShard: bucketToShard,
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
	bucketID := s.getBucketForUser(userID)
	shardID := s.bucketToShard[bucketID]
	return s.shards[shardID]
}

func (s *ShardedStorage) getBucketForUser(userID string) int {
	hash := s.hashUserID(userID)
	bucketID := hash % s.bucketCount
	return bucketID
}

func (s *ShardedStorage) hashUserID(userID string) int {
	h := fnv.New32a()
	h.Write([]byte(userID))
	hashValue := int(h.Sum32())
	if hashValue < 0 {
		hashValue = -hashValue
	}
	return hashValue
}

func (s *ShardedStorage) GetAllShards() []*pgxpool.Pool {
	return s.shards
}

func (s *ShardedStorage) Close() {
	for _, shard := range s.shards {
		shard.Close()
	}
}


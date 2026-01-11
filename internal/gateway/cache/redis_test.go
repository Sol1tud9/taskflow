package cache

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/Sol1tud9/taskflow/pkg/config"
)

type RedisCacheSuite struct {
	suite.Suite
	ctx    context.Context
	cache  *RedisCache
	client redismock.ClientMock
}

func (s *RedisCacheSuite) SetupTest() {
	s.ctx = context.Background()
	cfg := config.RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
		CacheTTL: 300,
	}
	
	db, mock := redismock.NewClientMock()
	s.cache = &RedisCache{
		client: db,
		ttl:    time.Duration(cfg.CacheTTL) * time.Second,
	}
	s.client = mock
}

func (s *RedisCacheSuite) TestGet_Success() {
	key := "test:key"
	value := `{"id":"123","name":"test"}`

	s.client.ExpectGet(key).SetVal(value)

	var result map[string]string
	err := s.cache.Get(s.ctx, key, &result)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "123", result["id"])
	assert.Equal(s.T(), "test", result["name"])
}

func (s *RedisCacheSuite) TestGet_NotFound() {
	key := "test:key"

	s.client.ExpectGet(key).RedisNil()

	var result map[string]string
	err := s.cache.Get(s.ctx, key, &result)

	assert.Error(s.T(), err)
}

func (s *RedisCacheSuite) TestSet_Success() {
	key := "test:key"
	value := map[string]string{"id": "123", "name": "test"}
	expectedJSON := `{"id":"123","name":"test"}`

	s.client.ExpectSet(key, []byte(expectedJSON), s.cache.ttl).SetVal("OK")

	err := s.cache.Set(s.ctx, key, value)

	assert.NoError(s.T(), err)
}

func (s *RedisCacheSuite) TestDelete_Success() {
	key := "test:key"

	s.client.ExpectDel(key).SetVal(1)

	err := s.cache.Delete(s.ctx, key)

	assert.NoError(s.T(), err)
}

func TestRedisCacheSuite(t *testing.T) {
	suite.Run(t, new(RedisCacheSuite))
}


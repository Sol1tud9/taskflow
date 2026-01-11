package config

import (
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Name     string `mapstructure:"name"`
	LogLevel string `mapstructure:"log_level"`
}

type ServerConfig struct {
	GRPCPort int `mapstructure:"grpc_port"`
	HTTPPort int `mapstructure:"http_port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type KafkaConfig struct {
	Brokers        []string          `mapstructure:"brokers"`
	Topics         map[string]string `mapstructure:"topics"`
	ConsumerGroups map[string]string `mapstructure:"consumer_groups"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	CacheTTL int    `mapstructure:"cache_ttl"`
}

type ShardConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type ShardingConfig struct {
	Enabled    bool          `mapstructure:"enabled"`
	ShardCount int           `mapstructure:"shard_count"`
	Shards     []ShardConfig `mapstructure:"shards"`
}

type ServiceEndpoint struct {
	GRPCAddr string `mapstructure:"grpc_addr"`
}

type ServicesConfig struct {
	User     ServiceEndpoint `mapstructure:"user"`
	Task     ServiceEndpoint `mapstructure:"task"`
	Activity ServiceEndpoint `mapstructure:"activity"`
}

type UserServiceConfig struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

type TaskServiceConfig struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

type ActivityServiceConfig struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Sharding ShardingConfig `mapstructure:"sharding"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

type GatewayConfig struct {
	App          AppConfig      `mapstructure:"app"`
	Server       ServerConfig   `mapstructure:"server"`
	Services     ServicesConfig `mapstructure:"services"`
	Redis        RedisConfig    `mapstructure:"redis"`
	UserDB       DatabaseConfig `mapstructure:"user_db"`
	TaskDB       DatabaseConfig `mapstructure:"task_db"`
	ActivityDB   ShardingConfig `mapstructure:"activity_db"`
	Kafka        KafkaConfig    `mapstructure:"kafka"`
}

func Load[T any](path string) (*T, error) {
	viper.SetConfigFile(path)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg T
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}


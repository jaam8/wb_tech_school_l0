package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jaam8/wb_tech_school_l0/pkg/kafka"
	"github.com/jaam8/wb_tech_school_l0/pkg/postgres"
	"github.com/jaam8/wb_tech_school_l0/pkg/redis"
)

type Config struct {
	Kafka    kafka.Config    `yaml:"kafka" env-prefix:"KAFKA_"`
	Redis    redis.Config    `yaml:"redis" env-prefix:"REDIS_"`
	Postgres postgres.Config `yaml:"postgres" env-prefix:"POSTGRES_"`
	Service  ServiceConfig   `yaml:"service" env-prefix:"SERVICE_"`
}

type ServiceConfig struct {
	Port           uint16 `yaml:"port" env-prefix:"PORT"`
	KafkaTopic     string `yaml:"kafka_topic" env-prefix:"KAFKA_TOPIC"`
	KafkaGroupID   string `yaml:"kafka_group_id" env-prefix:"KAFKA_GROUP"`
	RedisDB        int    `yaml:"redis_db" env:"REDIS_DB" env-default:"0"`
	BatchSize      int    `yaml:"batch_size" env:"BATCH_SIZE"`
	FlushTimeout   int    `yaml:"flush_timeout" env:"FLUSH_TIMEOUT"`
	MaxRetries     uint   `yaml:"max_retries" env:"MAX_RETRIES"`
	BaseRetryDelay int    `yaml:"base_retry_delay" env:"BASE_RETRY_DELAY"`
	MigrationsPath string `yaml:"migrations_path" env:"MIGRATIONS_PATH"`
}

func New() (Config, error) {
	var cfg Config
	// docker workdir - app/
	if err := cleanenv.ReadConfig("../configs/config.yaml", &cfg); err != nil {
		fmt.Println(err.Error())
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return Config{}, fmt.Errorf("failed to read env vars: %v", err)
		}
	}

	return cfg, nil
}

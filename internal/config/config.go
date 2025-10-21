package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jaam8/wb_tech_school_l0/pkg/kafka"
	lrucache "github.com/jaam8/wb_tech_school_l0/pkg/lru-cache"
	"github.com/jaam8/wb_tech_school_l0/pkg/postgres"
)

type Config struct {
	Kafka    kafka.Config    `env-prefix:"KAFKA_"    yaml:"kafka"`
	Cache    lrucache.Config `env-prefix:"CACHE_"    yaml:"cache"`
	Postgres postgres.Config `env-prefix:"POSTGRES_" yaml:"postgres"`
	Service  AppConfig       `env-prefix:"APP_"      yaml:"service"`

	LogLevel        string `env:"LOG_LEVEL"         env-default:"info"         yaml:"log_level"`
	MigrationsPath  string `env:"MIGRATIONS_PATH"   env-default:"./migrations" yaml:"migrations_path"`
	FakeOrdersCount int    `env:"FAKE_ORDERS_COUNT" env-default:"10"           yaml:"fake_orders_count"`
}

type AppConfig struct {
	Port                   uint16 `env:"PORT"                     env-default:"8080"      yaml:"port"`
	KafkaTopic             string `env:"KAFKA_TOPIC"              yaml:"kafka_topic"`
	KafkaGroupID           string `env:"KAFKA_GROUP_ID"           yaml:"kafka_group_id"`
	BatchSize              int    `env:"BATCH_SIZE"               env-default:"1"         yaml:"batch_size"`
	FlushTimeout           int    `env:"FLUSH_TIMEOUT"            env-default:"1"         yaml:"flush_timeout"`
	KafkaNumPartitions     int    `env:"KAFKA_NUM_PARTITIONS"     env-default:"1"         yaml:"kafka_num_partitions"`
	KafkaReplicationFactor int    `env:"KAFKA_REPLICATION_FACTOR" env-default:"1"         yaml:"kafka_replication_factor"`
	MaxRetries             int    `env:"MAX_RETRIES"              env-default:"5"         yaml:"max_retries"`
	BaseRetryDelay         int    `env:"BASE_RETRY_DELAY"         yaml:"base_retry_delay"`
}

func New() (Config, error) {
	var cfg Config
	// docker workdir - app/
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		log.Println(err)
		if err = cleanenv.ReadEnv(&cfg); err != nil {
			return Config{}, fmt.Errorf("failed to read env vars: %v", err)
		}
	}

	return cfg, nil
}

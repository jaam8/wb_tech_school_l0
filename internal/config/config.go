package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jaam8/wb_tech_school_l0/pkg/kafka"
	lrucache "github.com/jaam8/wb_tech_school_l0/pkg/lru-cache"
	"github.com/jaam8/wb_tech_school_l0/pkg/postgres"
)

type Config struct {
	Kafka          kafka.Config    `yaml:"kafka" env-prefix:"KAFKA_"`
	Cache          lrucache.Config `yaml:"cache" env-prefix:"CACHE_"`
	Postgres       postgres.Config `yaml:"postgres" env-prefix:"POSTGRES_"`
	Service        AppConfig       `yaml:"service" env-prefix:"APP_"`
	LogLevel       string          `yaml:"log_level" env:"LOG_LEVEL" env-default:"info"`
	MigrationsPath string          `yaml:"migrations_path" env:"MIGRATIONS_PATH" env-default:"./migrations"`
}

type AppConfig struct {
	Port           uint16 `yaml:"port" env:"PORT" env-default:"8080"`
	KafkaTopic     string `yaml:"kafka_topic" env:"KAFKA_TOPIC"`
	KafkaGroupID   string `yaml:"kafka_group_id" env:"KAFKA_GROUP_ID"`
	BatchSize      int    `yaml:"batch_size" env:"BATCH_SIZE" env-default:"1"`
	FlushTimeout   int    `yaml:"flush_timeout" env:"FLUSH_TIMEOUT" env-default:"1"`
	MaxRetries     uint   `yaml:"max_retries" env:"MAX_RETRIES"`
	BaseRetryDelay int    `yaml:"base_retry_delay" env:"BASE_RETRY_DELAY"`
}

func New() (Config, error) {
	var cfg Config
	// docker workdir - app/
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		fmt.Println(err.Error())
		if err = cleanenv.ReadEnv(&cfg); err != nil {
			return Config{}, fmt.Errorf("failed to read env vars: %v", err)
		}
	}

	return cfg, nil
}

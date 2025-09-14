package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type Config struct {
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port     uint16 `yaml:"port" env:"PORT" env-default:"5432"`
	Username string `yaml:"user" env:"USER" env-default:"postgres"`
	Password string `yaml:"password" env:"PASSWORD" env-default:"1234"`
	Database string `yaml:"db" env:"DB" env-default:"postgres"`
	MaxConns int32  `yaml:"max_conns" env:"MAX_CONNS" env-default:"10"`
	MinConns int32  `yaml:"min_conns" env:"MIN_CONNS" env-default:"5"`
}

func New(ctx context.Context, config Config) (*pgxpool.Pool, error) {
	dsn := config.GetDsn()
	dsn += fmt.Sprintf("&pool_max_conns=%d&pool_min_conns=%d",
		config.MaxConns,
		config.MinConns,
	)

	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return conn, nil
}

func Migrate(ctx context.Context, config Config, migrationsPath string) error {
	dsn := config.GetDsn()

	db, err := goose.OpenDBWithDriver("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err = goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}

	if err = goose.Up(db, migrationsPath); err != nil && !errors.Is(err, goose.ErrNoMigrations) {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("migrated successfully")
	return nil
}

func (c *Config) GetDsn() string {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
	return dsn
}

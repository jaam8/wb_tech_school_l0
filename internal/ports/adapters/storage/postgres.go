package storage

import (
	"github.com/jaam8/wb_tech_school_l0/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAdapter struct {
	pool *pgxpool.Pool
}

func NewPostgresAdapter(pool *pgxpool.Pool) *PostgresAdapter {
	return &PostgresAdapter{
		pool: pool,
	}
}

func GetOrder(id string) *models.Order {
	return nil
}

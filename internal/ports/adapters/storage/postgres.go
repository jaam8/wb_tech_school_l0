package storage

import (
	"context"

	"github.com/jaam8/wb_tech_school_l0/internal/models"
	"github.com/jackc/pgx/v5"
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

func (a *PostgresAdapter) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	query := `
	SELECT * FROM orders 
	WHERE order_uid = $1;
`
	rows, err := a.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	order, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Order])
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (a *PostgresAdapter) SaveOrder(ctx context.Context, order *models.Order) error {
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var deliveryId uint64
	query := `
		INSERT INTO deliveries (name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err = tx.QueryRow(ctx, query,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	).Scan(&deliveryId)
	if err != nil {
		return err
	}

	query = `
        INSERT INTO payments (transaction, request_id, currency, provider, amount,
                              payment_dt, bank, delivery_cost, goods_total, custom_fee)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err = tx.Exec(ctx, query,
		order.Payment.Transaction,
		order.Payment.RequestId,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	if err != nil {
		return err
	}

	query = `
        INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_transaction,
                            locale, internal_signature, customer_id, delivery_service,
                            shardkey, sm_id, date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err = tx.Exec(ctx, query,
		order.OrderUid,
		order.TrackNumber,
		order.Entry,
		deliveryId,
		order.Payment.Transaction,
		order.Locale,
		order.InternalSignature,
		order.CustomerId,
		order.DeliveryService,
		order.Shardkey,
		order.SmId,
		order.DateCreated,
		order.OofShard,
	)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO items (chrt_id, track_number, price, rid, name, sale,
						   size, total_price, nm_id, brand, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (chrt_id) DO NOTHING
	`
	batch := &pgx.Batch{}
	for _, item := range order.Items {
		batch.Queue(query,
			item.ChrtId,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmId,
			item.Brand,
			item.Status,
		)

		batch.Queue(`
        INSERT INTO order_items (order_uid, item_chrt_id)
        VALUES ($1,$2)`,
			order.OrderUid, item.ChrtId,
		)
	}

	err = tx.SendBatch(ctx, batch).Close()
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

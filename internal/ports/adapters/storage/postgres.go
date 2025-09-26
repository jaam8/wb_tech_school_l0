package storage

import (
	"context"
	"errors"

	"github.com/jaam8/wb_tech_school_l0/internal/models"
	errs "github.com/jaam8/wb_tech_school_l0/pkg/errors"
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrOrderNotFound
		}
		return nil, err
	}

	return &order, nil
}

func (a *PostgresAdapter) SaveOrders(ctx context.Context, orders ...*models.Order) error {
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// queries
	deliveriesQuery := `
		INSERT INTO deliveries (name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	paymentsQuery := `
        INSERT INTO payments (transaction, request_id, currency, provider, amount,
                              payment_dt, bank, delivery_cost, goods_total, custom_fee)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	ordersQuery := `
        INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_transaction,
                            locale, internal_signature, customer_id, delivery_service,
                            shardkey, sm_id, date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	itemsQuery := `
		INSERT INTO items (chrt_id, track_number, price, rid, name, sale,
						   size, total_price, nm_id, brand, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (chrt_id) DO NOTHING
	`

	for _, order := range orders {
		var deliveryId uint64
		err = tx.QueryRow(ctx, deliveriesQuery,
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

		_, err = tx.Exec(ctx, paymentsQuery,
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

		_, err = tx.Exec(ctx, ordersQuery,
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

		batch := &pgx.Batch{}
		for _, item := range order.Items {
			batch.Queue(itemsQuery,
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
	}

	return tx.Commit(ctx)
}

package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/karambo3a/wbtech_test_task/internal/model"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

type dbOrder struct {
	OrderUID          string    `db:"order_uid"`
	TrackNumber       string    `db:"track_number"`
	Entry             string    `db:"entry"`
	Locale            string    `db:"locale"`
	InternalSignature string    `db:"internal_signature"`
	CustomerID        string    `db:"customer_id"`
	DeliveryService   string    `db:"delivery_service"`
	Shardkey          string    `db:"shardkey"`
	SmID              int       `db:"sm_id"`
	DateCreated       time.Time `db:"date_created"`
	OofShard          string    `db:"oof_shard"`

	DeliveryName    string `db:"delivery_name"`
	DeliveryPhone   string `db:"delivery_phone"`
	DeliveryZip     string `db:"delivery_zip"`
	DeliveryCity    string `db:"delivery_city"`
	DeliveryAddress string `db:"delivery_address"`
	DeliveryRegion  string `db:"delivery_region"`
	DeliveryEmail   string `db:"delivery_email"`

	PaymentTransaction  string `db:"payment_transaction"`
	PaymentRequestID    string `db:"payment_request_id"`
	PaymentCurrency     string `db:"payment_currency"`
	PaymentProvider     string `db:"payment_provider"`
	PaymentAmount       int    `db:"payment_amount"`
	PaymentPaymentDt    int64  `db:"payment_payment_dt"`
	PaymentBank         string `db:"payment_bank"`
	PaymentDeliveryCost int    `db:"payment_delivery_cost"`
	PaymentGoodsTotal   int    `db:"payment_goods_total"`
	PaymentCustomFee    int    `db:"payment_custom_fee"`
}

const (
	getOrderQuery = `SELECT
            o.order_uid,
            o.track_number,
            o.entry,
            o.locale,
            o.internal_signature,
            o.customer_id,
            o.delivery_service,
            o.shardkey,
            o.sm_id,
            o.date_created,
            o.oof_shard,
            d.name AS delivery_name,
            d.phone AS delivery_phone,
            d.zip AS delivery_zip,
            d.city AS delivery_city,
            d.address AS delivery_address,
            d.region AS delivery_region,
            d.email AS delivery_email,
            p.transaction AS payment_transaction,
            p.request_id AS payment_request_id,
            p.currency AS payment_currency,
            p.provider AS payment_provider,
            p.amount AS payment_amount,
            p.payment_dt AS payment_payment_dt,
            p.bank AS payment_bank,
            p.delivery_cost AS payment_delivery_cost,
            p.goods_total AS payment_goods_total,
            p.custom_fee AS payment_custom_fee
        FROM orders o
        LEFT JOIN deliveries d ON o.delivery_id = d.id
        LEFT JOIN payments p ON o.payment_id = p.id`
	getItemsQuery = `SELECT
            i.chrt_id,
            i.track_number,
            i.price,
            i.rid,
            i.name,
            i.sale,
            i.size,
            i.total_price,
            i.nm_id,
            i.brand,
            i.status
        FROM items i
		LEFT JOIN orders_x_items oi ON i.id = oi.item_id
		LEFT JOIN orders o ON o.order_uid = oi.order_uid
		WHERE o.order_uid = $1`
	insertDeliveryQuery = `INSERT INTO deliveries (name, phone, zip, city, address, region, email)
							VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (name, phone, zip, city, address, region, email) DO NOTHING RETURNING id;`
	getDeliveryQuery   = `SELECT id FROM deliveries WHERE name=$1 AND phone=$2 AND zip=$3 AND city=$4 AND address=$5 AND region=$6 AND email=$7`
	insertPaymentQuery = `INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id;`
	insertOrderQuery = `INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
					 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);`
	insertItemQuery = `INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;`
	insertOrdersItemsQuery = `INSERT INTO orders_x_items (order_uid, item_id)
							VALUES ($1, $2);`
)

func (r *OrderRepository) GetOrder(orderUID string) (*model.Order, error) {
	var dbOrd dbOrder

	tx, err := r.db.Beginx()
	if err != nil {
		return &model.Order{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			return
		}
	}()

	err = tx.Get(&dbOrd, getOrderQuery+" WHERE o.order_uid = $1", orderUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &model.Order{}, fmt.Errorf("order %s not found: %w", orderUID, err)
		}
		return &model.Order{}, fmt.Errorf("failed to get order %s: %w", orderUID, err)
	}

	order := model.Order{
		OrderUID:          dbOrd.OrderUID,
		TrackNumber:       dbOrd.TrackNumber,
		Entry:             dbOrd.Entry,
		Locale:            dbOrd.Locale,
		InternalSignature: dbOrd.InternalSignature,
		CustomerID:        dbOrd.CustomerID,
		DeliveryService:   dbOrd.DeliveryService,
		Shardkey:          dbOrd.Shardkey,
		SmID:              dbOrd.SmID,
		DateCreated:       dbOrd.DateCreated,
		OofShard:          dbOrd.OofShard,
		Delivery: model.Delivery{
			Name:    dbOrd.DeliveryName,
			Phone:   dbOrd.DeliveryPhone,
			Zip:     dbOrd.DeliveryZip,
			City:    dbOrd.DeliveryCity,
			Address: dbOrd.DeliveryAddress,
			Region:  dbOrd.DeliveryRegion,
			Email:   dbOrd.DeliveryEmail,
		},
		Payment: model.Payment{
			Transaction:  dbOrd.PaymentTransaction,
			RequestID:    dbOrd.PaymentRequestID,
			Currency:     dbOrd.PaymentCurrency,
			Provider:     dbOrd.PaymentProvider,
			Amount:       dbOrd.PaymentAmount,
			PaymentDt:    dbOrd.PaymentPaymentDt,
			Bank:         dbOrd.PaymentBank,
			DeliveryCost: dbOrd.PaymentDeliveryCost,
			GoodsTotal:   dbOrd.PaymentGoodsTotal,
			CustomFee:    dbOrd.PaymentCustomFee,
		},
	}

	var items []model.Item
	err = tx.Select(&items, getItemsQuery, orderUID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &model.Order{}, fmt.Errorf("failed to get items: %w", err)
	}
	order.Items = items

	if err := tx.Commit(); err != nil {
		return &model.Order{}, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &order, nil
}

func (r *OrderRepository) SaveOrder(order *model.Order) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			return
		}
	}()

	var deliveryID int64
	err = tx.Get(&deliveryID, insertDeliveryQuery,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email)
	if err != nil {
		err = tx.Get(&deliveryID, getDeliveryQuery,
			order.Delivery.Name,
			order.Delivery.Phone,
			order.Delivery.Zip,
			order.Delivery.City,
			order.Delivery.Address,
			order.Delivery.Region,
			order.Delivery.Email)
		if err != nil {
			return fmt.Errorf("failed to insert new delivery: %w", err)
		}
	} else {
		log.Println("new delivery inserted")
	}

	var paymentID int
	err = tx.Get(&paymentID, insertPaymentQuery,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee)

	if err != nil {
		return fmt.Errorf("failed to insert new payment: %w", err)
	}

	log.Printf("delivery_id=%d", deliveryID)
	log.Printf("payment_id=%d", paymentID)

	_, err = tx.Exec(insertOrderQuery,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		deliveryID,
		paymentID,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard)

	if err != nil {
		return fmt.Errorf("failed to insert new order: %w", err)
	}

	for _, item := range order.Items {
		var itemID int64
		err := tx.Get(&itemID, insertItemQuery,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err != nil {
			return fmt.Errorf("failed to insert new item: %w", err)
		}

		_, err = tx.Exec(insertOrdersItemsQuery, order.OrderUID, itemID)
		if err != nil {
			return fmt.Errorf("failed to insert new item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *OrderRepository) GetAllOrders(limit int64) ([]*model.Order, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			return
		}
	}()

	var dbOrds []dbOrder
	err = tx.Select(&dbOrds, getOrderQuery+" ORDER BY o.date_created DESC LIMIT $1", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	var orders []*model.Order
	for _, dbOrd := range dbOrds {
		order := &model.Order{
			OrderUID:          dbOrd.OrderUID,
			TrackNumber:       dbOrd.TrackNumber,
			Entry:             dbOrd.Entry,
			Locale:            dbOrd.Locale,
			InternalSignature: dbOrd.InternalSignature,
			CustomerID:        dbOrd.CustomerID,
			DeliveryService:   dbOrd.DeliveryService,
			Shardkey:          dbOrd.Shardkey,
			SmID:              dbOrd.SmID,
			DateCreated:       dbOrd.DateCreated,
			OofShard:          dbOrd.OofShard,
			Delivery: model.Delivery{
				Name:    dbOrd.DeliveryName,
				Phone:   dbOrd.DeliveryPhone,
				Zip:     dbOrd.DeliveryZip,
				City:    dbOrd.DeliveryCity,
				Address: dbOrd.DeliveryAddress,
				Region:  dbOrd.DeliveryRegion,
				Email:   dbOrd.DeliveryEmail,
			},
			Payment: model.Payment{
				Transaction:  dbOrd.PaymentTransaction,
				RequestID:    dbOrd.PaymentRequestID,
				Currency:     dbOrd.PaymentCurrency,
				Provider:     dbOrd.PaymentProvider,
				Amount:       dbOrd.PaymentAmount,
				PaymentDt:    dbOrd.PaymentPaymentDt,
				Bank:         dbOrd.PaymentBank,
				DeliveryCost: dbOrd.PaymentDeliveryCost,
				GoodsTotal:   dbOrd.PaymentGoodsTotal,
				CustomFee:    dbOrd.PaymentCustomFee,
			},
		}

		var items []model.Item
		err = tx.Select(&items,
			getItemsQuery,
			dbOrd.OrderUID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get items for order %s: %w", dbOrd.OrderUID, err)
		}
		order.Items = items

		orders = append(orders, order)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return orders, nil
}

package database

import (
	"context"
	"fmt"
	"main/internal/log"
	"main/internal/model"
)

// GetItems performs SQL query which gets data from Items table by provided orderUID.
func (db Postgres) GetItems(orderUID string) ([]model.Items, error) {
	rows, err := db.pool.Query(context.Background(),
		`SELECT i.chrt_id, i.track_number, i.price, i.rid, i.name,
      			i.sale, i.size, i.total_price, i.nm_id, i.brand,
				i.status
			FROM order_items JOIN items i on i.chrt_id = order_items.chrt_id
			WHERE order_items.order_uid = $1`, orderUID)

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	items := make([]model.Items, 0)
	for rows.Next() {
		item := model.Items{}

		rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price,
			&item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID,
			&item.Brand, &item.Status)

		if rows.Err() != nil {
			return nil, rows.Err()
		}

		items = append(items, item)
	}

	return items, nil
}

// GetRecords performs SQL query which gets all records from Postgres database.
func (db Postgres) GetRecords() []*model.RecordModel {
	query := `SELECT wb_order.order_uid, wb_order.track_number, wb_order.entry,
					   d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
					   p.transaction, p.request_id, p.currency, p.provider, p.amount,
					   p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee,
					   wb_order.locale, wb_order.internal_signature, wb_order.customer_id, wb_order.delivery_service,
					   wb_order.shardkey, wb_order.sm_id, wb_order.date_created, wb_order.oof_shard
					FROM wb_order
					JOIN delivery d on d.id = wb_order.delivery_id
					JOIN payment p on p.id = wb_order.payment_id`

	rows, err := db.pool.Query(context.Background(), query)
	if err != nil {
		log.Logger.Fatalln(fmt.Errorf("GetRecords: Query: %w", err))
	}
	defer rows.Close()

	records := make([]*model.RecordModel, 0)
	for rows.Next() {
		record := &model.RecordModel{}
		err = rows.Scan(&record.OrderUID, &record.TrackNumber, &record.Entry,
			&record.Delivery.Name, &record.Delivery.Phone, &record.Delivery.Zip, &record.Delivery.City,
			&record.Delivery.Address, &record.Delivery.Region, &record.Delivery.Email, &record.Payment.Transaction,
			&record.Payment.RequestID, &record.Payment.Currency, &record.Payment.Provider, &record.Payment.Amount,
			&record.Payment.PaymentDt, &record.Payment.Bank, &record.Payment.DeliveryCost, &record.Payment.GoodsTotal,
			&record.Payment.CustomFee, &record.Locale, &record.InternalSignature, &record.CustomerID,
			&record.DeliveryService, &record.Shardkey, &record.SmID, &record.DateCreated, &record.OofShard)

		if err != nil {
			log.Logger.Fatalln(rows.Err())
		}

		record.Items, err = db.GetItems(record.OrderUID)
		if err != nil {
			log.Logger.Fatalln(fmt.Errorf("GetItems: %w", err))
		}
		records = append(records, record)
	}

	return records
}

// CheckDuplicate checks if orderUID exists in database.
func (db Postgres) CheckDuplicate(orderUID string) (bool, error) {
	tag, err := db.pool.Exec(context.Background(),
		`SELECT order_uid FROM wb_order WHERE order_uid = $1`, orderUID)

	if err != nil {
		return false, fmt.Errorf("exec select: %w", err)
	} else if tag.RowsAffected() == 1 {
		return true, nil
	}

	return false, nil
}

// AddRecord loads model into database.
func (db Postgres) AddRecord(record model.RecordModel) error {
	exist, err := db.CheckDuplicate(record.OrderUID)

	if err != nil {
		return fmt.Errorf("CheckDuplicate: %w", err)
	} else if exist {
		log.Logger.Traceln("AddRecord: message is not unique")
		return nil
	}

	delKey, err := db.addDelivery(record.Delivery)
	if err != nil {
		return fmt.Errorf("AddDelivery: %w", err)
	}

	payKey, err := db.addPayment(record.Payment)
	if err != nil {
		return fmt.Errorf("AddPayment: %w", err)
	}

	err = db.addItems(record.Items)
	if err != nil {
		return fmt.Errorf("AddItems: %w", err)
	}

	err = db.addWbOrder(payKey, delKey, record)
	if err != nil {
		return fmt.Errorf("AddWbOrder: %w", err)
	}

	err = db.addOrderItems(record)
	if err != nil {
		return fmt.Errorf("AddOrderItems: %w", err)
	}

	return nil
}

func (db Postgres) addDelivery(delivery model.Delivery) (int, error) {
	row := db.pool.QueryRow(context.Background(),
		`INSERT INTO delivery(name, phone, zip, city, address, region, email)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id;`,
		delivery.Name, delivery.Phone, delivery.Zip, delivery.City,
		delivery.Address, delivery.Region, delivery.Email)

	var delKey int
	if err := row.Scan(&delKey); err != nil {
		return 0, fmt.Errorf("row.Scan: %w", err)
	}

	return delKey, nil
}

func (db Postgres) addPayment(payment model.Payment) (int, error) {
	row := db.pool.QueryRow(context.Background(),
		`INSERT INTO payment(transaction,  request_id, currency, provider, amount, 
                    payment_dt, bank, delivery_cost, goods_total, custom_fee)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id;`,
		payment.Transaction, payment.RequestID, payment.Currency, payment.Provider,
		payment.Amount, payment.PaymentDt, payment.Bank, payment.DeliveryCost,
		payment.GoodsTotal, payment.CustomFee)

	var payKey int
	if err := row.Scan(&payKey); err != nil {
		return 0, fmt.Errorf("row.Scan: %w", err)
	}

	return payKey, nil
}

func (db Postgres) addItems(items []model.Items) error {
	for i := range items {
		_, err := db.pool.Exec(context.Background(),
			`INSERT INTO items (chrt_id, track_number, price, rid,
			name, sale, size, total_price, nm_id, brand, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			items[i].ChrtID, items[i].TrackNumber, items[i].Price, items[i].Rid,
			items[i].Name, items[i].Sale, items[i].Size, items[i].TotalPrice,
			items[i].NmID, items[i].Brand, items[i].Status)

		if err != nil {
			return fmt.Errorf("exec: %w", err)
		}
	}
	return nil
}

func (db Postgres) addWbOrder(payKey, delKey int, order model.RecordModel) error {
	_, err := db.pool.Exec(context.Background(),
		`INSERT INTO wb_order (order_uid, track_number ,entry, delivery_id,
			payment_id , locale, internal_signature,
			customer_id, delivery_service, shardkey,
			sm_id, date_created ,oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		order.OrderUID, order.TrackNumber, order.Entry, delKey, payKey,
		order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.Shardkey, order.SmID, order.DateCreated, order.OofShard)

	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

func (db Postgres) addOrderItems(model model.RecordModel) error {
	for i := range model.Items {
		_, err := db.pool.Exec(context.Background(),
			`INSERT INTO order_items (order_uid, chrt_id)
				VALUES ($1, $2)`, model.OrderUID, model.Items[i].ChrtID)

		if err != nil {
			return fmt.Errorf("exec: %w", err)
		}
	}

	return nil
}

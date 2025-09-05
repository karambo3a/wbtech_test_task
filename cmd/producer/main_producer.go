package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(os.Getenv("KAFKA_BROKERS_PROD")),
		Topic:                  "order",
		AllowAutoTopicCreation: true,
	}

	var err error
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		err = w.WriteMessages(ctx, kafka.Message{
			Key:   []byte("order"),
			Value: []byte(generateOrder()),
		})
		cancel()

		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err != nil {
			log.Fatalf("unexpected error %v", err)
		}

		time.Sleep(time.Second * 2)
	}
}

func generateOrder() string {
	timestamp := time.Now().Unix()
	randomNum := rand.Intn(10000)
	orderUID := fmt.Sprintf("order_%d_%04d", timestamp, randomNum)

	template := `{
		"order_uid": "%s",
		"track_number": "WBILTRACK123456",
		"entry": "WBIL",
		"delivery": {
			"name": "Иван Петров",
			"phone": "+74951234567",
			"zip": "101000",
			"city": "Москва",
			"address": "ул. Тверская, д. 25",
			"region": "Московская область",
			"email": "ivan.petrov@mail.ru"
		},
		"payment": {
			"transaction": "%s",
			"request_id": "req_987654",
			"currency": "RUB",
			"provider": "sberpay",
			"amount": 5420,
			"payment_dt": 1672534891,
			"bank": "sberbank",
			"delivery_cost": 500,
			"goods_total": 4920,
			"custom_fee": 0
		},
		"items": [
			{
				"chrt_id": 87654321,
				"track_number": "WBILTRACK123456",
				"price": 2460,
				"rid": "cd5632a1e8947f12prod",
				"name": "Кроссовки спортивные",
				"sale": 15,
				"size": "42",
				"total_price": 2091,
				"nm_id": 6543210,
				"brand": "Nike",
				"status": 200
			}
		],
		"locale": "ru",
		"internal_signature": "signature_123",
		"customer_id": "customer_456",
		"delivery_service": "cdek",
		"shardkey": "5",
		"sm_id": 42,
		"date_created": "2023-01-15T14:28:31Z",
		"oof_shard": "2"
	}`

	return fmt.Sprintf(template, orderUID, orderUID)
}

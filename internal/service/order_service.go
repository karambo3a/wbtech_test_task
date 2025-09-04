package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/karambo3a/wbtech_test_task/internal/cache"
	"github.com/karambo3a/wbtech_test_task/internal/consumer"
	"github.com/karambo3a/wbtech_test_task/internal/model"
	"github.com/karambo3a/wbtech_test_task/internal/repository"
	"github.com/redis/go-redis/v9"
)

type OrderService struct {
	repository *repository.Repository
	consumer   consumer.Consumer
	cache      cache.RedisCache
}

func NewOrderService(repository *repository.Repository, consumer consumer.Consumer, cache cache.RedisCache) *OrderService {
	service := &OrderService{
		repository: repository,
		consumer:   consumer,
		cache: cache,
	}

	if err := service.cache.Init(repository); err != nil {
		log.Fatalln("failed to get cache from db")
	}

	service.consumer.StartConsuming(service.SaveOrder)
	return service
}

func (s *OrderService) SaveOrder(msg []byte) error {
	var order model.Order
	err := json.Unmarshal(msg, &order)
	if err != nil {
		return fmt.Errorf("failed to parse order json: %w", err)
	}

	err = s.repository.SaveOrder(order)
	if err != nil {
		return fmt.Errorf("failed to save new order: %w", err)
	}

	if err = s.cache.Set(context.Background(), order.OrderUID, order, 24*time.Hour); err != nil {
		return fmt.Errorf("failed to save in cache order_uid=%s: %w", order.OrderUID, err)
	}
	log.Println("order saved")
	return nil
}

func (s *OrderService) GetOrder(orderUID string) (model.Order, error) {
	order, err := s.cache.Get(context.Background(), orderUID)
	if err == nil {
		return order, nil
	}

	if errors.Is(err, redis.Nil) {
		order, err = s.repository.GetOrder(orderUID)
		if err != nil {
			return model.Order{}, fmt.Errorf("order with order_uid=%s is not found: %w", orderUID, err)
		}

		go s.cacheOrderAsync(orderUID, order)

		log.Println("got order by id")
		return order, nil
	}

	return model.Order{}, fmt.Errorf("failed to get order_uid=%s: %w", orderUID, err)
}

func (s *OrderService) cacheOrderAsync(orderUID string, order model.Order) {
	if err := s.cache.Set(context.Background(), orderUID, order, 24*time.Hour); err != nil {
		log.Printf("failed to save in cache order_uid=%s: %v", orderUID, err)
	}
}

func (s *OrderService) CloseConsumer() {
	s.consumer.Close()
}

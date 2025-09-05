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

func NewOrderService(repository *repository.Repository, consumer consumer.Consumer, cache cache.RedisCache, initLimit int64) *OrderService {
	service := &OrderService{
		repository: repository,
		consumer:   consumer,
		cache:      cache,
	}

	orders, err := repository.GetAllOrders(initLimit)
	if err != nil {
		log.Fatalf("failed to get orders from database to cache: %v", err)
	}

	if err := service.cache.Init(orders); err != nil {
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

	err = s.repository.SaveOrder(&order)
	if err != nil {
		return fmt.Errorf("failed to save new order: %w", err)
	}

	go s.cacheOrderAsync(order.OrderUID, &order)

	log.Println("order saved")
	return nil
}

func (s *OrderService) GetOrder(orderUID string) (*model.Order, error) {
	order, err := s.cache.Get(context.TODO(), orderUID)
	if err == nil {
		return order, nil
	}

	if !errors.Is(err, redis.Nil) {
		return &model.Order{}, fmt.Errorf("failed to get order_uid=%s: %w", orderUID, err)
	}

	order, err = s.repository.GetOrder(orderUID)
	if err != nil {
		return &model.Order{}, fmt.Errorf("order with order_uid=%s is not found: %w", orderUID, err)
	}

	orderAsync := *order
	go s.cacheOrderAsync(orderUID, &orderAsync)

	log.Println("got order by id")
	return order, nil
}

func (s *OrderService) cacheOrderAsync(orderUID string, order *model.Order) {
	if err := s.cache.Set(context.TODO(), orderUID, order, 24*time.Hour); err != nil {
		log.Printf("failed to save in cache order_uid=%s: %v", orderUID, err)
	}
}

func (s *OrderService) CloseConsumer() {
	s.consumer.Close()
}

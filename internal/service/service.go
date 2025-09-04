package service

import (
	"github.com/karambo3a/wbtech_test_task/internal/cache"
	"github.com/karambo3a/wbtech_test_task/internal/consumer"
	"github.com/karambo3a/wbtech_test_task/internal/model"
	"github.com/karambo3a/wbtech_test_task/internal/repository"
)

//go:generate mockgen -source=service.go -destination=../../test/mocks/service_mock.go

type OrderServiceInterface interface {
	SaveOrder(msg []byte) error
	GetOrder(orderUID string) (model.Order, error)
	CloseConsumer()
}

type Service struct {
	OrderServiceInterface
}

func NewService(repository *repository.Repository, consumer consumer.Consumer, cache cache.RedisCache) *Service {
	return &Service{
		OrderServiceInterface: NewOrderService(repository, consumer, cache),
	}
}

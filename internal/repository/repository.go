package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/karambo3a/wbtech_test_task/internal/model"
)

//go:generate mockgen -source=repository.go -destination=../../test/mocks/repository_mock.go

type OrderRepositoryInterface interface {
	GetOrder(orderUID string) (model.Order, error)
	SaveOrder(order model.Order) error
	GetAllOrders() ([]model.Order, error)
}

type Repository struct {
	OrderRepositoryInterface
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		OrderRepositoryInterface: NewOrderRepository(db),
	}
}

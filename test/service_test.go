package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/karambo3a/wbtech_test_task/internal/model"
	"github.com/karambo3a/wbtech_test_task/internal/repository"
	"github.com/karambo3a/wbtech_test_task/internal/service"
	mock "github.com/karambo3a/wbtech_test_task/test/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestServiceGetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetOrderRepository := mock.NewMockOrderRepositoryInterface(ctrl)
	mockRepository := &repository.Repository{OrderRepositoryInterface: mockGetOrderRepository}
	mockConsumer := mock.NewMockConsumer(ctrl)
	mockRedisCache := mock.NewMockRedisCache(ctrl)

	mockRedisCache.EXPECT().Init(mockRepository).Return(nil)
	mockConsumer.EXPECT().StartConsuming(gomock.Any()).AnyTimes()
	s := service.NewService(mockRepository, mockConsumer, mockRedisCache, int64(100))

	testOrder := model.Order{
		OrderUID:          "order_uid1",
		TrackNumber:       "a",
		Entry:             "a",
		Locale:            "a",
		InternalSignature: "",
		CustomerID:        "a",
		DeliveryService:   "a",
		Shardkey:          "1",
		SmID:              1,
		DateCreated:       time.Date(2021, 11, 26, 6, 22, 19, 0, time.UTC),
		OofShard:          "1",
		Delivery: model.Delivery{
			Name:    "a a",
			Phone:   "+9720000000",
			Zip:     "1",
			City:    "a",
			Address: "a",
			Region:  "a",
			Email:   "a@a.a",
		},
		Payment: model.Payment{
			Transaction:  "12",
			RequestID:    "",
			Currency:     "a",
			Provider:     "a",
			Amount:       1,
			PaymentDt:    1,
			Bank:         "a",
			DeliveryCost: 1,
			GoodsTotal:   1,
			CustomFee:    1,
		},
		Items: []model.Item{
			{
				ChrtID:      1,
				TrackNumber: "a",
				Price:       1,
				Rid:         "a",
				Name:        "a",
				Sale:        1,
				Size:        "1",
				TotalPrice:  1,
				NmID:        1,
				Brand:       "a",
				Status:      1,
			},
		},
	}

	tests := []struct {
		name          string
		orderUID      string
		mockBehavior  func()
		expectedOrder model.Order
		expectedErr   error
	}{
		{
			name:     "order in cache successful",
			orderUID: "order_uid1",
			mockBehavior: func() {
				mockRedisCache.EXPECT().Get(context.Background(), "order_uid1").Return(testOrder, nil)
			},
			expectedOrder: testOrder,
			expectedErr:   nil,
		},
		{
			name:     "order not in cache, found in repository",
			orderUID: "order_uid1",
			mockBehavior: func() {
				mockRedisCache.EXPECT().Get(context.Background(), "order_uid1").Return(model.Order{}, redis.Nil)
				mockGetOrderRepository.EXPECT().GetOrder("order_uid1").Return(testOrder, nil)
				mockRedisCache.EXPECT().Set(context.Background(), "order_uid1", testOrder, 24*time.Hour).Return(nil)
			},
			expectedOrder: testOrder,
			expectedErr:   nil,
		},
		{
			name:     "order not in cache and not in repository",
			orderUID: "order_uid1",
			mockBehavior: func() {
				mockRedisCache.EXPECT().Get(context.Background(), "order_uid1").Return(model.Order{}, redis.Nil)
				mockGetOrderRepository.EXPECT().GetOrder("order_uid1").Return(model.Order{}, errors.New("failed to get order"))
			},
			expectedOrder: model.Order{},
			expectedErr:   errors.New("order with order_uid=order_uid1 is not found: failed to get order"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior()

			order, err := s.GetOrder(test.orderUID)

			assert.Equal(t, test.expectedOrder, order)
			if err != nil {
				assert.Equal(t, test.expectedErr.Error(), err.Error())
			}
		})
	}
}

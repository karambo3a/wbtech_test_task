package test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/karambo3a/wbtech_test_task/internal/handlers"
	"github.com/karambo3a/wbtech_test_task/internal/model"
	"github.com/karambo3a/wbtech_test_task/internal/service"
	mock "github.com/karambo3a/wbtech_test_task/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandlerGetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetOrderService := mock.NewMockOrderServiceInterface(ctrl)
	mockService := &service.Service{OrderServiceInterface: mockGetOrderService}
	h := handlers.NewHandler(mockService)

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
	expectedJSON := `{"order_uid":"order_uid1","track_number":"a","entry":"a","delivery":{"name":"a a","phone":"+9720000000","zip":"1","city":"a","address":"a","region":"a","email":"a@a.a"},"payment":{"transaction":"12","request_id":"","currency":"a","provider":"a","amount":1,"payment_dt":1,"bank":"a","delivery_cost":1,"goods_total":1,"custom_fee":1},"items":[{"chrt_id":1,"track_number":"a","price":1,"rid":"a","name":"a","sale":1,"size":"1","total_price":1,"nm_id":1,"brand":"a","status":1}],"locale":"a","internal_signature":"","customer_id":"a","delivery_service":"a","shardkey":"1","sm_id":1,"date_created":"2021-11-26T06:22:19Z","oof_shard":"1"}`

	tests := []struct {
		name           string
		orderUID       string
		mockBehavior   func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "success",
			orderUID: "order_uid1",
			mockBehavior: func() {
				mockGetOrderService.EXPECT().GetOrder("order_uid1").Return(testOrder, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   expectedJSON,
		},
		{
			name:     "not found",
			orderUID: "order_not_found",
			mockBehavior: func() {
				mockGetOrderService.EXPECT().GetOrder("order_not_found").Return(model.Order{}, errors.New("order not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"order not found"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior()

			r := chi.NewRouter()
			r.Get("/order/{order_uid}", h.GetOrder)

			body := ""
			req, err := http.NewRequest(http.MethodGet, "/order/"+test.orderUID, strings.NewReader(body))
			if err != nil {
				t.Fatalf("error occurred while testing")
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatus, w.Code)
			assert.JSONEq(t, test.expectedBody, w.Body.String())
		})
	}
}

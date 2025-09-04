package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/karambo3a/wbtech_test_task/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) InitRouts() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/order/{order_uid}", h.GetOrder)
	return r
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderUID := chi.URLParam(r, "order_uid")
	log.Printf("order_uid=%s\n", orderUID)

	w.Header().Set("Content-Type", "application/json")
	if orderUID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "empty order_uid"})
		return
	}

	order, err := h.service.GetOrder(orderUID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// to get response from client
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err = json.NewEncoder(w).Encode(order); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	log.Printf("order with order_uid=%s sent\n", orderUID)
}

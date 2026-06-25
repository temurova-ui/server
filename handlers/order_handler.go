package handlers

import (
	"encoding/json"
	"net/http"
	
	internalCtx "myservice/internal/context"
	"myservice/internal/models"
	"myservice/internal/service"
	"myservice/pkg/logger"
)

type OrderHandler struct {
	svc *service.OrderService
	log *logger.Logger
}

func NewOrderHandler(svc *service.OrderService, log *logger.Logger) *OrderHandler {
	return &OrderHandler{svc: svc, log: log}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(internalCtx.UserIDKey).(int)
	if !ok {
		h.log.Error("User ID missing from context in Order Create")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json format", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := h.svc.CreateOrder(r.Context(), userID, req)
	if err != nil {
		h.log.Error("Failed to create order: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}
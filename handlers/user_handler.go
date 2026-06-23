package handlers

import (
	"encoding/json"
	"errors"
	"myservice/internal/service"
	"myservice/pkg/errs"
	"myservice/pkg/logger"
	"net/http"
)

type UserHandler struct{
	svc *service.UserService
	log *logger.Logger
}

func NewUserHandler(svc *service.UserService, log *logger.Logger) *UserHandler{
	return &UserHandler{svc: svc, log: log}
}

type RegisterRequest struct{
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler)Register (w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		h.log.Error("Failed to decode JSON: %v", err)
		http.Error(w, "invalid json format", http.StatusBadRequest)
		return
	}

	user, err := h.svc.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, errs.ErrEmailConflict):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			h.log.Error("Database error: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

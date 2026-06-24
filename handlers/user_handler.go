package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"myservice/internal/models"
	"myservice/internal/service"
	"myservice/pkg/errs"
	"myservice/pkg/logger"
)

type UserHandler struct{
	svc *service.UserService
	log *logger.Logger
}

func NewUserHandler(svc *service.UserService, log *logger.Logger) *UserHandler{
	return &UserHandler{svc: svc, log: log}
}

func (h *UserHandler)Register (w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, "invalid json format", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.svc.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, errs.ErrEmailConflict) {
			http.Error(w, "email already registered", http.StatusConflict)
		}else{
			h.log.Error("Registration database failure: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.AuthResponse{Token: token})
}

func (h *UserHandler)Login(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, "invalid json format", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.svc.Login(r.Context(), req)
	if err != nil{
		http.Error(w, "unauthorized: invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Context-Type", "application/json")
	json.NewEncoder(w).Encode(models.AuthResponse{Token: token})
}
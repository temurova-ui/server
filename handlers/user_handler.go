package handlers

import (
	"encoding/json"
	"errors"
	// "fmt"
	// "log"
	internalCtx "myservice/internal/context"
	"myservice/internal/models"
	"myservice/internal/service"
	"myservice/pkg/errs"
	"myservice/pkg/logger"
	"net/http"
)

type UserHandler struct{
	svc service.UserService
	log *logger.Logger
}

func NewUserHandler(svc service.UserService, log *logger.Logger) *UserHandler{
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
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AuthResponse{Token: token})
}

func (h *UserHandler)GetMe(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(internalCtx.UserIDKey).(int)
	if !ok {
		h.log.Error("User ID not found in context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.svc.GetProfile(r.Context(), userID)
	if err != nil {
		h.log.Error("Failed to get user profile: %v", err)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler)UpdateMe(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPut{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(internalCtx.UserIDKey).(int)
	if !ok {
		h.log.Error("User ID missing from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, "invalid json format", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := h.svc.UpdateProfile(r.Context(), userID, req); err != nil{
		h.log.Error("Failed to update user profile: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "profile updated successfully"})
}

func (h *UserHandler)DeleteMe(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodDelete{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(internalCtx.UserIDKey).(int)
	if !ok {
		h.log.Error("User ID missing from context in DeleteMe")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	
	if err := h.svc.DeleteAccount(r.Context(), userID); err != nil{
		h.log.Error("Failed to delete user profile: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "account successfully deleted"})

}

func (h *UserHandler)ChangePassword(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPut{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(internalCtx.UserIDKey).(int)
	if !ok {
		h.log.Error("User ID missing from context in ChangePassword")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.ChangePassword
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, "invalid json format", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.svc.ChangePassword(r.Context(), userID, req); err != nil{
		h.log.Error("Failed to change password: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "password successfully changed"})
}
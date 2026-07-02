package handlers

import (
	"encoding/json"
	"net/http"
	internalCtx "myservice/internal/context"
)


func (h *UserHandler)GetUserLoginHistory(w http.ResponseWriter, r *http.Request){
	userID, ok := r.Context().Value(internalCtx.UserIDKey).(int64)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	histories, err := h.svc.GetLoginHistory(r.Context(), userID)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(histories)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
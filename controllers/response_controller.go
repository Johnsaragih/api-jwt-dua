package controllers

import (
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := map[string]interface{}{
		"status":  status,
		"message": message,
		"data":    data,
	}
	json.NewEncoder(w).Encode(response)
}

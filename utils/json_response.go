package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

//Buat Fungsi Respon supaya urut menggunakan struct bukan Map

type JSONResponseFormat struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,,omitempty"`
	Token   string      `json:"token,omitempty"`
}

func JSONResponse(w http.ResponseWriter, status int, message string, data interface{}, token string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := JSONResponseFormat{
		Status:  status,
		Message: message,
		Data:    data,
		Token:   token,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("JSON encode Error :", err)
	}
}

/*
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
*/

package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string  `json:"status"`
	Message *string `json:"message"`
	Data    any     `json:"data"`
}

func Success(w http.ResponseWriter, data any) {
	response := &Response{
		Status:  "success",
		Message: nil,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	json.NewEncoder(w).Encode(response)
}

func Error(w http.ResponseWriter, code int, message string) {
	response := &Response{
		Status:  "error",
		Message: &message,
		Data:    nil,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(response)
}

package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func InternalError(w http.ResponseWriter, err error) {
	log.Printf("internal error: %v", err)
	http.Error(w, "internal error", http.StatusInternalServerError)
}

func Unauthorized(w http.ResponseWriter) {
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}

func BadRequest(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusBadRequest)
}

package handler

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, code int, message []string) {
	WriteJSON(w, code, map[string][]string{"error": message})
}

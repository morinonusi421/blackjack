package handlers

import (
	"encoding/json"
	"net/http"
)

// HelloHandler は "Hello, World!" を返します。
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello, World!"})
}

package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthResponse はヘルスチェックのレスポンス構造体
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthHandler はシンプルなヘルスチェックエンドポイントを提供します。
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
}

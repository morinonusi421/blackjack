package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockRandomNumberGenerator はテスト用に固定値を返す実装です。
type mockRandomNumberGenerator struct {
	value int
}

func (m mockRandomNumberGenerator) Generate() int {
	return m.value
}

func TestRandomNumberHandler_ReturnsCorrectJSON(t *testing.T) {
	// 準備: 期待する値を返すモックとハンドラ生成
	expectedNumber := 7
	generator := mockRandomNumberGenerator{value: expectedNumber}
	handler := NewRandomNumberHandler(generator)

	// リクエストとレスポンスレコーダを用意
	req := httptest.NewRequest(http.MethodGet, "/api/random_number", nil)
	rr := httptest.NewRecorder()

	// ハンドラ呼び出し
	handler.ServeHTTP(rr, req)

	// ステータスコードの確認
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Content-Type ヘッダの確認
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", got)
	}

	// レスポンスボディの確認
	var resp RandomNumberResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Number != expectedNumber {
		t.Fatalf("expected number %d, got %d", expectedNumber, resp.Number)
	}
}

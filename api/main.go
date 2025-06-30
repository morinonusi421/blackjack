package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// CORSを許可するためのミドルウェア
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 許可するオリジンを指定します。'*'は全てのオリジンを許可します。
		// 本番環境では特定のドメイン（Next.jsアプリのURL）を指定することが推奨されます。
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// OPTIONSリクエストの場合はここで処理を終了
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// レスポンスのContent-TypeをJSONに設定
	w.Header().Set("Content-Type", "application/json")

	// 返却するデータを作成
	response := map[string]string{"message": "Hello, World!"}

	// JSONにエンコードしてレスポンスとして書き出す
	json.NewEncoder(w).Encode(response)
}

func main() {
	// render.comが設定するPORT環境変数を取得。なければ8080を使う
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ハンドラを登録
	mux := http.NewServeMux()
	mux.HandleFunc("/api/hello", helloHandler)

	// ミドルウェアを適用したハンドラ
	handlerWithCors := corsMiddleware(mux)

	log.Println("Server starting on port " + port)
	// サーバーを起動
	if err := http.ListenAndServe(":"+port, handlerWithCors); err != nil {
		log.Fatal(err)
	}
}

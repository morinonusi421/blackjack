package main

import (
	"log"
	"net/http"
	"os"

	"blackjack/api/game"
	"blackjack/api/handlers"
	"blackjack/api/services"

	"github.com/gorilla/mux"
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

func main() {
	// render.comが設定するPORT環境変数を取得。なければ8080を使う
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ルーターを作成
	router := mux.NewRouter()

	// 依存性の生成
	gameService := services.NewGameService(&game.RandomDeck{})
	strategyService := services.NewStrategyService()

	// ゲームエンドポイント
	router.HandleFunc("/api/game/new", handlers.NewGameHandler(gameService)).Methods("POST")

	// ヒットエンドポイント
	router.HandleFunc("/api/game/hit", handlers.HitHandler(gameService)).Methods("POST")

	// スタンドエンドポイント
	router.HandleFunc("/api/game/stand", handlers.StandHandler(gameService)).Methods("POST")

	// サレンダーエンドポイント
	router.HandleFunc("/api/game/surrender", handlers.SurrenderHandler(gameService)).Methods("POST")

	// ヘルスチェックエンドポイント
	router.HandleFunc("/api/health", handlers.HealthHandler).Methods("GET")

	// 戦略アドバイスエンドポイント
	router.HandleFunc("/api/strategy/advise", handlers.StrategyHandler(strategyService)).Methods("POST")

	// ミドルウェアを適用したハンドラ
	handlerWithCors := corsMiddleware(router)

	log.Println("Server starting on port " + port)
	// サーバーを起動
	if err := http.ListenAndServe(":"+port, handlerWithCors); err != nil {
		log.Fatal(err)
	}
}

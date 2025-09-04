package game

// 結果メッセージ定数（日本語）
const (
	// ブラックジャック（初手21）でプレイヤー勝利
	MessageBlackjackPlayerWin = "ブラックジャック！(初手が21の特殊勝利)プレイヤーの勝ちです"

	// プレイヤーがバーストしてディーラー勝利
	MessagePlayerBustDealerWin = "プレイヤーがバースト！ディーラーの勝ちです"

	// ディーラーがバーストしてプレイヤー勝利
	MessageDealerBustPlayerWin = "ディーラーがバースト！プレイヤーの勝ちです"

	// 通常の勝敗（非ブラックジャック・非バースト）
	MessagePlayerWin = "プレイヤーの勝ちです"
	MessageDealerWin = "ディーラーの勝ちです"

	// 引き分け（プッシュ）
	MessagePush = "引き分けです"

	// サレンダー
	MessagePlayerSurrendered = "プレイヤーがサレンダーしました。"
)

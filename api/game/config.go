package game

// GameConfig はゲーム全体の設定を表します。
type GameConfig struct {
	DealerStandThreshold int `json:"dealer_stand_threshold"` // ディーラーがスタンドする閾値
}


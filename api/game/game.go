package game

type Suit string

const (
	Spade   Suit = "Spade"
	Heart   Suit = "Heart"
	Diamond Suit = "Diamond"
	Club    Suit = "Club"
)

type Rank string

// 利便性のために配列でも保持しておく。
var ranks = []Rank{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
var suits = []Suit{Spade, Heart, Diamond, Club}

// Card は 1 枚のカードを表します。
type Card struct {
	Suit Suit `json:"suit"`
	Rank Rank `json:"rank"`
}

// Hand は手札を表します。
type Hand struct {
	Cards []Card `json:"cards"`
	Score int    `json:"score"`
}

// GameState はゲームの進行状況を表します。
type GameState string

const (
	PlayerTurn GameState = "PlayerTurn"
	Finished   GameState = "Finished"
)

// Result はゲームの結果を表します。
type Result string

const (
	Pending   Result = "Pending"
	PlayerWin Result = "PlayerWin"
	DealerWin Result = "DealerWin"
	Push      Result = "Push"
	Surrender Result = "Surrender"
)

// Game はゲーム全体の状態を保持します。
type Game struct {
	PlayerHand    Hand      `json:"player_hand"`
	DealerHand    Hand      `json:"dealer_hand"`
	State         GameState `json:"state"`
	Result        Result    `json:"result"`
	ResultMessage string    `json:"result_message"`
	Bet           int       `json:"bet"`    // 掛け金
	Payout        int       `json:"payout"` // 払戻金（勝利額／Push はベット返却）
}

var rankScoreMap = map[Rank]int{
	"A":  11,
	"K":  10,
	"Q":  10,
	"J":  10,
	"10": 10,
	"9":  9,
	"8":  8,
	"7":  7,
	"6":  6,
	"5":  5,
	"4":  4,
	"3":  3,
	"2":  2,
}

// ランクからデフォルトスコアへの変換を行うヘルパー関数。
func rankToScore(rank Rank) int {
	return rankScoreMap[rank]
}

// CalculateScore は手札のスコアを計算して返します。
// J, Q, K は 10、A は 1 もしくは 11 として扱います。
func CalculateScore(cards []Card) int {
	score := 0
	aceCount := 0

	for _, c := range cards {
		score += rankToScore(c.Rank)
		if c.Rank == "A" {
			aceCount++
		}
	}

	// A の調整: 合計が 21 を超える場合は 1 として数えるように減算する
	for score > 21 && aceCount > 0 {
		score -= 10 // 11 -> 1 にするため 10 引く
		aceCount--
	}

	return score
}

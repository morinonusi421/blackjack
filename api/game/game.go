package game

import "errors"

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

// ValidateCore はゲーム状態の基本整合性を検証する
// - 手札の枚数（プレイヤー>=2, ディーラー>=1）
// - State/Result の矛盾がない（PlayerTurn↔Pending, Finished↔非Pending）
// - Bet と Payout の簡易整合性
func (g *Game) ValidateCore() error {
	if len(g.PlayerHand.Cards) < 2 {
		return errors.New("invalid state: player must have at least 2 cards")
	}
	if len(g.DealerHand.Cards) < 1 {
		return errors.New("invalid state: dealer must have at least 1 card")
	}
	if g.State == PlayerTurn && g.Result != Pending {
		return errors.New("invalid state: player turn but result is not pending")
	}
	if g.State == Finished && g.Result == Pending {
		return errors.New("invalid state: finished but result is pending")
	}
	if g.Bet <= 0 {
		return errors.New("invalid state: bet must be positive")
	}
	if g.Payout < 0 {
		return errors.New("invalid state: payout must be non-negative")
	}
	return nil
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
func RankToScore(rank Rank) int {
	switch rank {
	case "A":
		return 1 // エースは1として扱う（戦略計算・内部統一）
	case "K", "Q", "J":
		return 10
	case "10":
		return 10
	case "9":
		return 9
	case "8":
		return 8
	case "7":
		return 7
	case "6":
		return 6
	case "5":
		return 5
	case "4":
		return 4
	case "3":
		return 3
	case "2":
		return 2
	default:
		return 0
	}
}

// CalculateScore は手札のスコアを計算して返します。
// J, Q, K は 10、A は 1 もしくは 11 として扱います。
func CalculateScore(cards []Card) int {
	score := 0
	hasAce := false

	for _, c := range cards {
		s := RankToScore(c.Rank)
		score += s
		if c.Rank == "A" {
			hasAce = true
		}
	}

	// 21以上でバースト
	if score > 21 {
		return 0
	}

	// エースを持っていて、合計+10が21以下なら+10する
	if hasAce && score+10 <= 21 {
		score += 10
	}

	return score
}

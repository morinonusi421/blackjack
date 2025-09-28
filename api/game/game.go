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

// ValidateCore は終了状態を含めた基本整合性のみを検証する。
func (g *Game) ValidateCore() error {
	if len(g.PlayerHand.Cards) < 2 {
		return errors.New("invalid state: player must have at least 2 cards")
	}
	if len(g.DealerHand.Cards) < 1 {
		return errors.New("invalid state: dealer must have at least 1 card")
	}

	// スコアとカード配列の整合性
	computedPlayer := CalculateScore(g.PlayerHand.Cards)
	if g.PlayerHand.Score != computedPlayer {
		return errors.New("invalid state: player score does not match cards")
	}
	computedDealer := CalculateScore(g.DealerHand.Cards)
	if g.DealerHand.Score != computedDealer {
		return errors.New("invalid state: dealer score does not match cards")
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

// ValidateInProgress は進行中（PlayerTurn/Pending）の前提を検証する。
func (g *Game) ValidateInProgress() error {
	if g.State != PlayerTurn {
		return errors.New("invalid state: expected player turn in progress")
	}
	if g.Result != Pending {
		return errors.New("invalid state: expected pending result in progress")
	}
	if len(g.DealerHand.Cards) != 1 {
		return errors.New("invalid state: dealer must have exactly 1 card during player turn")
	}
	if g.PlayerHand.Score == 0 {
		return errors.New("invalid state: player already busted during player turn")
	}
	return nil
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

// RecalculateScores は手札スコアを Cards から再計算して反映する。
func (g *Game) RecalculateScores() {
	g.PlayerHand.Score = CalculateScore(g.PlayerHand.Cards)
	g.DealerHand.Score = CalculateScore(g.DealerHand.Cards)
}

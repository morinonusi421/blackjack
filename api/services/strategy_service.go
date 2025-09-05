package services

import (
	"blackjack/api/game"
	"blackjack/api/strategy"
)

// StrategyAdvisor は現在のゲーム状態から各アクションの期待払い戻しを返すインタフェース
type StrategyAdvisor interface {
	// Advise はゲーム状態と設定を入力に、期待払い戻しを返す
	Advise(g game.Game, config *game.GameConfig) (strategy.StrategyExpectedPayouts, error)
}

type strategyService struct {
	calc *strategy.Calculator
}

// NewStrategyService は最適戦略計算用のサービスを生成します。
func NewStrategyService() StrategyAdvisor {
	return &strategyService{calc: strategy.NewCalculator()}
}

// Advise は game.Game から strategy.StrategyState に変換し、期待払い戻しを計算して返します。
func (s *strategyService) Advise(g game.Game, config *game.GameConfig) (strategy.StrategyExpectedPayouts, error) {
	// 基本整合性
	if err := (&g).ValidateCore(); err != nil {
		return strategy.StrategyExpectedPayouts{}, err
	}

	// プレイヤー手札
	playerSum := 0
	playerHasAce := false
	for _, c := range g.PlayerHand.Cards {
		playerSum += game.RankToScore(c.Rank)
		if c.Rank == "A" {
			playerHasAce = true
		}
	}

	// ディーラー手札（アップカードのみ使用）
	dealerUpcard := g.DealerHand.Cards[0]
	dealerSum := game.RankToScore(dealerUpcard.Rank)
	dealerHasAce := dealerUpcard.Rank == "A"

	hasHit := len(g.PlayerHand.Cards) > 2

	st := strategy.StrategyState{
		Player: strategy.StrategyHand{Sum: playerSum, HasAce: playerHasAce},
		Dealer: strategy.StrategyHand{Sum: dealerSum, HasAce: dealerHasAce},
		HasHit: hasHit,
	}

	payouts := s.calc.CalculateAllExpectedPayouts(st, config)

	// 実際の払戻額を返すために、サービス層でスケーリング
	betF := float64(g.Bet)
	payouts.HitPayout *= betF
	payouts.StandPayout *= betF
	payouts.SurrenderPayout *= betF
	payouts.BestPayout *= betF
	return payouts, nil
}

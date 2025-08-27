package services

import (
	"testing"

	"blackjack/api/game"
)

func TestStrategyService_Advise_ConvertsGameToStrategyState(t *testing.T) {
	svc := NewStrategyService()

	g := game.Game{
		PlayerHand: game.Hand{Cards: []game.Card{{Suit: game.Spade, Rank: "A"}, {Suit: game.Heart, Rank: "9"}}},
		DealerHand: game.Hand{Cards: []game.Card{{Suit: game.Diamond, Rank: "7"}}},
		State:      game.PlayerTurn,
		Result:     game.Pending,
		Bet:        100,
	}

	payouts, err := svc.Advise(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// スケーリング後は bet 倍（ここでは100倍）の金額で返ることをゆるく検証
	if payouts.HitPayout < 0.0 || payouts.HitPayout > 200.0 {
		t.Fatalf("HitPayout out of range: %v", payouts.HitPayout)
	}
	if payouts.StandPayout < 0.0 || payouts.StandPayout > 200.0 {
		t.Fatalf("StandPayout out of range: %v", payouts.StandPayout)
	}
	if payouts.SurrenderPayout != 50.0 {
		t.Fatalf("expected SurrenderPayout 50.0 when HasHit=false, got %v", payouts.SurrenderPayout)
	}
}

func TestStrategyService_Advise_InvalidInput(t *testing.T) {
	svc := NewStrategyService()

	// ディーラーのカードが無い
	g1 := game.Game{PlayerHand: game.Hand{Cards: []game.Card{{Suit: game.Spade, Rank: "A"}, {Suit: game.Heart, Rank: "9"}}}, DealerHand: game.Hand{Cards: []game.Card{}}}
	if _, err := svc.Advise(g1); err == nil {
		t.Fatalf("expected error for missing dealer upcard")
	}

	// プレイヤーのカードが1枚
	g2 := game.Game{PlayerHand: game.Hand{Cards: []game.Card{{Suit: game.Spade, Rank: "A"}}}, DealerHand: game.Hand{Cards: []game.Card{{Suit: game.Diamond, Rank: "7"}}}}
	if _, err := svc.Advise(g2); err == nil {
		t.Fatalf("expected error for insufficient player cards")
	}
}

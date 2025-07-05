package services

import (
	"errors"

	"blackjack/api/game"
)

// GameService はブラックジャックゲームのユースケースを提供します。
// 例えばゲーム開始、ヒット、スタンドなど。
// 現時点ではゲーム開始のみ実装します。

type GameService interface {
	NewGame(bet int) (game.Game, error)
	// Stand はプレイヤーがスタンドした後のディーラー処理を行い、current を更新して結果を反映します。
	Stand(g *game.Game) error
}

type gameService struct {
	deck game.Deck
}

// NewGameServiceは Deck を受け取りゲームサービスを生成します。
func NewGameService(deck game.Deck) GameService {
	if deck == nil {
		panic("deck must not be nil")
	}
	return &gameService{deck: deck}
}

// NewGame は掛け金を受け取り、新しいゲームを初期化して返します。
// bet が 1 未満の場合はエラーを返します。
func (s *gameService) NewGame(bet int) (game.Game, error) {
	if bet <= 0 {
		return game.Game{}, errors.New("bet must be positive")
	}

	playerCards := []game.Card{s.deck.Deal(), s.deck.Deal()}
	dealerCards := []game.Card{s.deck.Deal()}

	playerScore := game.CalculateScore(playerCards)
	dealerScore := game.CalculateScore(dealerCards)

	playerHand := game.Hand{Cards: playerCards, Score: playerScore}
	dealerHand := game.Hand{Cards: dealerCards, Score: dealerScore}

	g := game.Game{
		PlayerHand:    playerHand,
		DealerHand:    dealerHand,
		State:         game.PlayerTurn,
		Bet:           bet,
		Result:        game.Pending,
		ResultMessage: "",
		Payout:        0,
	}

	// ブラックジャック判定
	if playerScore == 21 {
		g.State = game.Finished
		g.Result = game.PlayerWin
		g.ResultMessage = "Blackjack! Player wins."
		g.Payout = bet * 3 / 2 // 1.5 倍
	}

	return g, nil
}

// Stand はプレイヤーターン終了後、ディーラーが17以上になるまでカードを引き、
// 最終結果を判定して Game を返します。
// g.State が PlayerTurn でない場合はエラーを返します。
func (s *gameService) Stand(g *game.Game) error {
	// 引数チェック
	if g == nil {
		return errors.New("game must not be nil")
	}

	// プレイヤーは 2 枚、ディーラーは 1 枚以上の手札が必要
	if len(g.PlayerHand.Cards) != 2 {
		return errors.New("invalid state: player must have exactly 2 cards")
	}
	if len(g.DealerHand.Cards) < 1 {
		return errors.New("invalid state: dealer must have at least 1 card")
	}

	// ゲーム結果がまだ確定していないことを確認
	if g.Result != game.Pending {
		return errors.New("invalid state: game already finished")
	}
	if g.State != game.PlayerTurn {
		return errors.New("invalid state: game is not in player turn")
	}

	// ディーラーは17以上になるまでヒット
	for g.DealerHand.Score < 17 {
		card := s.deck.Deal()
		g.DealerHand.Cards = append(g.DealerHand.Cards, card)
		g.DealerHand.Score = game.CalculateScore(g.DealerHand.Cards)
	}

	dealerScore := g.DealerHand.Score

	// 結果判定
	var (
		result game.Result
		payout int
		msg    string
	)

	playerScore := g.PlayerHand.Score

	switch {
	case dealerScore > 21:
		result = game.PlayerWin
		payout = g.Bet * 2
		msg = "Dealer busts! Player wins."
	case dealerScore < playerScore:
		result = game.PlayerWin
		payout = g.Bet * 2
		msg = "Player wins."
	case dealerScore > playerScore:
		result = game.DealerWin
		payout = 0
		msg = "Dealer wins."
	default: // 引き分け
		result = game.Push
		payout = g.Bet // 掛け金を返却
		msg = "Push. Bet returned."
	}

	g.State = game.Finished
	g.Result = result
	g.ResultMessage = msg
	g.Payout = payout

	return nil
}

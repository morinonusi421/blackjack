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
}

type gameService struct{}

// NewGameService はゲームサービスのデフォルト実装を返します。
func NewGameService() GameService {
	return &gameService{}
}

// NewGame は掛け金を受け取り、新しいゲームを初期化して返します。
// bet が 1 未満の場合はエラーを返します。
func (s *gameService) NewGame(bet int) (game.Game, error) {
	if bet <= 0 {
		return game.Game{}, errors.New("bet must be positive")
	}

	playerCards := []game.Card{game.NewRandomCard(), game.NewRandomCard()}
	dealerCards := []game.Card{game.NewRandomCard()}

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
		BalanceChange: 0,
	}

	// ブラックジャック判定
	if playerScore == 21 {
		g.State = game.Finished
		g.Result = game.PlayerWin
		g.ResultMessage = "Blackjack! Player wins."
		g.BalanceChange = bet * 3 / 2 // 1.5 倍
	}

	return g, nil
}

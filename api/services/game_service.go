package services

import (
	"errors"

	"blackjack/api/game"
)

// GameStarter は新規ゲーム開始のみを表す最小インタフェース
type GameStarter interface {
	NewGame(bet int) (game.Game, error)
}

// Hitter はヒット（カードを引く）処理のみを表す最小インタフェース
type Hitter interface {
	Hit(*game.Game) error
}

// Stander はスタンド処理のみを表す最小インタフェース
type Stander interface {
	Stand(*game.Game) error
}

// Surrenderer はサレンダー処理のみを表す最小インタフェース
type Surrenderer interface {
	Surrender(*game.Game) error
}

// GameService はブラックジャックに必要な全ての処理を提供するインターフェース
type GameService interface {
	GameStarter
	Hitter
	Stander
	Surrenderer
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
	if len(g.PlayerHand.Cards) < 2 {
		return errors.New("invalid state: player must have at least 2 cards")
	}
	if len(g.DealerHand.Cards) != 1 {
		return errors.New("invalid state: dealer must have exactly 1 card")
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

func (s *gameService) Hit(g *game.Game) error {
	// 引数チェック
	if g == nil {
		return errors.New("game must not be nil")
	}

	// ゲーム状態がプレイヤーターンであることを確認
	if g.State != game.PlayerTurn {
		return errors.New("invalid state: game is not in player turn")
	}

	// 既に結果が確定していないか確認
	if g.Result != game.Pending {
		return errors.New("invalid state: game already finished")
	}

	// 1 枚カードを配る
	card := s.deck.Deal()
	g.PlayerHand.Cards = append(g.PlayerHand.Cards, card)
	g.PlayerHand.Score = game.CalculateScore(g.PlayerHand.Cards)

	playerScore := g.PlayerHand.Score

	// バーストチェック
	if playerScore > 21 {
		g.State = game.Finished
		g.Result = game.DealerWin
		g.ResultMessage = "Player busts! Dealer wins."
		g.Payout = 0
		return nil
	}

	// 21 ちょうどの場合は自動的にスタンド相当の処理を行う
	if playerScore == 21 {
		return s.Stand(g)
	}

	// それ以外（21 未満）の場合は引き続きプレイヤーターン
	return nil
}

// Surrender はプレイヤーがサレンダー（降参）を選択した時の処理を行います。
// 掛け金の半分を失い、ゲームを終了します。
// プレイヤーは最初の2枚のカードを受け取った後にのみサレンダーできます。
func (s *gameService) Surrender(g *game.Game) error {
	// 引数チェック
	if g == nil {
		return errors.New("game must not be nil")
	}

	// ゲーム状態がプレイヤーターンであることを確認
	if g.State != game.PlayerTurn {
		return errors.New("invalid state: game is not in player turn")
	}

	// 既に結果が確定していないか確認
	if g.Result != game.Pending {
		return errors.New("invalid state: game already finished")
	}

	// プレイヤーが最初の2枚のカードしか持っていないことを確認
	if len(g.PlayerHand.Cards) != 2 {
		return errors.New("invalid state: surrender is only allowed with initial 2 cards")
	}

	// サレンダー処理
	g.State = game.Finished
	g.Result = game.Surrender
	g.ResultMessage = "Player surrendered."
	g.Payout = g.Bet / 2 // 掛け金の半分を返却

	return nil
}

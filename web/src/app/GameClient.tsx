'use client';

import { useState } from 'react';
import useGame from '../hooks/useGame';
import GameInfo from '../components/GameInfo';
import ErrorMessage from '../components/ErrorMessage';
import StartGameForm from '../components/StartGameForm';
import Balance from '../components/Balance';
import ActionButtons from '../components/ActionButtons';
import StrategyAdvice from '../components/StrategyAdvice';

export default function GameClient() {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL;
  const { game, loading, error, startGame, stand, hit, surrender, balance, advice, getAdvice } = useGame(apiUrl);
  const [bet, setBet] = useState(100);

  const handleStart = () => {
    startGame(bet);
  };

  // サレンダーは最初の2枚のカードを受け取った後にのみ可能
  const canSurrender = game?.player_hand.cards.length === 2;

  const gameInProgress = game?.state === 'PlayerTurn';

  return (
    <section
      style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        gap: '24px',
      }}
    >
      {/* 所持金表示 */}
      <Balance balance={balance} />

      {/* 掛け金入力 & ゲーム開始ボタン */}
      <StartGameForm
        bet={bet}
        onBetChange={setBet}
        loading={loading}
        onStart={handleStart}
        disabled={gameInProgress}
      />

      {/* エラーメッセージ */}
      <ErrorMessage message={error ?? ''} />

      {/* ゲーム状況表示 */}
      <GameInfo game={game} />

      {/* 戦略の期待払い戻し表示 */}
      {game && game.state === 'PlayerTurn' && advice && (
        <StrategyAdvice advice={advice} canSurrender={canSurrender} bet={game.bet} />
      )}

      {/* Hit / Stand / Surrender アクション */}
      {game && game.state === 'PlayerTurn' && (
        <ActionButtons 
          onHit={hit} 
          onStand={stand} 
          onSurrender={surrender} 
          canSurrender={canSurrender}
          disabled={loading} 
          onShowAdvice={getAdvice}
        />
      )}
    </section>
  );
} 

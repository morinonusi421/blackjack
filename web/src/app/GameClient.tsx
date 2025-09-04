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
  const canSurrender = !!(game && game.player_hand.cards.length === 2);

  const gameInProgress = game?.state === 'PlayerTurn';
  const controlsDisabled = !gameInProgress || loading;

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

      {/* Hit / Stand / Surrender アクション（常時表示・未開始時は非活性） */}
      <ActionButtons 
        onHit={hit} 
        onStand={stand} 
        onSurrender={surrender} 
        canSurrender={canSurrender}
        disabled={controlsDisabled} 
        onShowAdvice={getAdvice}
      />

      {/* 二列レイアウト: 左=ゲーム状況, 右=助言（固定幅プレースホルダを常設） */}
      <div className="gc-two-col">
        {/* 左カラム: ゲーム状況 */}
        <div className="gc-left">
          <GameInfo game={game} />
        </div>

        {/* 右カラム: 助言パネル */}
        <div className="gc-right">
          {game && game.state === 'PlayerTurn' && advice ? (
            <StrategyAdvice advice={advice} canSurrender={canSurrender} bet={game.bet} />
          ) : (
            <section
              aria-label="strategy-advice"
              style={{
                width: '100%',
                border: '1px solid #e5e7eb',
                borderRadius: '8px',
                padding: '16px',
                background: '#fafafa',
                color: '#6b7280',
                minHeight: '120px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                textAlign: 'center',
              }}
            >
              答えを見る を押すと、ここに期待払い戻しが表示されます。
            </section>
          )}
        </div>
      </div>
    </section>
  );
} 

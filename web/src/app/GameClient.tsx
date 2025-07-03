'use client';

import { useState } from 'react';
import useGame from '../hooks/useGame';
import GameInfo from '../components/GameInfo';
import ErrorMessage from '../components/ErrorMessage';
import StartGameForm from '../components/StartGameForm';

export default function GameClient() {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL;
  const { game, loading, error, startGame } = useGame(apiUrl);
  const [bet, setBet] = useState(100);

  const handleStart = () => {
    startGame(bet);
  };

  return (
    <section
      style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        gap: '24px',
      }}
    >
      {/* 掛け金入力 & ゲーム開始ボタン */}
      <StartGameForm
        bet={bet}
        onBetChange={setBet}
        loading={loading}
        onStart={handleStart}
      />

      {/* エラーメッセージ */}
      <ErrorMessage message={error ?? ''} />

      {/* ゲーム状況表示 */}
      <GameInfo game={game} />
    </section>
  );
} 

'use client';

import { useState } from 'react';
import GameClient from '../app/GameClient';
import GameRules from './GameRules';
import { DEFAULT_DEALER_THRESHOLD } from '../constants/config';

export default function HomeContent() {
  const [currentThreshold, setCurrentThreshold] = useState(DEFAULT_DEALER_THRESHOLD);

  return (
    <>
      <GameClient onThresholdChange={setCurrentThreshold} />
      <p style={{ color: '#6b7280', fontSize: '13px' }}>
        ボタンがうまく反応しないときは、バックエンドサーバーがスリープ中の可能性があります（無料サーバーのため）。数分お待ちいただくと復帰します。
      </p>
      <GameRules dealerThreshold={currentThreshold} />
    </>
  );
}
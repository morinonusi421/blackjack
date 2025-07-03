'use client';

import { useState, useCallback } from 'react';
import type { Game } from '../types/game';

/**
 * ブラックジャックの新規ゲーム開始を扱うカスタムフック。
 *
 * @param apiUrl API のベース URL（例: https://example.com）
 */
export default function useGame(apiUrl?: string | null) {
  const [game, setGame] = useState<Game | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  /**
   * ゲームを開始します。
   *
   * @param bet 掛け金
   */
  const startGame = useCallback(
    async (bet: number) => {
      if (!apiUrl) {
        setError('APIのURLが設定されていません。');
        return;
      }
      if (bet <= 0) {
        setError('掛け金は 1 以上で入力してください。');
        return;
      }

      setLoading(true);
      setError(null);

      try {
        const res = await fetch(`${apiUrl}/api/game/new`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ bet }),
        });

        if (!res.ok) throw new Error('APIからの応答がありませんでした');
        const result: Game = await res.json();
        setGame(result);
      } catch (err: unknown) {
        if (err instanceof Error) {
          setError(err.message);
        } else {
          setError('不明なエラーが発生しました');
        }
      } finally {
        setLoading(false);
      }
    },
    [apiUrl],
  );

  return {
    game,
    loading,
    error,
    startGame,
  } as const;
} 

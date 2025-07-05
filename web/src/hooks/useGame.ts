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
  const [balance, setBalance] = useState(10000);

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
      if (balance <= 0) {
        setError('所持金がありません。');
        return;
      }
      if (bet > balance) {
        setError('所持金が不足しています。');
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
        setBalance((prev) => prev - bet + result.payout);
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
    [apiUrl, balance],
  );

  /**
   * スタンド（カード引き止め）を行う。
   * 現在の game を API に送り、結果を受け取って状態と所持金を更新する。
   */
  const stand = useCallback(async () => {
    if (!apiUrl) {
      setError('APIのURLが設定されていません。');
      return;
    }
    if (!game) {
      setError('ゲームが開始されていません。');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const res = await fetch(`${apiUrl}/api/game/stand`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(game),
      });

      if (!res.ok) throw new Error('APIからの応答がありませんでした');
      const result: Game = await res.json();
      setGame(result);
      // 返ってきたチップを所持金に加算
      setBalance((prev) => prev + result.payout);
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('不明なエラーが発生しました');
      }
    } finally {
      setLoading(false);
    }
  }, [apiUrl, game]);

  return {
    game,
    loading,
    error,
    startGame,
    stand,
    balance,
  } as const;
} 

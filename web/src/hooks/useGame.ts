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
   * ゲーム状態を取得して state / balance / error / loading を一括で更新する共通関数
   *
   * @param endpoint API のパス（例: '/api/game/new'）
   * @param payload  リクエストボディ
   * @param betAmount 掛け金（新規ゲーム開始時のみ指定）
   */
  const fetchAndUpdateGame = useCallback(
    async (endpoint: string, payload: unknown, betAmount?: number) => {
      if (!apiUrl) {
        setError('APIのURLが設定されていません。');
        return;
      }

      setLoading(true);
      setError(null);

      try {
        const res = await fetch(`${apiUrl}${endpoint}`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload),
        });

        if (!res.ok) throw new Error('APIからの応答がありませんでした');
        const result: Game = await res.json();
        setGame(result);

        // `betAmount` が指定されている場合は掛け金を差し引いた上で払い戻しを加算、それ以外は払い戻しのみ加算
        setBalance((prev) =>
          betAmount !== undefined ? prev - betAmount + result.payout : prev + result.payout
        );
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

  /**
   * ゲームを開始します。
   *
   * @param bet 掛け金
   */
  const startGame = useCallback(
    async (bet: number) => {
      // すでにゲームが進行中の場合は新しいゲームを開始しない
      if (game?.state === 'PlayerTurn') {
        setError('ゲームが進行中です。');
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

      await fetchAndUpdateGame('/api/game/new', { bet }, bet);
    },
    [balance, game, fetchAndUpdateGame],
  );

  /**
   * スタンド（カード引き止め）を行う。
   */
  const stand = useCallback(async () => {
    if (!game) {
      setError('ゲームが開始されていません。');
      return;
    }

    await fetchAndUpdateGame('/api/game/stand', game);
  }, [game, fetchAndUpdateGame]);

  /**
   * ヒット（カードを1枚引く）を行う。
   */
  const hit = useCallback(async () => {
    if (!game) {
      setError('ゲームが開始されていません。');
      return;
    }

    await fetchAndUpdateGame('/api/game/hit', game);
  }, [game, fetchAndUpdateGame]);

  /**
   * サレンダー（降参）を行う。
   * 掛け金の半分を失い、ゲームを終了する。
   */
  const surrender = useCallback(async () => {
    if (!game) {
      setError('ゲームが開始されていません。');
      return;
    }

    await fetchAndUpdateGame('/api/game/surrender', game);
  }, [game, fetchAndUpdateGame]);

  return {
    game,
    loading,
    error,
    startGame,
    stand,
    hit,
    surrender,
    balance,
  } as const;
} 

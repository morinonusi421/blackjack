'use client';

import { useState, useCallback, useRef } from 'react';
import type { Game, StrategyAdvice } from '../types/game';
import { DEFAULT_DEALER_THRESHOLD } from '../constants/config';

/**
 * ブラックジャックの新規ゲーム開始を扱うカスタムフック。
 *
 * @param apiUrl API のベース URL（例: https://example.com）
 * @param dealerThreshold ディーラーのスタンド閾値
 */
export default function useGame(apiUrl?: string | null, dealerThreshold?: number) {
  const [game, setGame] = useState<Game | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [balance, setBalance] = useState(1000);
  const [advice, setAdvice] = useState<StrategyAdvice | null>(null);

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
        // ゲーム進行のAPIには設定を含める
        let requestBody;
        if (endpoint === '/api/game/new') {
          requestBody = payload as object; // 新規ゲーム開始時はconfigは不要
        } else {
          requestBody = {
            game: payload,
            config: { dealer_stand_threshold: dealerThreshold || DEFAULT_DEALER_THRESHOLD }
          };
        }

        const res = await fetch(`${apiUrl}${endpoint}`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(requestBody),
        });

        if (!res.ok) throw new Error('APIからの応答がありませんでした');
        const result: Game = await res.json();
        setGame(result);
        // アクションを実際に行ったので、過去のアドバイスはクリア
        setAdvice(null);

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
    [apiUrl, dealerThreshold],
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

  /**
   * 現在のゲーム状態に対する戦略アドバイス（期待払い戻し）を取得
   */
  const getAdvice = useCallback(async () => {
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
      // ローカル設定を使用
      const config = {
        dealer_stand_threshold: dealerThreshold || DEFAULT_DEALER_THRESHOLD,
      };

      // ゲーム状態と設定を一緒に送信
      const res = await fetch(`${apiUrl}/api/strategy/advise`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          game: game,
          config: config,
        }),
      });
      if (!res.ok) throw new Error('APIからの応答がありませんでした');
      const data: StrategyAdvice = await res.json();
      setAdvice(data);
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('不明なエラーが発生しました');
      }
    } finally {
      setLoading(false);
    }
  }, [apiUrl, game, dealerThreshold]);

  const clearAdvice = useCallback(() => {
    setAdvice(null);
  }, []);

  const refreshAdviceIfVisible = useCallback(async () => {
    if (advice && game) {
      // 既にアドバイスが表示されている場合は、新しい設定で再計算
      await getAdvice();
    }
  }, [advice, game, getAdvice]);

  // デバウンス用のタイマーref
  const debounceTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const debouncedRefreshAdviceIfVisible = useCallback((newThreshold: number) => {
    // 既存のタイマーをクリア
    if (debounceTimeoutRef.current) {
      clearTimeout(debounceTimeoutRef.current);
    }
    
    // 新しいタイマーを設定（100ms後に実行）
    debounceTimeoutRef.current = setTimeout(() => {
      if (advice && game) {
        // 最新の閾値を使用してアドバイスを再計算
        const config = {
          dealer_stand_threshold: newThreshold || DEFAULT_DEALER_THRESHOLD,
        };

        setLoading(true);
        setError(null);
        
        fetch(`${apiUrl}/api/strategy/advise`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            game: game,
            config: config,
          }),
        })
        .then(res => {
          if (!res.ok) throw new Error('APIからの応答がありませんでした');
          return res.json();
        })
        .then((data: StrategyAdvice) => {
          setAdvice(data);
        })
        .catch((err: unknown) => {
          if (err instanceof Error) {
            setError(err.message);
          } else {
            setError('不明なエラーが発生しました');
          }
        })
        .finally(() => {
          setLoading(false);
        });
      }
    }, 100);
  }, [advice, game, apiUrl]);

  return {
    game,
    loading,
    error,
    startGame,
    stand,
    hit,
    surrender,
    balance,
    advice,
    getAdvice,
    clearAdvice,
    refreshAdviceIfVisible,
    debouncedRefreshAdviceIfVisible,
  } as const;
} 

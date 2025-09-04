'use client';

import type { StrategyAdvice as StrategyAdviceType } from '../types/game';

interface StrategyAdviceProps {
  advice: StrategyAdviceType;
  canSurrender: boolean;
  bet: number;
}

export default function StrategyAdvice({ advice, canSurrender, bet }: StrategyAdviceProps) {
  const items = [
    { label: 'Hit', value: advice.hit_payout },
    { label: 'Stand', value: advice.stand_payout },
    ...(canSurrender ? [{ label: 'Surrender', value: advice.surrender_payout }] : []),
  ];

  const best = Math.max(...items.map((i) => i.value));

  return (
    <section
      aria-label="strategy-advice"
      style={{
        width: '100%',
        border: '1px solid #ccc',
        borderRadius: '8px',
        padding: '16px',
        background: '#fafafa',
        minHeight: '120px',
      }}
    >
      <h3 style={{ marginTop: 0, marginBottom: '8px' }}>期待払い戻し（ベット{bet}）</h3>
      <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
        {items.map((i) => {
          const isBest = i.value === best;
          return (
            <div
              key={i.label}
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                padding: '10px 12px',
                borderRadius: '6px',
                border: '1px solid #e5e5e5',
                backgroundColor: isBest ? '#e6ffed' : '#fff',
                fontWeight: isBest ? 700 : 400,
              }}
            >
              <span>{i.label}</span>
              <span>{i.value.toFixed(2)}</span>
            </div>
          );
        })}
      </div>
      <p style={{ marginTop: '10px', color: '#555' }}>緑が現在の最適行動です。</p>
    </section>
  );
}

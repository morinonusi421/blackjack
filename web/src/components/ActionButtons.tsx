'use client';

// no-op
import ShowAdviceButton from './ShowAdviceButton';
import HitButton from './HitButton';
import StandButton from './StandButton';
import SurrenderButton from './SurrenderButton';

interface ActionButtonsProps {
  onStand: () => void;
  onHit: () => void;
  onSurrender: () => void;
  canSurrender: boolean;
  disabled?: boolean;
  onShowAdvice?: () => void;
}

export default function ActionButtons({ onStand, onHit, onSurrender, canSurrender, disabled = false, onShowAdvice }: ActionButtonsProps) {
  // サレンダーは初回アクションでのみ可能
  const surrenderDisabled = disabled || !canSurrender;

  // 各ボタンのスタイルは個別コンポーネント側で定義

  // SurrenderButton にスタイル委譲のため個別スタイルは不要

  return (
    <div className="action-buttons">
      <ShowAdviceButton onClick={onShowAdvice} disabled={disabled} />
      <HitButton onClick={onHit} disabled={disabled} />
      <StandButton onClick={onStand} disabled={disabled} />
      <SurrenderButton onClick={onSurrender} disabled={surrenderDisabled} />
    </div>
  );
}

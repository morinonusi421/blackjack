'use client';

import type { CSSProperties } from 'react';

interface SurrenderButtonProps {
  onClick?: () => void;
  disabled?: boolean;
}

export default function SurrenderButton({ onClick, disabled = false }: SurrenderButtonProps) {
  const baseStyle: CSSProperties = {
    padding: '10px 20px',
    backgroundColor: disabled ? '#999' : '#ff4444',
    color: '#fff',
    border: 'none',
    borderRadius: '4px',
    fontSize: '16px',
    cursor: disabled ? 'not-allowed' : 'pointer',
    opacity: disabled ? 0.6 : 1,
    transition: 'filter 0.2s',
  };

  return (
    <button
      onClick={onClick}
      disabled={disabled}
      style={baseStyle}
      onMouseOver={(e) => {
        if (!disabled) (e.currentTarget.style.filter = 'brightness(1.1)');
      }}
      onMouseOut={(e) => {
        if (!disabled) (e.currentTarget.style.filter = 'brightness(1)');
      }}
    >
      Surrender
    </button>
  );
}

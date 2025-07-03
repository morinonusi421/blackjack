'use client';

interface BalanceProps {
  balance: number;
}

export default function Balance({ balance }: BalanceProps) {
  return (
    <div
      style={{
        fontSize: '24px',
        fontWeight: 'bold',
        textAlign: 'center',
      }}
    >
      所持金: {balance}
    </div>
  );
} 

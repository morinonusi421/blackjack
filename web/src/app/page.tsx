import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'ブラックジャック'
};

export default function Home() {
  return (
    <main
      style={{
        fontFamily: 'sans-serif',
        textAlign: 'center',
        marginTop: '40px',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        gap: '24px',
      }}
    >
      <p>
        &quot;エラー: Failed to fetch&quot;が発生した時は、バックエンドサーバーがまだ眠っています。（無料サーバーなのですぐスリープしてしまいます）壊れてなければ、数分待つと動きます。
      </p>
    </main>
  );
}

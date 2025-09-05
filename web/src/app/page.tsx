import type { Metadata } from 'next';
import HomeContent from '../components/HomeContent';

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
      <HomeContent />
    </main>
  );
}

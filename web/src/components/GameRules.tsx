'use client';

interface GameRulesProps {
  dealerThreshold: number;
}

export default function GameRules({ dealerThreshold }: GameRulesProps) {
  return (
    <section
      style={{
        maxWidth: '840px',
        textAlign: 'left',
        lineHeight: 1.9,
        fontSize: '15px',
        background: '#f9fafb',
        padding: '16px 20px',
        borderRadius: '8px',
        boxShadow: '0 1px 2px rgba(0,0,0,0.04)'
      }}
    >
      <h2 style={{ fontSize: '20px', margin: '16px 0 8px' }}>ゲームの目的と流れ</h2>
      <p>
        ブラックジャックは、短いラウンドで運と判断を楽しむカードゲームです。目的は自分の手札合計を21に近づけ、ディーラーより高くすることです。
      </p>

      <ol style={{ paddingLeft: '1.2em', listStyleType: 'decimal' }}>
        <li style={{ margin: '6px 0' }}>配札: プレイヤーは2枚、ディーラーは1枚の公開カードを受け取ります。</li>
        <li style={{ margin: '6px 0' }}>プレイヤーの手番: 21を超えないように、ギリギリまでカードを引いていきます。21を超えるとその時点で敗北。</li>
        <li style={{ margin: '6px 0' }}>ディーラーの手番: 合計が<span style={{ color: 'red', fontWeight: 'bold' }}>{dealerThreshold}</span>以上になるまで自動でカードを引き、プレイヤーと勝負をします。</li>
        <li style={{ margin: '6px 0' }}>勝敗判定: バーストしていなければ、21により近い方が勝ち。同じなら引き分け。</li>
      </ol>
      <p>
        補足：J、Q、Kは10として数え、11,12,13とは数えません。
      </p>
      <p>
        補足：Aは1または11として、自分にとって都合の良いほうで数えます。
      </p>
      <hr style={{ margin: '20px 0', border: 0, borderTop: '2px solid #cbd5e1' }} />

      <h3 style={{ fontSize: '18px', margin: '16px 0 8px' }}>プレイヤーの手番で行えるアクション</h3>
      <ul style={{ paddingLeft: '1.2em', listStyleType: 'disc' }}>
        <li style={{ margin: '6px 0' }}>
          <strong>ヒット（Hit）</strong>: もう1枚カードを引きます。合計が21を超えると即バーストで負け。手札が弱いとき、まだ余裕があると判断したときに選びます。
        </li>
        <li style={{ margin: '6px 0' }}>
          <strong>スタンド（Stand）</strong>: カードを引くのをやめ、手番を終了します。現在の合計でディーラーと勝負します。21に近く、これ以上引くと危険だと判断したときに選びます。
        </li>
        <li style={{ margin: '6px 0' }}>
          <strong>サレンダー（Surrender）</strong>: ラウンドを降ります。掛け金の半分だけは払い戻されます。勝ち目が薄い初期手札のときに損失を抑える選択です。
        </li>
      </ul>
      <hr style={{ margin: '20px 0', border: 0, borderTop: '2px solid #cbd5e1' }} />

      <h3 style={{ fontSize: '18px', margin: '16px 0 8px' }}>ディーラーの動作</h3>
      <ul style={{ paddingLeft: '1.2em', listStyleType: 'disc' }}>
        <li>ディーラーは合計が<span style={{ color: 'red', fontWeight: 'bold' }}>{dealerThreshold}</span>以上になるまで自動でヒットし、以降はスタンドします。</li>
      </ul>
      <hr style={{ margin: '20px 0', border: 0, borderTop: '2px solid #cbd5e1' }} />

      <h3 style={{ fontSize: '18px', margin: '16px 0 8px' }}>勝敗の決まり</h3>
      <ul style={{ paddingLeft: '1.2em', listStyleType: 'disc' }}>
        <li style={{ margin: '6px 0' }}>プレイヤーが21を超えるとバーストで即負け。（ディラーにターンは回らない）</li>
        <li style={{ margin: '6px 0' }}>ディラーがバーストしたらディラーの負け</li>
        <li style={{ margin: '6px 0' }}>両者がバーストしなければ、21により近いほうが勝ち。同点はプッシュ（引き分け）</li>
      </ul>
      <hr style={{ margin: '20px 0', border: 0, borderTop: '2px solid #cbd5e1' }} />

      <h3 style={{ fontSize: '18px', margin: '16px 0 8px' }}>配当と残高の動き</h3>
      <p>
        このゲームでは、ベットは開始時に所持金から差し引かれ、ラウンド終了時に結果に応じて以下の金額が払い戻されます。
      </p>
      <ul style={{ paddingLeft: '1.2em', listStyleType: 'disc' }}>
        <li style={{ margin: '6px 0' }}><strong>勝ち</strong>: ベット×2 を受け取り（差し引かれたベットが戻り、同額の利益）</li>
        <li style={{ margin: '6px 0' }}><strong>ブラックジャック</strong>: ベット×2.5 を受け取り（最初の2枚が A と10点札。例: ベット100なら250が戻り＝利益150）</li>
        <li style={{ margin: '6px 0' }}><strong>プッシュ（引き分け）</strong>: ベット×1 を受け取り（収支±0）</li>
        <li style={{ margin: '6px 0' }}><strong>サレンダー</strong>: ベット×0.5 を受け取り</li>
        <li style={{ margin: '6px 0' }}><strong>負け・バースト</strong>: 受け取り 0</li>
      </ul>
      <hr style={{ margin: '20px 0', border: 0, borderTop: '2px solid #cbd5e1' }} />

      <h3 style={{ fontSize: '18px', margin: '16px 0 8px' }}>補足</h3>
      <ul style={{ paddingLeft: '1.2em', listStyleType: 'disc' }}>
        <li style={{ margin: '6px 0' }}>
          山札は無限（無限デッキ）として扱います。各ドローは独立で、すべてのカードが常に同じ確率で引けます。
        </li>
        <li style={{ margin: '6px 0' }}>
          「答えを見る」ボタンは、現在の手札・ディーラーの公開カード・ベット額をもとに、Hit / Stand /（初手のみ）Surrender の期待払い戻しを計算して表示します。
        </li>
      </ul>
    </section>
  );
}
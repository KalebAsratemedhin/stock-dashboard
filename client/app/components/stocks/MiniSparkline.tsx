'use client';

import { StockQuote } from '@/types';

interface MiniSparklineProps {
  quotes: StockQuote[];
}

export function MiniSparkline({ quotes }: MiniSparklineProps) {
  const points = quotes.slice(-10).map((q) => q.close);
  if (points.length < 2) {
    return <span className="inline-block w-[60px] h-6 text-slate-600 text-xs">—</span>;
  }
  const min = Math.min(...points);
  const max = Math.max(...points);
  const range = max - min || 1;
  const w = 60;
  const h = 24;
  const path = points
    .map((p, i) => {
      const x = (i / (points.length - 1)) * w;
      const y = h - ((p - min) / range) * h;
      return `${i === 0 ? 'M' : 'L'} ${x} ${y}`;
    })
    .join(' ');
  const isUp = points[points.length - 1] >= points[0];
  return (
    <svg width={w} height={h} className="inline-block">
      <path
        d={path}
        fill="none"
        stroke={isUp ? 'var(--cyber-orange)' : '#94a3b8'}
        strokeWidth={1.5}
      />
    </svg>
  );
}

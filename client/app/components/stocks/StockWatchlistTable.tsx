'use client';

import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { BarChart3 } from 'lucide-react';
import { cn } from '@/lib/utils';
import { StockQuote } from '@/types';
import { SYMBOL_TO_COMPANY } from '@/lib/constants';
import { MiniSparkline } from './MiniSparkline';

interface StockWatchlistTableProps {
  quotes: StockQuote[];
  selectedSymbol: string;
  onSelectSymbol: (symbol: string) => void;
  getQuotesForSparkline: (symbol: string) => StockQuote[];
}

export function StockWatchlistTable({
  quotes,
  selectedSymbol,
  onSelectSymbol,
  getQuotesForSparkline,
}: StockWatchlistTableProps) {
  return (
    <Card className="bg-surface-slate border-border-slate">
      <CardContent className="p-6">
        <div className="flex flex-wrap items-center justify-between gap-4 mb-4">
          <h3 className="text-cyber-orange text-[10px] font-black uppercase tracking-[0.2em] flex items-center gap-2">
            <BarChart3 className="h-4 w-4" />
            Live Matrix Watchlist
          </h3>
          <Button
            variant="outline"
            size="sm"
            className="text-[10px] font-bold uppercase border-border-slate text-slate-400 hover:text-cyber-orange hover:border-cyber-orange"
          >
            + Inject New Asset
          </Button>
        </div>
        <div className="overflow-x-auto border border-border-slate rounded-sm">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border-slate">
                <th className="text-left py-3 px-4 text-slate-500 text-[10px] font-bold uppercase tracking-wider">
                  Node
                </th>
                <th className="text-left py-3 px-4 text-slate-500 text-[10px] font-bold uppercase tracking-wider">
                  Identity
                </th>
                <th className="text-right py-3 px-4 text-slate-500 text-[10px] font-bold uppercase tracking-wider">
                  Value
                </th>
                <th className="text-right py-3 px-4 text-slate-500 text-[10px] font-bold uppercase tracking-wider">
                  Shift (24H)
                </th>
                <th className="text-right py-3 px-4 text-slate-500 text-[10px] font-bold uppercase tracking-wider">
                  Spread
                </th>
                <th className="text-right py-3 px-4 text-slate-500 text-[10px] font-bold uppercase tracking-wider">
                  Volume
                </th>
                <th className="text-right py-3 px-4 text-slate-500 text-[10px] font-bold uppercase tracking-wider">
                  Trend
                </th>
              </tr>
            </thead>
            <tbody>
              {quotes.length > 0 ? (
                quotes.map((q) => (
                  <tr
                    key={q.symbol}
                    className={cn(
                      'border-b border-border-slate/50 hover:bg-white/5 cursor-pointer',
                      selectedSymbol === q.symbol && 'bg-cyber-orange/10'
                    )}
                    onClick={() => onSelectSymbol(q.symbol)}
                  >
                    <td className="py-3 px-4 font-mono font-bold text-white">{q.symbol}</td>
                    <td className="py-3 px-4 text-slate-400">
                      {SYMBOL_TO_COMPANY[q.symbol] || q.symbol}
                    </td>
                    <td className="py-3 px-4 text-right font-mono text-white">
                      ${q.close.toFixed(2)}
                    </td>
                    <td
                      className={cn(
                        'py-3 px-4 text-right font-mono font-bold',
                        (q.change_pct ?? 0) >= 0 ? 'text-cyber-orange' : 'text-rose-500'
                      )}
                    >
                      {(q.change_pct ?? 0) >= 0 ? '+' : ''}
                      {(q.change_pct ?? 0).toFixed(2)}%
                    </td>
                    <td className="py-3 px-4 text-right font-mono text-slate-400 text-xs">
                      {(q.bid ?? q.close - 0.01).toFixed(2)} // {(q.ask ?? q.close + 0.01).toFixed(2)}
                    </td>
                    <td className="py-3 px-4 text-right font-mono text-slate-400">
                      {(q.volume / 1e6).toFixed(1)}M
                    </td>
                    <td className="py-3 px-4 text-right">
                      <MiniSparkline quotes={getQuotesForSparkline(q.symbol)} />
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={7} className="py-8 text-center text-slate-500">
                    No live quotes yet. Waiting for feed.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  );
}

'use client';

import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { TrendingUp } from 'lucide-react';
import { cn } from '@/lib/utils';
import { StockQuote } from '@/types';
import { SYMBOL_TO_COMPANY } from '@/lib/constants';

interface StockQuoteHeaderProps {
  quote: StockQuote;
}

/** Header row for stock detail: symbol, exchange, company, price, change, action. No card wrapper. */
export function StockQuoteHeader({ quote }: StockQuoteHeaderProps) {
  const isUp = (quote.change_pct ?? 0) >= 0;
  return (
    <div className="flex flex-wrap items-start justify-between gap-6">
      <div>
        <div className="flex items-center gap-2 flex-wrap">
          <h2 className="text-3xl font-black tracking-tighter text-white">{quote.symbol}</h2>
          <Badge className="bg-border-slate text-slate-400 border-0 text-[10px] font-bold uppercase">
            NASDAQ
          </Badge>
        </div>
        <p className="text-slate-500 text-xs font-bold uppercase tracking-wider mt-1">
          {SYMBOL_TO_COMPANY[quote.symbol] || quote.symbol} // MARKET INTELLIGENCE FEED V2.4
        </p>
        <div className="flex items-baseline gap-3 mt-3">
          <span className="text-4xl font-black text-white font-mono">
            ${quote.close.toFixed(2)}
          </span>
          <span
            className={cn(
              'text-lg font-bold flex items-center gap-1',
              isUp ? 'text-cyber-orange' : 'text-rose-500'
            )}
          >
            {isUp ? (
              <TrendingUp className="h-5 w-5" />
            ) : (
              <span className="rotate-180">
                <TrendingUp className="h-5 w-5" />
              </span>
            )}
            {(quote.change ?? 0) >= 0 ? '+' : ''}
            {(quote.change ?? 0).toFixed(2)} (
            {(quote.change_pct ?? 0) >= 0 ? '+' : ''}
            {(quote.change_pct ?? 0).toFixed(2)}%)
          </span>
        </div>
      </div>
      <Button className="bg-cyber-orange text-midnight-slate hover:bg-cyber-orange/90 font-black text-xs uppercase tracking-wider">
        Execute Matrix Trade
      </Button>
    </div>
  );
}

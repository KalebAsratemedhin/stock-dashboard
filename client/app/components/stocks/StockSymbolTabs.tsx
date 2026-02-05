'use client';

import { Button } from '@/components/ui/button';
import { ScrollArea, ScrollBar } from '@/components/ui/scroll-area';
import { cn } from '@/lib/utils';
import { StockQuote } from '@/types';
import { STOCK_SYMBOLS } from '@/lib/constants';

interface StockSymbolTabsProps {
  symbols?: readonly string[];
  selectedSymbol: string;
  onSelect: (symbol: string) => void;
  latestQuotes: Record<string, StockQuote>;
}

export function StockSymbolTabs({
  symbols = STOCK_SYMBOLS,
  selectedSymbol,
  onSelect,
  latestQuotes,
}: StockSymbolTabsProps) {
  return (
    <ScrollArea className="w-full">
      <div className="flex gap-2 pb-2 whitespace-nowrap">
        {symbols.map((symbol) => {
          const quote = latestQuotes[symbol];
          const isSelected = symbol === selectedSymbol;
          return (
            <Button
              key={symbol}
              onClick={() => onSelect(symbol)}
              variant={isSelected ? 'default' : 'outline'}
              size="sm"
              className={cn(
                'shrink-0 bg-surface-slate border-border-slate text-node-white hover:bg-surface-slate hover:text-cyber-orange',
                isSelected &&
                  'bg-cyber-orange/10 border-cyber-orange text-cyber-orange hover:bg-cyber-orange/10 hover:text-cyber-orange'
              )}
            >
              <span className="font-bold">{symbol}</span>
              {quote && (
                <>
                  <span className="ml-2 font-mono text-white">${quote.close.toFixed(2)}</span>
                  <span
                    className={cn(
                      'ml-2 text-xs font-bold',
                      (quote.change_pct || 0) >= 0 ? 'text-cyan-accent' : 'text-rose-500'
                    )}
                  >
                    {(quote.change_pct || 0) >= 0 ? '+' : ''}
                    {quote.change_pct?.toFixed(2)}%
                  </span>
                </>
              )}
            </Button>
          );
        })}
      </div>
      <ScrollBar orientation="horizontal" />
    </ScrollArea>
  );
}

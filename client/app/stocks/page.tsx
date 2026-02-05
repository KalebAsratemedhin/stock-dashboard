'use client';

import { useEffect, useState, useMemo } from 'react';
import Layout from '../components/Layout';
import { PageHeader } from '../components/PageHeader';
import { EmptyState } from '../components/EmptyState';
import {
  StockSymbolTabs,
  StockQuoteHeader,
  StockChartControls,
  StockChart,
  StockWatchlistTable,
  type StockChartDataPoint,
} from '../components/stocks';
import { useWebSocket } from '@/hooks/useWebSocket';
import { useStockStore } from '@/lib/store';
import { fetchStocks } from '@/lib/api';
import { StockQuote } from '@/types';
import { Card, CardContent } from '@/components/ui/card';
import { ema } from '@/lib/chart-utils';
import { STOCK_SYMBOLS } from '@/lib/constants';

const CHART_POINTS = 35;

export default function StocksPage() {
  const { isConnected, subscribe } = useWebSocket();
  const { latestQuotes, quotes, setQuotes } = useStockStore();
  const [selectedSymbol, setSelectedSymbol] = useState('AAPL');
  const [localQuotes, setLocalQuotes] = useState<StockQuote[]>([]);
  const [timeRange, setTimeRange] = useState('1D');
  const [chartType, setChartType] = useState<'bar' | 'line'>('bar');

  useEffect(() => {
    subscribe('stock_quotes');
  }, [subscribe]);

  useEffect(() => {
    const loadData = async () => {
      try {
        const data = await fetchStocks(selectedSymbol, 150);
        const fetchedQuotes = Array.isArray(data) ? data : [];
        if (fetchedQuotes.length > 0) {
          setLocalQuotes(fetchedQuotes);
          setQuotes(fetchedQuotes);
        } else {
          const allLatest = await fetchStocks();
          if (Array.isArray(allLatest) && allLatest.length > 0) {
            setQuotes(allLatest);
            const symbolQuotes = allLatest.filter((q) => q.symbol === selectedSymbol);
            setLocalQuotes(symbolQuotes.length > 0 ? symbolQuotes : []);
          } else {
            setLocalQuotes([]);
          }
        }
      } catch (error) {
        console.error('Error loading stocks:', error);
        setLocalQuotes([]);
      }
    };
    loadData();
  }, [selectedSymbol, setQuotes]);

  const displayQuotes =
    quotes.length > 0 ? quotes.filter((q) => q.symbol === selectedSymbol) : localQuotes;
  const currentQuote =
    latestQuotes[selectedSymbol] ||
    (displayQuotes.length > 0 ? displayQuotes[displayQuotes.length - 1] : null);

  const chartData = useMemo((): StockChartDataPoint[] => {
    const slice = displayQuotes.slice(-CHART_POINTS);
    const closes = slice.map((q) => q.close);
    const ema20 = ema(closes, 20);
    const ema50 = ema(closes, 50);
    return slice.map((quote, i) => ({
      index: i,
      time: new Date(quote.timestamp).toLocaleTimeString('en-US', {
        hour: '2-digit',
        minute: '2-digit',
        hour12: false,
      }),
      close: quote.close,
      open: quote.open,
      high: quote.high,
      low: quote.low,
      isUp: quote.close >= quote.open,
      ema20: ema20[i],
      ema50: ema50[i],
    }));
  }, [displayQuotes]);

  const watchlistQuotes = useMemo(() => {
    return STOCK_SYMBOLS.map((symbol) => latestQuotes[symbol] || null).filter(
      Boolean
    ) as StockQuote[];
  }, [latestQuotes]);

  const symbolQuotesForSparkline = (symbol: string) => {
    const fromDisplay = displayQuotes.filter((q) => q.symbol === symbol);
    if (fromDisplay.length > 0) return fromDisplay;
    return (quotes.length ? quotes : localQuotes).filter((q) => q.symbol === symbol);
  };

  return (
    <Layout>
      <div className="p-8 flex flex-col gap-8">
        <PageHeader
          title="Market Intelligence"
          subtitle="Real-Time Stock Market Data"
          isLive={isConnected}
        />

        <StockSymbolTabs
          selectedSymbol={selectedSymbol}
          onSelect={setSelectedSymbol}
          latestQuotes={latestQuotes}
        />

        {currentQuote ? (
          <>
            <Card className="bg-surface-slate border-border-slate overflow-hidden">
              <CardContent className="p-6">
                <StockQuoteHeader quote={currentQuote} />
                <StockChartControls
                  timeRange={timeRange}
                  onTimeRangeChange={setTimeRange}
                  chartType={chartType}
                  onChartTypeChange={setChartType}
                />
                {chartData.length > 0 && (
                  <StockChart data={chartData} chartType={chartType} />
                )}
              </CardContent>
            </Card>

            <StockWatchlistTable
              quotes={watchlistQuotes}
              selectedSymbol={selectedSymbol}
              onSelectSymbol={setSelectedSymbol}
              getQuotesForSparkline={symbolQuotesForSparkline}
            />
          </>
        ) : (
          <EmptyState
            title="No stock data available"
            description="Waiting for data generation..."
            hint="Stock data is generated every 30 seconds. Ensure simulator and worker are running."
          />
        )}
      </div>
    </Layout>
  );
}

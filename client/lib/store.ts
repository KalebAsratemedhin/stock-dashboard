import { create } from 'zustand';
import { StockQuote, Sale, UserEvent, FinancialMetric } from '@/types';

interface StockStore {
  quotes: StockQuote[];
  latestQuotes: Record<string, StockQuote>;
  addQuote: (quote: StockQuote) => void;
  setQuotes: (quotes: StockQuote[]) => void;
  updateLatestQuotes: () => void;
}

interface SalesStore {
  sales: Sale[];
  totalRevenue: number;
  addSale: (sale: Sale) => void;
  setSales: (sales: Sale[]) => void;
  setRevenue: (revenue: number) => void;
  updateRevenue: () => void;
}

interface UserEventsStore {
  events: UserEvent[];
  addEvent: (event: UserEvent) => void;
  setEvents: (events: UserEvent[]) => void;
}

interface FinancialStore {
  metrics: FinancialMetric[];
  addMetric: (metric: FinancialMetric) => void;
  setMetrics: (metrics: FinancialMetric[]) => void;
}

export const useStockStore = create<StockStore>((set, get) => ({
  quotes: [],
  latestQuotes: {},
  addQuote: (quote) =>
    set((state) => {
      const newQuotes = [...state.quotes.slice(-999), quote];
      const latestQuotes = { ...state.latestQuotes, [quote.symbol]: quote };
      return {
        quotes: newQuotes,
        latestQuotes,
      };
    }),
  setQuotes: (quotes) => {
    const latestQuotes: Record<string, StockQuote> = {};
    quotes.forEach((quote) => {
      if (!latestQuotes[quote.symbol] || new Date(quote.timestamp) > new Date(latestQuotes[quote.symbol].timestamp)) {
        latestQuotes[quote.symbol] = quote;
      }
    });
    set({ quotes, latestQuotes });
  },
  updateLatestQuotes: () => {
    const { quotes } = get();
    const latestQuotes: Record<string, StockQuote> = {};
    quotes.forEach((quote) => {
      if (!latestQuotes[quote.symbol] || new Date(quote.timestamp) > new Date(latestQuotes[quote.symbol].timestamp)) {
        latestQuotes[quote.symbol] = quote;
      }
    });
    set({ latestQuotes });
  },
}));

export const useSalesStore = create<SalesStore>((set, get) => ({
  sales: [],
  totalRevenue: 0,
  addSale: (sale) =>
    set((state) => {
      const newSales = [...state.sales.slice(-999), sale];
      const totalRevenue = newSales.reduce((sum, s) => sum + s.revenue, 0);
      return {
        sales: newSales,
        totalRevenue,
      };
    }),
  setSales: (sales) => {
    const totalRevenue = sales.reduce((sum, s) => sum + s.revenue, 0);
    set({ sales, totalRevenue });
  },
  setRevenue: (revenue) => set({ totalRevenue: revenue }),
  updateRevenue: () => {
    const { sales } = get();
    const totalRevenue = sales.reduce((sum, s) => sum + s.revenue, 0);
    set({ totalRevenue });
  },
}));

export const useUserEventsStore = create<UserEventsStore>((set) => ({
  events: [],
  addEvent: (event) =>
    set((state) => ({
      events: [...state.events.slice(-999), event],
    })),
  setEvents: (events) => set({ events }),
}));

export const useFinancialStore = create<FinancialStore>((set) => ({
  metrics: [],
  addMetric: (metric) =>
    set((state) => ({
      metrics: [...state.metrics.slice(-999), metric],
    })),
  setMetrics: (metrics) => set({ metrics }),
}));

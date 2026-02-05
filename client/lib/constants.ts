export const STOCK_SYMBOLS = [
  'AAPL',
  'GOOGL',
  'MSFT',
  'AMZN',
  'TSLA',
  'META',
  'NVDA',
  'NFLX',
  'AMD',
  'INTC',
] as const;

export const SYMBOL_TO_COMPANY: Record<string, string> = {
  AAPL: 'APPLE INC.',
  GOOGL: 'ALPHABET INC.',
  MSFT: 'MICROSOFT CORP.',
  AMZN: 'AMAZON.COM INC.',
  TSLA: 'TESLA INC.',
  META: 'META PLATFORMS INC.',
  NVDA: 'NVIDIA CORP.',
  NFLX: 'NETFLIX INC.',
  AMD: 'ADVANCED MICRO DEVICES',
  INTC: 'INTEL CORP.',
};

export const TIME_HORIZONS = ['1D', '1W', '1M', '1Y', '5Y', 'MAX'] as const;

export const CHART_TOOLTIP_STYLE = {
  backgroundColor: '#1A1D23',
  border: '1px solid #2D323B',
  borderRadius: '0.5rem',
  color: '#F8FAFC',
} as const;

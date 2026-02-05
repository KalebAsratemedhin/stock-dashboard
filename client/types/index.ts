export interface StockQuote {
  id: number;
  symbol: string;
  timestamp: string;
  open: number;
  high: number;
  low: number;
  close: number;
  volume: number;
  bid?: number;
  ask?: number;
  change?: number;
  change_pct?: number;
  created_at: string;
  updated_at: string;
}

export interface Sale {
  id: number;
  timestamp: string;
  product_id: string;
  product_name: string;
  category?: string;
  customer_id?: string;
  region?: string;
  quantity: number;
  unit_price: number;
  discount: number;
  revenue: number;
  created_at: string;
  updated_at: string;
}

export interface UserEvent {
  id: number;
  timestamp: string;
  event_type: string;
  user_id?: string;
  session_id?: string;
  page?: string;
  device?: string;
  browser?: string;
  country?: string;
  city?: string;
  referrer?: string;
  metadata?: string;
  created_at: string;
  updated_at: string;
}

export interface FinancialMetric {
  id: number;
  timestamp: string;
  metric_type: string;
  department?: string;
  category?: string;
  amount: number;
  budget?: number;
  variance?: number;
  variance_pct?: number;
  period?: string;
  created_at: string;
  updated_at: string;
}

export interface WebSocketMessage {
  type: string;
  channel?: string;
  data?: any;
}




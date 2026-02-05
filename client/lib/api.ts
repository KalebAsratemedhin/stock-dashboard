const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export async function fetchStocks(symbol?: string, limit?: number) {
  const params = new URLSearchParams();
  if (symbol) params.append('symbol', symbol);
  if (limit) params.append('limit', limit.toString());
  
  const response = await fetch(`${API_URL}/api/stocks?${params}`);
  if (!response.ok) throw new Error('Failed to fetch stocks');
  return response.json();
}

export async function fetchStocksByRange(symbol: string, start: string, end: string) {
  const params = new URLSearchParams({
    symbol,
    start,
    end,
  });
  
  const response = await fetch(`${API_URL}/api/stocks/range?${params}`);
  if (!response.ok) throw new Error('Failed to fetch stocks');
  return response.json();
}

export async function fetchSales(start?: string, end?: string, category?: string, region?: string) {
  const params = new URLSearchParams();
  if (start) params.append('start', start);
  if (end) params.append('end', end);
  if (category) params.append('category', category);
  if (region) params.append('region', region);
  
  const response = await fetch(`${API_URL}/api/sales?${params}`);
  if (!response.ok) throw new Error('Failed to fetch sales');
  return response.json();
}

export async function fetchSalesRevenue(start: string, end: string) {
  const params = new URLSearchParams({ start, end });
  const response = await fetch(`${API_URL}/api/sales/revenue?${params}`);
  if (!response.ok) throw new Error('Failed to fetch revenue');
  return response.json();
}

export async function fetchUserEvents(start?: string, end?: string, type?: string) {
  const params = new URLSearchParams();
  if (start) params.append('start', start);
  if (end) params.append('end', end);
  if (type) params.append('type', type);
  
  const response = await fetch(`${API_URL}/api/events?${params}`);
  if (!response.ok) throw new Error('Failed to fetch events');
  return response.json();
}

export async function fetchFinancialMetrics(start?: string, end?: string, type?: string, department?: string) {
  const params = new URLSearchParams();
  if (start) params.append('start', start);
  if (end) params.append('end', end);
  if (type) params.append('type', type);
  if (department) params.append('department', department);
  
  const response = await fetch(`${API_URL}/api/metrics?${params}`);
  if (!response.ok) throw new Error('Failed to fetch metrics');
  return response.json();
}




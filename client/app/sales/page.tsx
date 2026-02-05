'use client';

import { useEffect, useState } from 'react';
import Layout from '../components/Layout';
import MetricCard from '../components/MetricCard';
import { PageHeader } from '../components/PageHeader';
import { ChartContainer } from '../components/ChartContainer';
import { useWebSocket } from '@/hooks/useWebSocket';
import { useSalesStore } from '@/lib/store';
import { fetchSales, fetchSalesRevenue } from '@/lib/api';
import { format } from 'date-fns';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ScrollArea } from '@/components/ui/scroll-area';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { ShoppingBag } from 'lucide-react';

export default function SalesPage() {
  const { isConnected, subscribe } = useWebSocket();
  const { sales, totalRevenue, setSales, setRevenue } = useSalesStore();
  const [timeRange, setTimeRange] = useState('1M');

  useEffect(() => {
    subscribe('sales');
  }, [subscribe]);

  useEffect(() => {
    const loadData = async () => {
      try {
        const end = new Date().toISOString();
        const start = new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString();
        const salesData = await fetchSales(start, end);
        const revenueData = await fetchSalesRevenue(start, end);
        setSales(Array.isArray(salesData) ? salesData : []);
        setRevenue(revenueData.revenue || 0);
      } catch (error) {
        console.error('Error loading sales:', error);
      }
    };
    loadData();
  }, [setSales, setRevenue]);

  const recentSales = sales.slice(-10);
  const avgOrderValue = sales.length > 0 ? totalRevenue / sales.length : 0;
  const chartData = Array.from({ length: 30 }, (_, i) => {
    const date = new Date(Date.now() - (29 - i) * 24 * 60 * 60 * 1000);
    const daySales = sales.filter((s) => new Date(s.timestamp).toDateString() === date.toDateString());
    return {
      date: format(date, 'MMM dd'),
      revenue: daySales.reduce((sum, s) => sum + s.revenue, 0),
      expenses: daySales.reduce((sum, s) => sum + s.revenue * 0.6, 0),
    };
  });

  return (
    <Layout>
      <div className="p-8 flex flex-col gap-8">
        <PageHeader
          title="Sales Radar"
          subtitle="Business Performance Analytics"
          isLive={isConnected}
        />

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <MetricCard
            title="Total Revenue"
            value={`$${(totalRevenue / 1000).toFixed(1)}k`}
            change={12.5}
            changeLabel="vs last month"
            icon="trending_up"
            trend="up"
          />
          <MetricCard
            title="Net Profit"
            value={`$${((totalRevenue * 0.34) / 1000).toFixed(1)}k`}
            change={8.4}
            changeLabel="vs last month"
            icon="insights"
            trend="up"
          />
          <MetricCard
            title="Avg Order Value"
            value={`$${avgOrderValue.toFixed(2)}`}
            change={5.1}
            changeLabel="vs last month"
            icon="shopping_cart"
            trend="up"
          />
        </div>

        <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
          <Card className="xl:col-span-2 bg-surface-slate border-border-slate">
            <CardHeader>
              <div className="flex justify-between items-center">
                <CardTitle className="text-lg font-black text-white">Revenue vs Expenses</CardTitle>
                <div className="flex gap-4">
                  <div className="flex items-center gap-2">
                    <div className="size-3 rounded-full bg-cyber-orange"></div>
                    <span className="text-xs font-medium opacity-70">Revenue</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <div className="size-3 rounded-full bg-amber-500"></div>
                    <span className="text-xs font-medium opacity-70">Expenses</span>
                  </div>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <ChartContainer height={256} withGrid>
                <ResponsiveContainer width="100%" height="100%" minHeight={256}>
                    <LineChart data={chartData}>
                      <CartesianGrid strokeDasharray="3 3" stroke="#2D323B" />
                      <XAxis 
                        dataKey="date" 
                        className="text-xs"
                        tick={{ fill: '#94a3b8', fontSize: 10 }}
                      />
                      <YAxis 
                        className="text-xs"
                        tick={{ fill: '#94a3b8', fontSize: 10 }}
                      />
                      <Tooltip 
                        contentStyle={{ 
                          backgroundColor: '#1A1D23',
                          border: '1px solid #2D323B',
                          borderRadius: '0.5rem',
                          color: '#F8FAFC'
                        }}
                      />
                      <Line 
                        type="monotone" 
                        dataKey="revenue" 
                        stroke="#FF8A00" 
                        strokeWidth={3}
                        name="Revenue"
                      />
                      <Line 
                        type="monotone" 
                        dataKey="expenses" 
                        stroke="#f59e0b" 
                        strokeWidth={2}
                        strokeDasharray="4 4"
                        name="Expenses"
                      />
                    </LineChart>
                </ResponsiveContainer>
              </ChartContainer>
            </CardContent>
          </Card>

          <Card className="bg-surface-slate border-border-slate flex flex-col h-full overflow-hidden">
            <CardHeader className="p-4 border-b border-border-slate flex justify-between items-center bg-slate-900/40">
              <CardTitle className="text-sm font-black text-cyber-orange tracking-tight">LIVE TRANSACTIONS</CardTitle>
              <div className="flex items-center gap-1.5">
                <span className="relative flex h-2 w-2">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-orange-400 opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-2 w-2 bg-cyber-orange"></span>
                </span>
                <span className="text-[10px] font-black text-cyber-orange uppercase tracking-widest">Active</span>
              </div>
            </CardHeader>
            <CardContent className="flex-1 p-0">
              <ScrollArea className="h-[300px]">
                <div className="divide-y divide-border-slate">
                  {recentSales.length > 0 ? (
                    recentSales.map((sale) => (
                      <div key={sale.id} className="p-4 flex items-center gap-4 hover:bg-orange-500/5 transition-colors group">
                        <div className="size-10 rounded-full bg-cyber-orange/10 flex items-center justify-center text-cyber-orange group-hover:bg-cyber-orange group-hover:text-white transition-all duration-300">
                          <ShoppingBag className="h-5 w-5" />
                        </div>
                        <div className="flex-1 min-w-0">
                          <p className="text-sm font-black truncate text-white">#{sale.id} - {sale.product_name}</p>
                          <p className="text-[11px] font-medium text-slate-500 truncate">
                            {sale.customer_id || 'Anonymous'} • {format(new Date(sale.timestamp), 'MMM dd, HH:mm')}
                          </p>
                        </div>
                        <div className="text-right flex-shrink-0">
                          <p className="text-sm font-black text-cyber-orange">${sale.revenue.toFixed(2)}</p>
                          <p className="text-[10px] font-black text-cyan-accent uppercase">Paid</p>
                        </div>
                      </div>
                    ))
                  ) : (
                    <p className="text-slate-500 text-sm text-center py-8">No recent sales</p>
                  )}
                </div>
              </ScrollArea>
            </CardContent>
          </Card>
        </div>
      </div>
    </Layout>
  );
}

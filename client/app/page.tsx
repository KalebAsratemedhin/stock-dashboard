'use client';

import { useEffect } from 'react';
import Layout from './components/Layout';
import MetricCard from './components/MetricCard';
import { PageHeader } from './components/PageHeader';
import { ChartContainer } from './components/ChartContainer';
import { useWebSocket } from '@/hooks/useWebSocket';
import { useStockStore, useSalesStore, useUserEventsStore, useFinancialStore } from '@/lib/store';
import { fetchStocks, fetchSalesRevenue, fetchUserEvents, fetchFinancialMetrics } from '@/lib/api';
import { formatEventType } from '@/lib/utils';
import { format } from 'date-fns';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { ScrollArea } from '@/components/ui/scroll-area';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { Radio, CreditCard, MousePointerClick, LogIn, Activity } from 'lucide-react';

const getEventIcon = (eventType: string) => {
  switch (eventType) {
    case 'purchase':
      return <CreditCard className="h-4 w-4" />;
    case 'page_view':
      return <MousePointerClick className="h-4 w-4" />;
    case 'login':
      return <LogIn className="h-4 w-4" />;
    default:
      return <Activity className="h-4 w-4" />;
  }
};

export default function Dashboard() {
  const { isConnected, subscribe } = useWebSocket();
  const { totalRevenue } = useSalesStore();
  const { events } = useUserEventsStore();

  useEffect(() => {
    subscribe('stock_quotes');
    subscribe('sales');
    subscribe('user_events');
    subscribe('financial_metrics');
  }, [subscribe]);

  useEffect(() => {
    const loadData = async () => {
      try {
        const end = new Date().toISOString();
        const start = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString();
        await fetchStocks();
        await fetchSalesRevenue(start, end);
        await fetchUserEvents(start, end);
        await fetchFinancialMetrics(start, end);
      } catch (error) {
        console.error('Error loading data:', error);
      }
    };
    loadData();
  }, []);

  const recentEvents = events.slice(-5);
  const chartData = Array.from({ length: 30 }, (_, i) => ({
    date: format(new Date(Date.now() - (29 - i) * 24 * 60 * 60 * 1000), 'MMM dd'),
    value: Math.random() * 100 + 50,
  }));

  return (
    <Layout>
      <div className="p-8 flex flex-col gap-8">
        <PageHeader
          title="Executive Overview"
          subtitle={`Real-Time Analytics Dashboard // ${format(new Date(), 'MMM dd, yyyy HH:mm:ss')}`}
          isLive={isConnected}
        />

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <MetricCard
            title="Total Revenue"
            value={`$${(totalRevenue / 1000).toFixed(1)}k`}
            change={12.5}
            changeLabel="vs last month"
            icon="database"
            trend="up"
          />
          <MetricCard
            title="Active Users"
            value={`${events.length > 0 ? (events.length / 100).toFixed(1) : '0'}k`}
            change={5.0}
            changeLabel="vs last month"
            icon="sensors"
            trend="up"
          />
          <MetricCard
            title="Total Sales"
            value={recentEvents.filter(e => e.event_type === 'purchase').length.toString()}
            change={-2.0}
            changeLabel="vs last month"
            icon="shopping_cart"
            trend="down"
          />
          <MetricCard
            title="System Load"
            value="18.4%"
            change={0.5}
            changeLabel="vs last hour"
            icon="memory"
            trend="up"
          />
        </div>

        <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
          <Card className="xl:col-span-2 bg-surface-slate border-border-slate">
            <CardHeader className="flex justify-between items-center">
              <CardTitle className="text-white text-xs font-black tracking-widest uppercase">Market Performance</CardTitle>
              <div className="flex gap-px bg-border-slate border border-border-slate rounded-sm">
                <Button variant="ghost" size="sm" className="px-3 py-1.5 text-[10px] font-bold hover:text-cyber-orange h-auto">1D</Button>
                <Button variant="default" size="sm" className="px-3 py-1.5 bg-cyber-orange text-midnight-slate text-[10px] font-black h-auto">1W</Button>
                <Button variant="ghost" size="sm" className="px-3 py-1.5 text-[10px] font-bold hover:text-cyber-orange h-auto">1M</Button>
                <Button variant="ghost" size="sm" className="px-3 py-1.5 text-[10px] font-bold hover:text-cyber-orange h-auto">1Y</Button>
              </div>
            </CardHeader>
            <CardContent>
              <ChartContainer height={256} withGrid>
                <ResponsiveContainer width="100%" height="100%" minHeight={256}>
                    <AreaChart data={chartData}>
                      <defs>
                        <linearGradient id="colorValue" x1="0" y1="0" x2="0" y2="1">
                          <stop offset="5%" stopColor="#FF8A00" stopOpacity={0.3}/>
                          <stop offset="95%" stopColor="#FF8A00" stopOpacity={0}/>
                        </linearGradient>
                      </defs>
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
                      <Area 
                        type="monotone" 
                        dataKey="value" 
                        stroke="#FF8A00" 
                        fillOpacity={1} 
                        fill="url(#colorValue)" 
                        strokeWidth={2}
                      />
                    </AreaChart>
                </ResponsiveContainer>
              </ChartContainer>
            </CardContent>
          </Card>

          <Card className="bg-surface-slate border-border-slate flex flex-col h-full overflow-hidden">
            <CardHeader className="p-6 border-b border-border-slate flex items-center justify-between">
              <div>
                <CardTitle className="text-xs font-bold uppercase tracking-widest text-white">Live Activity</CardTitle>
                <p className="text-[10px] text-cyber-orange font-mono uppercase">REALTIME_TRACKING</p>
              </div>
              <Radio className="h-4 w-4 text-cyan-accent animate-pulse" />
            </CardHeader>
            <CardContent className="flex-1 p-4">
              <ScrollArea className="h-full">
                <div className="flex flex-col gap-4">
                  {recentEvents.length > 0 ? (
                    recentEvents.map((event) => (
                      <div key={event.id} className="flex gap-4 p-3 rounded-sm hover:bg-midnight-slate/50 border-l-2 border-cyber-orange/20 hover:border-cyber-orange transition-all group">
                        <div className="size-8 rounded-sm bg-cyber-orange/10 text-cyber-orange flex items-center justify-center shrink-0">
                          {getEventIcon(event.event_type)}
                        </div>
                        <div className="flex-1 min-w-0">
                          <p className="text-[11px] font-bold uppercase truncate">{formatEventType(event.event_type)}</p>
                          <p className="text-[10px] text-slate-500 font-mono truncate">{event.user_id || 'Anonymous'}</p>
                          <p className="text-[9px] font-bold text-cyan-accent mt-1">
                            {format(new Date(event.timestamp), 'MMM dd, HH:mm')}
                          </p>
                        </div>
                      </div>
                    ))
                  ) : (
                    <p className="text-slate-500 text-sm text-center py-8">No recent activity</p>
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

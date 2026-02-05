'use client';

import { useEffect } from 'react';
import Layout from '../components/Layout';
import MetricCard from '../components/MetricCard';
import { PageHeader } from '../components/PageHeader';
import { useWebSocket } from '@/hooks/useWebSocket';
import { useFinancialStore } from '@/lib/store';
import { fetchFinancialMetrics } from '@/lib/api';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Database } from 'lucide-react';

export default function FinancialPage() {
  const { isConnected, subscribe } = useWebSocket();
  const { metrics, setMetrics } = useFinancialStore();

  useEffect(() => {
    subscribe('financial_metrics');
  }, [subscribe]);

  useEffect(() => {
    const loadData = async () => {
      try {
        const end = new Date().toISOString();
        const start = new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString();
        const metricsData = await fetchFinancialMetrics(start, end);
        setMetrics(Array.isArray(metricsData) ? metricsData : []);
      } catch (error) {
        console.error('Error loading metrics:', error);
      }
    };
    loadData();
  }, [setMetrics]);

  const revenue = metrics.filter((m) => m.metric_type === 'revenue').reduce((sum, m) => sum + m.amount, 0);
  const expenses = metrics.filter((m) => m.metric_type === 'expense').reduce((sum, m) => sum + m.amount, 0);
  const profit = revenue - expenses;
  const margin = revenue > 0 ? (profit / revenue) * 100 : 0;
  const departments = Array.from(new Set(metrics.map((m) => m.department).filter(Boolean)));

  return (
    <Layout>
      <div className="p-8 flex flex-col gap-8">
        <PageHeader
          title="Financial Performance"
          subtitle="Business & Financial Metrics"
          isLive={isConnected}
        />

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <MetricCard
            title="Total Revenue"
            value={`$${(revenue / 1000).toFixed(1)}k`}
            change={12.5}
            changeLabel="vs last month"
            icon="trending_up"
            trend="up"
          />
          <MetricCard
            title="Total Expenses"
            value={`$${(expenses / 1000).toFixed(1)}k`}
            change={-3.2}
            changeLabel="vs last month"
            icon="trending_down"
            trend="down"
          />
          <MetricCard
            title="Net Profit"
            value={`$${(profit / 1000).toFixed(1)}k`}
            change={8.4}
            changeLabel="vs last month"
            icon="insights"
            trend="up"
          />
          <MetricCard
            title="Profit Margin"
            value={`${margin.toFixed(1)}%`}
            change={2.1}
            changeLabel="vs last month"
            icon="percent"
            trend="up"
          />
        </div>

        <Card className="bg-surface-slate border-border-slate">
          <CardHeader className="px-8 py-5 border-b border-border-slate">
            <div className="flex items-center gap-3">
              <Database className="h-5 w-5 text-cyber-orange" />
              <CardTitle className="text-sm font-black tracking-widest uppercase text-white">Budget vs Actual</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="p-0">
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow className="text-[9px] uppercase tracking-[0.3em] text-slate-600 bg-midnight-slate font-black hover:bg-midnight-slate">
                    <TableHead className="px-8 py-4">Department</TableHead>
                    <TableHead className="px-8 py-4 text-right">Actual</TableHead>
                    <TableHead className="px-8 py-4 text-right">Budget</TableHead>
                    <TableHead className="px-8 py-4 text-right">Variance</TableHead>
                    <TableHead className="px-8 py-4 text-right">Variance %</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {departments.length > 0 ? (
                    departments.map((dept) => {
                      const deptMetrics = metrics.filter(m => m.department === dept);
                      const actual = deptMetrics.reduce((sum, m) => sum + m.amount, 0);
                      const budget = deptMetrics.reduce((sum, m) => sum + (m.budget || 0), 0);
                      const variance = actual - budget;
                      const variancePct = budget > 0 ? (variance / budget) * 100 : 0;
                      
                      return (
                        <TableRow key={dept} className="border-b border-border-slate hover:bg-white/[0.02] transition-all">
                          <TableCell className="px-8 py-4 text-white text-xs font-bold uppercase tracking-wide">{dept}</TableCell>
                          <TableCell className="px-8 py-4 text-right text-white font-black font-mono">${(actual / 1000).toFixed(1)}k</TableCell>
                          <TableCell className="px-8 py-4 text-right text-slate-400 font-mono">${(budget / 1000).toFixed(1)}k</TableCell>
                          <TableCell className={`px-8 py-4 text-right font-black font-mono ${variance >= 0 ? 'text-cyan-accent' : 'text-rose-500'}`}>
                            ${(variance / 1000).toFixed(1)}k
                          </TableCell>
                          <TableCell className={`px-8 py-4 text-right font-black ${variancePct >= 0 ? 'text-cyan-accent' : 'text-rose-500'}`}>
                            {variancePct >= 0 ? '+' : ''}{variancePct.toFixed(1)}%
                          </TableCell>
                        </TableRow>
                      );
                    })
                  ) : (
                    <TableRow>
                      <TableCell colSpan={5} className="px-8 py-8 text-center text-slate-500">
                        No financial data available
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>
      </div>
    </Layout>
  );
}

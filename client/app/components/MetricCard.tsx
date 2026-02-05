'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { cn } from '@/lib/utils';
import { Database, Activity, ShoppingCart, Cpu, TrendingUp, TrendingDown, DollarSign, Percent, BarChart3 } from 'lucide-react';

interface MetricCardProps {
  title: string;
  value: string | number;
  change?: number;
  changeLabel?: string;
  icon: string;
  trend?: 'up' | 'down';
}

const getIcon = (iconName: string) => {
  switch (iconName) {
    case 'database':
      return <Database className="h-5 w-5" />;
    case 'sensors':
      return <Activity className="h-5 w-5" />;
    case 'shopping_cart':
      return <ShoppingCart className="h-5 w-5" />;
    case 'memory':
      return <Cpu className="h-5 w-5" />;
    case 'trending_up':
      return <TrendingUp className="h-5 w-5" />;
    case 'trending_down':
      return <TrendingDown className="h-5 w-5" />;
    case 'insights':
      return <BarChart3 className="h-5 w-5" />;
    case 'percent':
      return <Percent className="h-5 w-5" />;
    default:
      return <Database className="h-5 w-5" />;
  }
};

export default function MetricCard({ title, value, change, changeLabel, icon, trend }: MetricCardProps) {
  const changeColor = trend === 'up' ? 'text-cyan-accent' : trend === 'down' ? 'text-rose-500' : 'text-slate-500';

  return (
    <Card className="bg-surface-slate border-border-slate hover:border-cyber-orange transition-colors">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-[10px] font-bold text-slate-500 uppercase tracking-widest">{title}</CardTitle>
        <div className="text-cyber-orange">{getIcon(icon)}</div>
      </CardHeader>
      <CardContent>
        <div className="flex items-end gap-2 mb-3">
          <h4 className="text-2xl font-black font-mono text-white">{value}</h4>
          {change !== undefined && (
            <span className={cn("text-[10px] font-bold mb-1 flex items-center gap-1", changeColor)}>
              {trend === 'up' && <TrendingUp className="h-3 w-3" />}
              {trend === 'down' && <TrendingDown className="h-3 w-3" />}
              {change > 0 ? '+' : ''}{change}% {changeLabel && <span className="text-slate-500 font-normal">{changeLabel}</span>}
            </span>
          )}
        </div>
        <div className="h-10 w-full opacity-60">
          <svg className="h-full w-full" viewBox="0 0 100 40">
            <path
              d="M0 35 L 20 25 L 40 30 L 60 15 L 80 20 L 100 5"
              fill="none"
              stroke={trend === 'up' ? '#00FFFF' : '#f43f5e'}
              strokeWidth="2"
            />
          </svg>
        </div>
      </CardContent>
    </Card>
  );
}

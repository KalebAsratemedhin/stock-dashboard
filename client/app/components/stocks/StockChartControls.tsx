'use client';

import { Button } from '@/components/ui/button';
import { BarChart3, LineChart as LineChartIcon } from 'lucide-react';
import { cn } from '@/lib/utils';
import { TIME_HORIZONS } from '@/lib/constants';

interface StockChartControlsProps {
  timeRange: string;
  onTimeRangeChange: (range: string) => void;
  chartType: 'bar' | 'line';
  onChartTypeChange: (type: 'bar' | 'line') => void;
}

export function StockChartControls({
  timeRange,
  onTimeRangeChange,
  chartType,
  onChartTypeChange,
}: StockChartControlsProps) {
  return (
    <div className="flex flex-wrap items-center gap-4 mt-6 pt-4 border-t border-border-slate">
      <div className="flex gap-px bg-midnight-slate border border-border-slate rounded-sm p-0.5">
        {TIME_HORIZONS.map((range) => (
          <Button
            key={range}
            onClick={() => onTimeRangeChange(range)}
            variant="ghost"
            size="sm"
            className={cn(
              'h-7 px-2.5 py-1 text-[10px] font-bold text-slate-400 hover:text-cyber-orange',
              timeRange === range &&
                'bg-cyber-orange text-midnight-slate hover:bg-cyber-orange hover:text-midnight-slate'
            )}
          >
            {range}
          </Button>
        ))}
      </div>
      <div className="flex gap-1">
        <Button
          variant="ghost"
          size="sm"
          className={cn('h-7 w-7 p-0', chartType === 'bar' && 'bg-cyber-orange/20 text-cyber-orange')}
          onClick={() => onChartTypeChange('bar')}
        >
          <BarChart3 className="h-4 w-4" />
        </Button>
        <Button
          variant="ghost"
          size="sm"
          className={cn('h-7 w-7 p-0', chartType === 'line' && 'bg-cyber-orange/20 text-cyber-orange')}
          onClick={() => onChartTypeChange('line')}
        >
          <LineChartIcon className="h-4 w-4" />
        </Button>
      </div>
      <div className="flex items-center gap-3 text-[10px] font-bold text-slate-500">
        <span className="flex items-center gap-1.5">
          <span className="w-1.5 h-1.5 rounded-full bg-cyber-orange" />
          EMA 20
        </span>
        <span className="flex items-center gap-1.5">
          <span className="w-1.5 h-1.5 rounded-full bg-slate-500" />
          EMA 50
        </span>
      </div>
    </div>
  );
}

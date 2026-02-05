'use client';

import {
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  ComposedChart,
  ReferenceArea,
} from 'recharts';
import { CHART_TOOLTIP_STYLE } from '@/lib/constants';

export interface StockChartDataPoint {
  index: number;
  time: string;
  close: number;
  open: number;
  high: number;
  low: number;
  isUp: boolean;
  ema20: number;
  ema50: number;
}

interface StockChartProps {
  data: StockChartDataPoint[];
  chartType: 'bar' | 'line';
  barHalfWidth?: number;
}

export function StockChart({
  data,
  chartType,
  barHalfWidth = 0.32,
}: StockChartProps) {
  if (data.length === 0) return null;

  return (
    <div className="h-[380px] w-full min-w-0 relative bg-midnight-slate/80 rounded-sm p-4 mt-4 overflow-hidden border border-border-slate">
      <ResponsiveContainer width="100%" height="100%" minHeight={360}>
        <ComposedChart data={data} margin={{ top: 8, right: 8, left: 8, bottom: 8 }}>
          <CartesianGrid strokeDasharray="3 3" stroke="#2D323B" />
          <XAxis
            dataKey="index"
            type="number"
            domain={[-0.5, data.length - 0.5]}
            tickFormatter={(i) => data[i]?.time ?? ''}
            tick={{ fill: '#94a3b8', fontSize: 10 }}
            axisLine={{ stroke: '#2D323B' }}
            allowDuplicatedCategory={false}
          />
          <YAxis
            domain={['auto', 'auto']}
            tick={{ fill: '#94a3b8', fontSize: 10 }}
            axisLine={{ stroke: '#2D323B' }}
            tickFormatter={(v) => `$${v.toFixed(0)}`}
          />
          <Tooltip
            contentStyle={CHART_TOOLTIP_STYLE}
            labelFormatter={(_, payload) =>
              payload?.[0]?.payload?.time ? `Time: ${payload[0].payload.time}` : ''
            }
            formatter={(value: number | undefined) => [`$${Number(value ?? 0).toFixed(2)}`, '']}
          />
          {chartType === 'bar' ? (
            data.map((d, i) => (
              <ReferenceArea
                key={i}
                x1={i - barHalfWidth}
                x2={i + barHalfWidth}
                y1={d.low}
                y2={d.high}
                fill={d.isUp ? '#FF9500' : '#FFFFFF'}
                fillOpacity={1}
                stroke={d.isUp ? '#FF9500' : '#E2E8F0'}
                strokeWidth={1}
              />
            ))
          ) : (
            <Line
              type="monotone"
              dataKey="close"
              stroke="var(--cyber-orange)"
              strokeWidth={2}
              dot={false}
            />
          )}
          <Line
            type="monotone"
            dataKey="ema20"
            stroke="var(--cyber-orange)"
            strokeWidth={1}
            strokeDasharray="4 2"
            dot={false}
          />
          <Line
            type="monotone"
            dataKey="ema50"
            stroke="#64748b"
            strokeWidth={1}
            strokeDasharray="4 2"
            dot={false}
          />
        </ComposedChart>
      </ResponsiveContainer>
    </div>
  );
}

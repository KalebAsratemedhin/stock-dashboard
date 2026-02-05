'use client';

import { ReactNode } from 'react';
import { cn } from '@/lib/utils';

interface ChartContainerProps {
  children: ReactNode;
  className?: string;
  height?: string | number;
  withGrid?: boolean;
}

export function ChartContainer({
  children,
  className,
  height = 256,
  withGrid = false,
}: ChartContainerProps) {
  const h = typeof height === 'number' ? `${height}px` : height;
  return (
    <div
      className={cn(
        'w-full min-w-0 relative bg-midnight-slate/80 rounded-sm p-4 overflow-hidden border border-border-slate',
        className
      )}
      style={{ height: h }}
    >
      {withGrid && (
        <div className="absolute inset-0 opacity-[0.05] pointer-events-none grid-bg" />
      )}
      <div className="relative z-10 h-full">{children}</div>
    </div>
  );
}

'use client';

import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';

interface PageHeaderProps {
  title: string;
  subtitle?: string;
  isLive?: boolean;
  className?: string;
}

export function PageHeader({ title, subtitle, isLive = false, className }: PageHeaderProps) {
  return (
    <div className={cn('flex flex-wrap justify-between items-end gap-6', className)}>
      <div>
        <h1 className="text-white text-4xl font-black tracking-tighter uppercase italic mb-2">
          {title}
        </h1>
        {subtitle && (
          <p className="text-slate-500 text-[10px] font-black uppercase tracking-[0.2em]">
            {subtitle}
          </p>
        )}
      </div>
      <Badge
        variant={isLive ? 'default' : 'secondary'}
        className="gap-2 bg-black/40 border-border-slate text-node-white hover:bg-black/40 shrink-0"
      >
        <div
          className={cn('w-2 h-2 rounded-full', isLive ? 'bg-cyan-accent animate-pulse' : 'bg-slate-500')}
        />
        {isLive ? 'Live' : 'Disconnected'}
      </Badge>
    </div>
  );
}

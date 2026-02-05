'use client';

import { Card, CardContent } from '@/components/ui/card';
import { cn } from '@/lib/utils';

interface EmptyStateProps {
  title: string;
  description?: string;
  hint?: string;
  className?: string;
}

export function EmptyState({ title, description, hint, className }: EmptyStateProps) {
  return (
    <Card className={cn('bg-surface-slate border-border-slate', className)}>
      <CardContent className="py-12">
        <div className="text-center text-slate-500">
          <p className="text-lg mb-2">{title}</p>
          {description && <p className="text-sm">{description}</p>}
          {hint && <p className="text-xs mt-4 text-slate-600">{hint}</p>}
        </div>
      </CardContent>
    </Card>
  );
}

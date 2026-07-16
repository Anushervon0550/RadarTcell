import type { ReactNode } from 'react';
import { cn } from '../lib/cn';

export interface EmptyStateProps {
  title?: string;
  description?: ReactNode;
  action?: ReactNode;
  icon?: ReactNode;
  className?: string;
}

export function EmptyState({ title = 'Пусто', description, action, icon, className }: EmptyStateProps) {
  return (
    <div
      className={cn(
        'flex flex-col items-center justify-center gap-3 rounded-card border border-dashed border-line-soft bg-bg-panel/40 px-6 py-10 text-center',
        className,
      )}
    >
      {icon && <div className="text-3xl text-ink-muted">{icon}</div>}
      <div className="text-base font-medium text-ink">{title}</div>
      {description && <div className="max-w-md text-sm text-ink-muted">{description}</div>}
      {action}
    </div>
  );
}

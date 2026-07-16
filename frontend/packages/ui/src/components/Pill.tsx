import type { HTMLAttributes, ReactNode } from 'react';
import { cn } from '../lib/cn';

export interface PillProps extends HTMLAttributes<HTMLSpanElement> {
  variant?: 'default' | 'accent' | 'success' | 'danger' | 'warn';
  children: ReactNode;
}

const VARIANTS = {
  default: 'border-line-soft/80 bg-bg-panel/60 text-ink-muted',
  accent: 'border-brand-600/50 bg-brand-600/15 text-brand-400',
  success: 'border-emerald-500/40 bg-emerald-500/15 text-emerald-300',
  danger: 'border-red-500/40 bg-red-500/15 text-red-300',
  warn: 'border-amber-500/40 bg-amber-500/15 text-amber-300',
} as const;

export function Pill({ variant = 'default', className, children, ...rest }: PillProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center gap-1 rounded-full border px-2.5 py-0.5 text-xs font-medium',
        VARIANTS[variant],
        className,
      )}
      {...rest}
    >
      {children}
    </span>
  );
}

import { cn } from '../lib/cn';

export function TrlChip({ trl, className }: { trl: number; className?: string }) {
  return (
    <span
      className={cn(
        'inline-flex items-center gap-1 rounded-full border border-brand-600/45 bg-brand-600/15 px-2.5 py-0.5 text-xs font-semibold text-brand-400',
        className,
      )}
    >
      TRL {trl}
    </span>
  );
}

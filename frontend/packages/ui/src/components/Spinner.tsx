import { cn } from '../lib/cn';

export function Spinner({ className, size = 32 }: { className?: string; size?: number }) {
  return (
    <span
      role="status"
      aria-label="Loading"
      style={{ width: size, height: size }}
      className={cn(
        'inline-block animate-spin rounded-full border-[3px] border-brand-500/25 border-t-brand-400',
        className,
      )}
    />
  );
}

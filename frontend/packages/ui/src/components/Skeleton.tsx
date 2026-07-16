import { cn } from '../lib/cn';

export function Skeleton({ className }: { className?: string }) {
  return (
    <div
      className={cn(
        'animate-pulse rounded-md bg-gradient-to-r from-bg-panel via-bg-panel-2 to-bg-panel',
        className,
      )}
    />
  );
}

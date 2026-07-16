import { Skeleton } from '@radartcell/ui';

export function LoadingGrid({ height = 220 }: { height?: number }) {
  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      {Array.from({ length: 6 }).map((_, i) => (
        <Skeleton key={i} style={{ height }} className="w-full" />
      ))}
    </div>
  );
}

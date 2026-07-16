import { Link } from 'react-router-dom';
import { cssUrl, EmptyState, PageHeader, Pill } from '@radartcell/ui';
import { useTrends } from '@/api/queries';
import { coverForTrend } from '@/pages/radar/constants';
import { LoadingGrid } from './_LoadingGrid';

export function TrendsPage() {
  const q = useTrends();
  return (
    <>
      <PageHeader title="Тренды" description="Ключевые направления развития технологий." />
      {q.isLoading ? (
        <LoadingGrid />
      ) : (q.data ?? []).length === 0 ? (
        <EmptyState title="Пусто" description="Тренды пока не заведены." />
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {(q.data ?? []).map((t) => (
            <Link
              key={t.id}
              to={`/trend/${encodeURIComponent(t.slug)}`}
              className="group flex flex-col overflow-hidden rounded-card border border-line bg-bg-panel/95 shadow-elevated transition-all hover:-translate-y-0.5 hover:border-brand-600/50"
            >
              <div
                className="aspect-[16/8] w-full bg-cover bg-center"
                style={{ backgroundImage: cssUrl(coverForTrend(t.slug)) }}
              />
              <div className="flex flex-col gap-2 p-4">
                <h3 className="text-base font-semibold">{t.name}</h3>
                <div className="text-sm text-ink-muted">Технологий: {t.technologies_count}</div>
                <Pill>Открыть радар тренда</Pill>
              </div>
            </Link>
          ))}
        </div>
      )}
    </>
  );
}

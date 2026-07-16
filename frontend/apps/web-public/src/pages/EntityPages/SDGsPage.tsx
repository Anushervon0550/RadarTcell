import { Link } from 'react-router-dom';
import { Card, EmptyState, PageHeader } from '@radartcell/ui';
import { useSDGs } from '@/api/queries';
import { LoadingGrid } from './_LoadingGrid';

export function SDGsPage() {
  const q = useSDGs();
  return (
    <>
      <PageHeader title="Цели устойчивого развития" />
      {q.isLoading ? (
        <LoadingGrid height={140} />
      ) : (q.data ?? []).length === 0 ? (
        <EmptyState title="Пусто" />
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {(q.data ?? []).map((s) => (
            <Link key={s.id} to={`/sdg/${encodeURIComponent(s.code)}`}>
              <Card interactive>
                <h3 className="text-base font-semibold">{s.code}</h3>
                <div className="mt-1 text-sm text-ink-muted">{s.title}</div>
                <div className="mt-2 text-xs text-ink-muted-2">
                  Технологий: {s.technologies_count}
                </div>
              </Card>
            </Link>
          ))}
        </div>
      )}
    </>
  );
}

import { Link } from 'react-router-dom';
import { Card, EmptyState, PageHeader } from '@radartcell/ui';
import { useTags } from '@/api/queries';
import { LoadingGrid } from './_LoadingGrid';

export function TagsPage() {
  const q = useTags();
  return (
    <>
      <PageHeader title="Теги" description="Технологические категории и метки." />
      {q.isLoading ? (
        <LoadingGrid height={140} />
      ) : (q.data ?? []).length === 0 ? (
        <EmptyState title="Пусто" />
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {(q.data ?? []).map((t) => (
            <Link key={t.id} to={`/tag/${encodeURIComponent(t.slug)}`}>
              <Card interactive>
                <h3 className="text-base font-semibold">{t.title}</h3>
                {t.category && <div className="mt-1 text-sm text-ink-muted">{t.category}</div>}
                {t.description && (
                  <p className="mt-2 line-clamp-3 text-sm text-ink-muted">{t.description}</p>
                )}
              </Card>
            </Link>
          ))}
        </div>
      )}
    </>
  );
}

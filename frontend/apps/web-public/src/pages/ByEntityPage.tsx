import { Link, useParams } from 'react-router-dom';
import { Button, EmptyState, PageHeader, Skeleton } from '@radartcell/ui';
import { useTechnologiesByEntity } from '@/api/queries';
import { TechCard } from '@/components/TechCard';

const TITLES = {
  trend: 'Тренд',
  tag: 'Тег',
  sdg: 'ЦУР',
  organization: 'Организация',
} as const;

type Kind = keyof typeof TITLES;

export function ByEntityPage({ kind }: { kind: Kind }) {
  const params = useParams<{ value: string }>();
  const value = params.value ?? '';
  const displayValue = safeDecode(value);
  const { data, isLoading, error } = useTechnologiesByEntity(kind, value);

  const items = data?.items ?? [];

  return (
    <>
      <PageHeader
        title={`${TITLES[kind]}: ${displayValue}`}
        description={`Технологий найдено: ${items.length}`}
        actions={
          <Link to="/">
            <Button variant="ghost">← К радару</Button>
          </Link>
        }
      />
      {isLoading ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {Array.from({ length: 6 }).map((_, i) => (
            <Skeleton key={i} className="h-72 w-full" />
          ))}
        </div>
      ) : error ? (
        <EmptyState title="Ошибка" description={(error as Error).message} />
      ) : items.length === 0 ? (
        <EmptyState title="Нет связанных технологий" />
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {items.map((t) => (
            <TechCard key={t.id} item={t} />
          ))}
        </div>
      )}
    </>
  );
}

function safeDecode(s: string) {
  try {
    return decodeURIComponent(s);
  } catch {
    return s;
  }
}

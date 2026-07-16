import { useMemo, useState } from 'react';
import {
  Button,
  Card,
  EmptyState,
  Input,
  PageHeader,
  Select,
  Skeleton,
} from '@radartcell/ui';
import {
  useOrganizations,
  useSDGs,
  useTags,
  useTechnologies,
  useTrends,
} from '@/api/queries';
import { TechCard } from '@/components/TechCard';

interface Filters {
  search: string;
  trend_id: string;
  tag_id: string;
  sdg_id: string;
  organization_id: string;
  trl_min: string;
  trl_max: string;
}

const EMPTY: Filters = {
  search: '',
  trend_id: '',
  tag_id: '',
  sdg_id: '',
  organization_id: '',
  trl_min: '',
  trl_max: '',
};

function normalize(filters: Filters) {
  const out: Record<string, string | number> = {};
  if (filters.search) out.search = filters.search;
  if (filters.trend_id) out.trend_id = filters.trend_id;
  if (filters.tag_id) out.tag_id = filters.tag_id;
  if (filters.sdg_id) out.sdg_id = filters.sdg_id;
  if (filters.organization_id) out.organization_id = filters.organization_id;
  if (filters.trl_min) out.trl_min = Number(filters.trl_min);
  if (filters.trl_max) out.trl_max = Number(filters.trl_max);
  out.limit = 200;
  return out;
}

export function CatalogPage() {
  const [pending, setPending] = useState<Filters>(EMPTY);
  const [applied, setApplied] = useState<Filters>(EMPTY);

  const trendsQ = useTrends();
  const tagsQ = useTags();
  const sdgsQ = useSDGs();
  const orgsQ = useOrganizations();
  const listQ = useTechnologies(normalize(applied));

  const trendOptions = useMemo(
    () => (trendsQ.data ?? []).map((t) => ({ value: t.id, label: t.name })),
    [trendsQ.data],
  );
  const tagOptions = useMemo(
    () => (tagsQ.data ?? []).map((t) => ({ value: t.id, label: t.title || t.slug })),
    [tagsQ.data],
  );
  const sdgOptions = useMemo(
    () => (sdgsQ.data ?? []).map((s) => ({ value: s.id, label: s.code })),
    [sdgsQ.data],
  );
  const orgOptions = useMemo(
    () => (orgsQ.data ?? []).map((o) => ({ value: o.id, label: o.name })),
    [orgsQ.data],
  );

  const items = listQ.data?.items ?? [];
  const total = listQ.data?.total ?? items.length;

  return (
    <>
      <PageHeader
        title="Каталог технологий"
        description={
          listQ.isFetching ? 'Загружаем…' : `Найдено: ${items.length} из ${total}`
        }
      />

      <Card className="mb-4">
        <div className="flex flex-wrap items-end gap-2">
          <div className="min-w-[220px] flex-1">
            <Input
              placeholder="Поиск по названию"
              value={pending.search}
              onChange={(e) => setPending((p) => ({ ...p, search: e.target.value }))}
            />
          </div>
          <div className="min-w-[180px]">
            <Select
              placeholder="Все тренды"
              value={pending.trend_id}
              options={trendOptions}
              onChange={(e) => setPending((p) => ({ ...p, trend_id: e.target.value }))}
            />
          </div>
          <div className="min-w-[180px]">
            <Select
              placeholder="Все теги"
              value={pending.tag_id}
              options={tagOptions}
              onChange={(e) => setPending((p) => ({ ...p, tag_id: e.target.value }))}
            />
          </div>
          <div className="min-w-[160px]">
            <Select
              placeholder="Все ЦУР"
              value={pending.sdg_id}
              options={sdgOptions}
              onChange={(e) => setPending((p) => ({ ...p, sdg_id: e.target.value }))}
            />
          </div>
          <div className="min-w-[200px]">
            <Select
              placeholder="Все организации"
              value={pending.organization_id}
              options={orgOptions}
              onChange={(e) =>
                setPending((p) => ({ ...p, organization_id: e.target.value }))
              }
            />
          </div>
          <div className="w-[110px]">
            <Input
              type="number"
              min={1}
              max={9}
              placeholder="TRL мин"
              value={pending.trl_min}
              onChange={(e) => setPending((p) => ({ ...p, trl_min: e.target.value }))}
            />
          </div>
          <div className="w-[110px]">
            <Input
              type="number"
              min={1}
              max={9}
              placeholder="TRL макс"
              value={pending.trl_max}
              onChange={(e) => setPending((p) => ({ ...p, trl_max: e.target.value }))}
            />
          </div>
          <Button variant="primary" onClick={() => setApplied(pending)}>
            Применить
          </Button>
          <Button
            variant="ghost"
            onClick={() => {
              setPending(EMPTY);
              setApplied(EMPTY);
            }}
          >
            Сброс
          </Button>
        </div>
      </Card>

      {listQ.isLoading ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {Array.from({ length: 8 }).map((_, i) => (
            <Skeleton key={i} className="h-72 w-full" />
          ))}
        </div>
      ) : items.length === 0 ? (
        <EmptyState
          title="Ничего не найдено"
          description="Попробуйте изменить фильтры или сбросить их."
        />
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {items.map((it) => (
            <TechCard key={it.id} item={it} />
          ))}
        </div>
      )}
    </>
  );
}

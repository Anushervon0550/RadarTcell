import { useQueries } from '@tanstack/react-query';
import { Card, PageHeader, Skeleton } from '@radartcell/ui';
import { api } from '@/api/client';
import { unwrapList } from '@/crud/types';

interface Stat {
  label: string;
  value: number | string;
  hint?: string;
}

export function DashboardPage() {
  const results = useQueries({
    queries: [
      { queryKey: ['dash', 'trends'], queryFn: () => api.get('/api/admin/trends') },
      { queryKey: ['dash', 'tags'], queryFn: () => api.get('/api/admin/tags') },
      { queryKey: ['dash', 'orgs'], queryFn: () => api.get('/api/admin/organizations') },
      { queryKey: ['dash', 'metrics'], queryFn: () => api.get('/api/admin/metrics') },
      { queryKey: ['dash', 'sdgs'], queryFn: () => api.get('/api/admin/sdgs') },
      { queryKey: ['dash', 'techs'], queryFn: () => api.get('/api/admin/technologies?limit=1') },
      { queryKey: ['dash', 'users'], queryFn: () => api.get('/api/admin/users') },
    ],
  });

  const [trends, tags, orgs, metrics, sdgs, techs, users] = results;

  const loading = results.some((r) => r.isLoading);

  const stats: Stat[] = [
    { label: 'Технологий', value: (techs.data as { total?: number } | undefined)?.total ?? 0 },
    { label: 'Трендов', value: unwrapList<unknown>(trends.data).length },
    { label: 'Тегов', value: unwrapList<unknown>(tags.data).length },
    { label: 'Организаций', value: unwrapList<unknown>(orgs.data).length },
    { label: 'Метрик', value: unwrapList<unknown>(metrics.data).length },
    { label: 'ЦУР', value: unwrapList<unknown>(sdgs.data).length },
    { label: 'Админов', value: unwrapList<unknown>(users.data).length },
  ];

  return (
    <>
      <PageHeader
        title="Панель управления"
        description="Ключевые метрики контента и системы."
      />
      <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4 xl:grid-cols-7">
        {stats.map((s) => (
          <Card key={s.label} className="min-h-[110px]">
            <div className="text-xs uppercase tracking-[0.08em] text-ink-muted">{s.label}</div>
            <div className="mt-2 text-3xl font-bold tracking-tight">
              {loading ? <Skeleton className="h-8 w-16" /> : s.value}
            </div>
          </Card>
        ))}
      </div>
    </>
  );
}

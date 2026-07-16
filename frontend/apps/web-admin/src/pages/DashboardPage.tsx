import { useQueries } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { Card, PageHeader, Skeleton } from '@radartcell/ui';
import { api } from '@/api/client';
import { unwrapList } from '@/crud/types';

interface Stat {
  label: string;
  value: number | string;
  to: string;
  icon: string;
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
    {
      label: 'Технологий',
      value: (techs.data as { total?: number } | undefined)?.total ?? 0,
      to: '/technologies',
      icon: '⚙',
    },
    { label: 'Трендов', value: unwrapList<unknown>(trends.data).length, to: '/trends', icon: '↗' },
    { label: 'Тегов', value: unwrapList<unknown>(tags.data).length, to: '/tags', icon: '#' },
    {
      label: 'Организаций',
      value: unwrapList<unknown>(orgs.data).length,
      to: '/organizations',
      icon: '▣',
    },
    { label: 'Метрик', value: unwrapList<unknown>(metrics.data).length, to: '/metrics', icon: '⋯' },
    { label: 'ЦУР', value: unwrapList<unknown>(sdgs.data).length, to: '/sdgs', icon: '◯' },
    { label: 'Админов', value: unwrapList<unknown>(users.data).length, to: '/users', icon: '⚿' },
  ];

  return (
    <>
      <PageHeader title="Панель управления" description="Ключевые метрики контента и системы." />

      <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4 xl:grid-cols-7">
        {stats.map((s) => (
          <Link key={s.label} to={s.to} className="group">
            <Card className="min-h-[120px] transition-all group-hover:border-brand-600/50 group-hover:bg-brand-600/5">
              <div className="flex items-center justify-between">
                <div className="text-xs uppercase tracking-[0.08em] text-ink-muted">{s.label}</div>
                <span aria-hidden className="text-ink-muted-2 group-hover:text-brand-400">
                  {s.icon}
                </span>
              </div>
              <div className="mt-2 text-3xl font-bold tracking-tight">
                {loading ? <Skeleton className="h-8 w-16" /> : s.value}
              </div>
              <div className="mt-2 text-xs text-ink-muted-2 group-hover:text-brand-400">
                Открыть →
              </div>
            </Card>
          </Link>
        ))}
      </div>

      <div className="mt-6">
        <h2 className="mb-3 text-sm font-bold uppercase tracking-[0.1em] text-ink-muted-2">
          Быстрые действия
        </h2>
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          <QuickAction to="/technologies" title="Добавить технологию" desc="Новая карточка на радаре" />
          <QuickAction to="/trends" title="Управлять трендами" desc="Секторы радара" />
          <QuickAction to="/organizations" title="Организации" desc="Компании и логотипы" />
        </div>
      </div>
    </>
  );
}

function QuickAction({ to, title, desc }: { to: string; title: string; desc: string }) {
  return (
    <Link to={to} className="group">
      <Card className="transition-all group-hover:border-brand-600/50 group-hover:bg-brand-600/5">
        <div className="font-semibold group-hover:text-white">{title}</div>
        <div className="mt-1 text-sm text-ink-muted">{desc}</div>
      </Card>
    </Link>
  );
}

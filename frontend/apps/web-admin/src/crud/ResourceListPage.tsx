import { useMemo, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Button, Card, EmptyState, Input, PageHeader, Pill, Skeleton } from '@radartcell/ui';
import { api } from '@/api/client';
import { useToast } from '@/components/Toast';
import { ResourceForm } from './ResourceForm';
import { unwrapList, type CrudResource } from './types';

interface Props<T extends Record<string, unknown>> {
  resource: CrudResource<T>;
}

export function ResourceListPage<T extends Record<string, unknown>>({ resource }: Props<T>) {
  const qc = useQueryClient();
  const toast = useToast();
  const [editing, setEditing] = useState<T | null>(null);
  const [creating, setCreating] = useState<boolean>(false);
  const [search, setSearch] = useState('');

  const listKey = ['admin', resource.name, resource.listPath] as const;

  const listQ = useQuery({
    queryKey: listKey,
    queryFn: async () => {
      const data = await api.get<unknown>(resource.listPath);
      return resource.parseList ? resource.parseList(data) : unwrapList<T>(data);
    },
  });

  const createM = useMutation({
    mutationFn: (values: Record<string, unknown>) =>
      api.post(resource.createPath, values),
    onSuccess: () => {
      toast.success('Создано');
      setCreating(false);
      qc.invalidateQueries({ queryKey: ['admin', resource.name] });
    },
    onError: (e) => toast.error(`Ошибка создания: ${(e as Error).message}`),
  });

  const updateM = useMutation({
    mutationFn: ({ key, values }: { key: string; values: Record<string, unknown> }) =>
      api.put(resource.updatePath(key), values),
    onSuccess: () => {
      toast.success('Сохранено');
      setEditing(null);
      qc.invalidateQueries({ queryKey: ['admin', resource.name] });
    },
    onError: (e) => toast.error(`Ошибка сохранения: ${(e as Error).message}`),
  });

  const removeM = useMutation({
    mutationFn: (key: string) => api.delete(resource.removePath(key)),
    onSuccess: () => {
      toast.success('Удалено');
      qc.invalidateQueries({ queryKey: ['admin', resource.name] });
    },
    onError: (e) => toast.error(`Ошибка удаления: ${(e as Error).message}`),
  });

  const restoreM = useMutation({
    mutationFn: (key: string) => api.put(resource.restorePath!(key)),
    onSuccess: () => {
      toast.success('Восстановлено');
      qc.invalidateQueries({ queryKey: ['admin', resource.name] });
    },
    onError: (e) => toast.error(`Ошибка восстановления: ${(e as Error).message}`),
  });

  const rows = useMemo(() => {
    if (!listQ.data) return [] as T[];
    if (!search.trim()) return listQ.data;
    const q = search.trim().toLowerCase();
    return listQ.data.filter((r) =>
      Object.values(r).some((v) => String(v ?? '').toLowerCase().includes(q)),
    );
  }, [listQ.data, search]);

  return (
    <>
      <PageHeader
        title={resource.title}
        description={
          listQ.isLoading
            ? 'Загружаем…'
            : `Всего записей: ${listQ.data?.length ?? 0}`
        }
        actions={
          <Button
            variant="primary"
            onClick={() => {
              setEditing(null);
              setCreating(true);
            }}
          >
            + Создать
          </Button>
        }
      />

      <div className="mb-4 flex flex-wrap items-center gap-2">
        <div className="min-w-[260px] flex-1">
          <Input
            placeholder="Поиск в таблице"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
      </div>

      <div className="grid gap-4 lg:grid-cols-[minmax(0,1.6fr)_minmax(320px,1fr)]">
        <Card className="overflow-hidden p-0">
          {listQ.isLoading ? (
            <div className="p-4">
              <Skeleton className="h-8 w-full" />
              <Skeleton className="mt-2 h-8 w-full" />
              <Skeleton className="mt-2 h-8 w-full" />
            </div>
          ) : rows.length === 0 ? (
            <div className="p-6">
              <EmptyState title="Нет записей" description="Создайте первую в правой панели." />
            </div>
          ) : (
            <div className="overflow-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="bg-brand-600/10 text-xs uppercase tracking-[0.06em] text-ink-muted">
                    {resource.columns.map((c) => (
                      <th
                        key={String(c.key)}
                        className="whitespace-nowrap border-b border-line px-3 py-2.5 text-left"
                        style={c.width ? { width: c.width } : undefined}
                      >
                        {c.header}
                      </th>
                    ))}
                    <th className="border-b border-line px-3 py-2.5 text-right">Действия</th>
                  </tr>
                </thead>
                <tbody>
                  {rows.map((row) => {
                    const key = String(row[resource.keyField as keyof T] ?? '');
                    const isDeleted =
                      resource.softDelete && Boolean(row['deleted_at']);
                    return (
                      <tr
                        key={key}
                        className={
                          'border-b border-line/60 transition-colors hover:bg-brand-600/5 ' +
                          (isDeleted ? 'opacity-60' : '')
                        }
                      >
                        {resource.columns.map((c) => (
                          <td key={String(c.key)} className="px-3 py-2 align-top">
                            {c.render
                              ? c.render(row)
                              : renderCell((row as Record<string, unknown>)[String(c.key)])}
                          </td>
                        ))}
                        <td className="px-3 py-2 text-right">
                          <div className="inline-flex gap-1">
                            <Button
                              size="sm"
                              variant="secondary"
                              onClick={() => {
                                setCreating(false);
                                setEditing(row);
                              }}
                            >
                              Ред.
                            </Button>
                            <Button
                              size="sm"
                              variant="danger"
                              onClick={() => {
                                if (window.confirm(`Удалить ${key}?`)) removeM.mutate(key);
                              }}
                            >
                              Удалить
                            </Button>
                            {isDeleted && resource.restorePath && (
                              <Button
                                size="sm"
                                variant="warn"
                                onClick={() => restoreM.mutate(key)}
                              >
                                Восстановить
                              </Button>
                            )}
                          </div>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          )}
        </Card>

        <Card>
          <h2 className="mb-3 text-base font-semibold">
            {creating ? 'Создание' : editing ? `Редактирование: ${String(editing[resource.keyField as keyof T])}` : 'Форма'}
          </h2>
          {creating ? (
            <ResourceForm
              resource={resource}
              onSubmit={(v) => createM.mutateAsync(v)}
              onCancel={() => setCreating(false)}
              submitting={createM.isPending}
            />
          ) : editing ? (
            <ResourceForm
              resource={resource}
              initial={editing}
              onSubmit={(v) =>
                updateM.mutateAsync({
                  key: String(editing[resource.keyField as keyof T]),
                  values: v,
                })
              }
              onCancel={() => setEditing(null)}
              submitting={updateM.isPending}
            />
          ) : (
            <EmptyState
              title="Выберите запись"
              description="Нажмите «Ред.» в таблице или «+ Создать» сверху."
            />
          )}
        </Card>
      </div>
    </>
  );
}

function renderCell(value: unknown): React.ReactNode {
  if (value == null) return <span className="text-ink-muted-2">—</span>;
  if (typeof value === 'boolean')
    return value ? (
      <Pill variant="success">да</Pill>
    ) : (
      <Pill>нет</Pill>
    );
  const s = String(value);
  return s.length > 60 ? s.slice(0, 60) + '…' : s;
}

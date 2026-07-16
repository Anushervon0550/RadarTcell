import { useMemo, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Button, Card, EmptyState, Input, Modal, PageHeader, Pill, Skeleton } from '@radartcell/ui';
import { api } from '@/api/client';
import { useToast } from '@/components/Toast';
import { ResourceForm } from './ResourceForm';
import { unwrapList, type ColumnDef, type CrudResource } from './types';

interface Props<T extends Record<string, unknown>> {
  resource: CrudResource<T>;
}

type SortDir = 'asc' | 'desc';

export function ResourceListPage<T extends Record<string, unknown>>({ resource }: Props<T>) {
  const qc = useQueryClient();
  const toast = useToast();
  const [editing, setEditing] = useState<T | null>(null);
  const [creating, setCreating] = useState<boolean>(false);
  const [search, setSearch] = useState('');
  const [sortKey, setSortKey] = useState<string | null>(null);
  const [sortDir, setSortDir] = useState<SortDir>('asc');
  const [showDeleted, setShowDeleted] = useState(true);

  const listKey = ['admin', resource.name, resource.listPath] as const;

  const listQ = useQuery({
    queryKey: listKey,
    queryFn: async () => {
      const data = await api.get<unknown>(resource.listPath);
      return resource.parseList ? resource.parseList(data) : unwrapList<T>(data);
    },
  });

  const closeForm = () => {
    setCreating(false);
    setEditing(null);
  };

  const createM = useMutation({
    mutationFn: (values: Record<string, unknown>) => api.post(resource.createPath, values),
    onSuccess: () => {
      toast.success('Создано');
      closeForm();
      qc.invalidateQueries({ queryKey: ['admin', resource.name] });
    },
    onError: (e) => toast.error(`Ошибка создания: ${(e as Error).message}`),
  });

  const updateM = useMutation({
    mutationFn: ({ key, values }: { key: string; values: Record<string, unknown> }) =>
      api.put(resource.updatePath(key), values),
    onSuccess: () => {
      toast.success('Сохранено');
      closeForm();
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

  const toggleSort = (key: string) => {
    if (sortKey === key) {
      setSortDir((d) => (d === 'asc' ? 'desc' : 'asc'));
    } else {
      setSortKey(key);
      setSortDir('asc');
    }
  };

  const rows = useMemo(() => {
    let data = listQ.data ?? [];

    if (resource.softDelete && !showDeleted) {
      data = data.filter((r) => !r['deleted_at']);
    }

    if (search.trim()) {
      const q = search.trim().toLowerCase();
      data = data.filter((r) =>
        Object.values(r).some((v) => String(v ?? '').toLowerCase().includes(q)),
      );
    }

    if (sortKey) {
      data = [...data].sort((a, b) => {
        const av = a[sortKey];
        const bv = b[sortKey];
        const an = typeof av === 'number' ? av : Number(av);
        const bn = typeof bv === 'number' ? bv : Number(bv);
        let cmp: number;
        if (!Number.isNaN(an) && !Number.isNaN(bn) && av !== '' && bv !== '') {
          cmp = an - bn;
        } else {
          cmp = String(av ?? '').localeCompare(String(bv ?? ''), 'ru');
        }
        return sortDir === 'asc' ? cmp : -cmp;
      });
    }

    return data;
  }, [listQ.data, search, sortKey, sortDir, showDeleted, resource.softDelete]);

  const total = listQ.data?.length ?? 0;
  const deletedCount = resource.softDelete
    ? (listQ.data ?? []).filter((r) => r['deleted_at']).length
    : 0;

  const formOpen = creating || editing !== null;

  return (
    <>
      <PageHeader
        title={resource.title}
        description={
          listQ.isLoading
            ? 'Загружаем…'
            : `Показано ${rows.length} из ${total}${deletedCount ? ` · удалённых: ${deletedCount}` : ''}`
        }
        actions={
          <Button
            variant="primary"
            icon={<span aria-hidden>+</span>}
            onClick={() => {
              setEditing(null);
              setCreating(true);
            }}
          >
            Создать
          </Button>
        }
      />

      <div className="mb-4 flex flex-wrap items-center gap-3">
        <div className="min-w-[260px] flex-1">
          <Input
            placeholder={`Поиск по ${resource.title.toLowerCase()}…`}
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
        {search && (
          <Button variant="ghost" size="sm" onClick={() => setSearch('')}>
            Сбросить
          </Button>
        )}
        {resource.softDelete && (
          <label className="flex select-none items-center gap-2 rounded-lg border border-line bg-bg-panel px-3 py-2 text-sm text-ink-muted">
            <input
              type="checkbox"
              className="h-4 w-4 rounded border-line bg-bg-soft"
              checked={showDeleted}
              onChange={(e) => setShowDeleted(e.target.checked)}
            />
            Показывать удалённые
          </label>
        )}
        <Button
          variant="secondary"
          size="sm"
          onClick={() => listQ.refetch()}
          loading={listQ.isFetching}
        >
          Обновить
        </Button>
      </div>

      <Card className="overflow-hidden p-0">
        {listQ.isLoading ? (
          <div className="p-4">
            <Skeleton className="h-8 w-full" />
            <Skeleton className="mt-2 h-8 w-full" />
            <Skeleton className="mt-2 h-8 w-full" />
            <Skeleton className="mt-2 h-8 w-full" />
          </div>
        ) : rows.length === 0 ? (
          <div className="p-6">
            <EmptyState
              title={search ? 'Ничего не найдено' : 'Нет записей'}
              description={
                search
                  ? 'Попробуйте изменить поисковый запрос.'
                  : 'Нажмите «Создать», чтобы добавить первую запись.'
              }
            />
          </div>
        ) : (
          <div className="max-h-[calc(100vh-260px)] overflow-auto">
            <table className="w-full text-sm">
              <thead className="sticky top-0 z-10">
                <tr className="bg-[#141d33] text-xs uppercase tracking-[0.06em] text-ink-muted">
                  {resource.columns.map((c) => (
                    <SortableHeader
                      key={String(c.key)}
                      column={c}
                      active={sortKey === String(c.key)}
                      dir={sortDir}
                      onSort={() => toggleSort(String(c.key))}
                    />
                  ))}
                  <th className="border-b border-line px-3 py-2.5 text-right">Действия</th>
                </tr>
              </thead>
              <tbody>
                {rows.map((row) => {
                  const key = String(row[resource.keyField as keyof T] ?? '');
                  const isDeleted = resource.softDelete && Boolean(row['deleted_at']);
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
                          {!isDeleted && (
                            <Button
                              size="sm"
                              variant="danger"
                              onClick={() => {
                                if (window.confirm(`Удалить «${key}»?`)) removeM.mutate(key);
                              }}
                            >
                              Удалить
                            </Button>
                          )}
                          {isDeleted && resource.restorePath && (
                            <Button size="sm" variant="warn" onClick={() => restoreM.mutate(key)}>
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

      <Modal open={formOpen} onClose={closeForm} size="lg" labelledBy="resource-form-title">
        <div className="flex items-center justify-between border-b border-line px-6 py-4">
          <h2 id="resource-form-title" className="text-lg font-semibold">
            {creating
              ? `Создание — ${resource.title}`
              : `Редактирование: ${editing ? String(editing[resource.keyField as keyof T]) : ''}`}
          </h2>
          <button
            type="button"
            onClick={closeForm}
            aria-label="Закрыть"
            className="h-8 w-8 rounded-lg border border-line text-ink-muted hover:bg-white/5 hover:text-white"
          >
            ×
          </button>
        </div>
        <div className="max-h-[72vh] overflow-y-auto px-6 py-5">
          {creating ? (
            <ResourceForm
              resource={resource}
              onSubmit={async (v) => {
                await createM.mutateAsync(v);
              }}
              onCancel={closeForm}
              submitting={createM.isPending}
            />
          ) : editing ? (
            <ResourceForm
              resource={resource}
              initial={editing}
              onSubmit={async (v) => {
                await updateM.mutateAsync({
                  key: String(editing[resource.keyField as keyof T]),
                  values: v,
                });
              }}
              onCancel={closeForm}
              submitting={updateM.isPending}
            />
          ) : null}
        </div>
      </Modal>
    </>
  );
}

function SortableHeader<T>({
  column,
  active,
  dir,
  onSort,
}: {
  column: ColumnDef<T>;
  active: boolean;
  dir: SortDir;
  onSort: () => void;
}) {
  return (
    <th
      className="whitespace-nowrap border-b border-line px-3 py-2.5 text-left"
      style={column.width ? { width: column.width } : undefined}
    >
      <button
        type="button"
        onClick={onSort}
        className="inline-flex items-center gap-1 uppercase tracking-[0.06em] transition-colors hover:text-white"
      >
        {column.header}
        <span aria-hidden className={active ? 'text-brand-400' : 'text-ink-muted-2'}>
          {active ? (dir === 'asc' ? '▲' : '▼') : '⇅'}
        </span>
      </button>
    </th>
  );
}

function renderCell(value: unknown): React.ReactNode {
  if (value == null) return <span className="text-ink-muted-2">—</span>;
  if (typeof value === 'boolean')
    return value ? <Pill variant="success">да</Pill> : <Pill>нет</Pill>;
  const s = String(value);
  return s.length > 60 ? s.slice(0, 60) + '…' : s;
}

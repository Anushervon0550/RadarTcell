import { useState, type FormEvent } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { AdminUser } from '@radartcell/api';
import {
  Button,
  Card,
  EmptyState,
  Input,
  PageHeader,
  Pill,
  Skeleton,
} from '@radartcell/ui';
import { api } from '@/api/client';
import { useToast } from '@/components/Toast';
import { unwrapList } from '@/crud/types';

export function UsersPage() {
  const qc = useQueryClient();
  const toast = useToast();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const listQ = useQuery({
    queryKey: ['admin', 'users'] as const,
    queryFn: async () => unwrapList<AdminUser>(await api.get('/api/admin/users')),
  });

  const createM = useMutation({
    mutationFn: (body: { username: string; password: string }) =>
      api.post('/api/admin/users', body),
    onSuccess: () => {
      toast.success('Пользователь создан');
      setUsername('');
      setPassword('');
      qc.invalidateQueries({ queryKey: ['admin', 'users'] });
    },
    onError: (e) => toast.error(`Ошибка: ${(e as Error).message}`),
  });

  const toggleM = useMutation({
    mutationFn: ({ user, active }: { user: string; active: boolean }) =>
      api.put(
        `/api/admin/users/${encodeURIComponent(user)}/${active ? 'activate' : 'deactivate'}`,
      ),
    onSuccess: () => {
      toast.success('Статус обновлён');
      qc.invalidateQueries({ queryKey: ['admin', 'users'] });
    },
    onError: (e) => toast.error(`Ошибка: ${(e as Error).message}`),
  });

  const onCreate = (e: FormEvent) => {
    e.preventDefault();
    if (!username.trim() || !password.trim()) return;
    createM.mutate({ username: username.trim(), password });
  };

  return (
    <>
      <PageHeader
        title="Пользователи"
        description="Управление админ-аккаунтами (bcrypt-хеши хранятся в БД)."
      />

      <div className="grid gap-4 lg:grid-cols-[minmax(0,1.6fr)_minmax(320px,1fr)]">
        <Card className="p-0">
          {listQ.isLoading ? (
            <div className="p-4">
              <Skeleton className="h-8 w-full" />
              <Skeleton className="mt-2 h-8 w-full" />
            </div>
          ) : (listQ.data ?? []).length === 0 ? (
            <div className="p-6">
              <EmptyState title="Нет пользователей" />
            </div>
          ) : (
            <table className="w-full text-sm">
              <thead>
                <tr className="bg-brand-600/10 text-xs uppercase tracking-[0.06em] text-ink-muted">
                  <th className="border-b border-line px-3 py-2.5 text-left">Логин</th>
                  <th className="border-b border-line px-3 py-2.5 text-left">Статус</th>
                  <th className="border-b border-line px-3 py-2.5 text-right">Действия</th>
                </tr>
              </thead>
              <tbody>
                {(listQ.data ?? []).map((u) => (
                  <tr key={u.username} className="border-b border-line/60 hover:bg-brand-600/5">
                    <td className="px-3 py-2 font-medium">{u.username}</td>
                    <td className="px-3 py-2">
                      {u.is_active ? (
                        <Pill variant="success">активен</Pill>
                      ) : (
                        <Pill variant="danger">заблокирован</Pill>
                      )}
                    </td>
                    <td className="px-3 py-2 text-right">
                      {u.is_active ? (
                        <Button
                          size="sm"
                          variant="warn"
                          onClick={() =>
                            toggleM.mutate({ user: u.username, active: false })
                          }
                        >
                          Заблокировать
                        </Button>
                      ) : (
                        <Button
                          size="sm"
                          variant="success"
                          onClick={() =>
                            toggleM.mutate({ user: u.username, active: true })
                          }
                        >
                          Активировать
                        </Button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </Card>

        <Card>
          <h2 className="mb-3 text-base font-semibold">Создать пользователя</h2>
          <form onSubmit={onCreate} className="space-y-3">
            <Input
              label="Логин"
              autoComplete="off"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
            <Input
              label="Пароль"
              type="password"
              autoComplete="new-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              hint="Минимум 12 символов, будет захеширован bcrypt-ом на сервере."
            />
            <Button variant="primary" type="submit" loading={createM.isPending} fullWidth>
              Создать
            </Button>
          </form>
        </Card>
      </div>
    </>
  );
}

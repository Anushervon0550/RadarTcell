import { useState, type FormEvent } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Button, Card, Input } from '@radartcell/ui';
import { api } from '@/api/client';
import { useAuthStore } from '@/auth/store';
import { useToast } from '@/components/Toast';
import type { AdminLoginResponse } from '@radartcell/api';

export function LoginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const setToken = useAuthStore((s) => s.setToken);
  const navigate = useNavigate();
  const [params] = useSearchParams();
  const toast = useToast();

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      const res = await api.post<AdminLoginResponse>('/api/admin/login', { username, password });
      if (!res?.token) throw new Error('empty token');
      setToken(res.token, username);
      toast.success('Вход выполнен');
      const next = params.get('next') || '/technologies';
      navigate(next, { replace: true });
    } catch (err) {
      toast.error(`Не удалось войти: ${(err as Error).message}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center px-4">
      <Card className="w-full max-w-md">
        <div className="mb-6 text-center">
          <div className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-gradient-to-br from-brand-600 to-accent text-xl font-extrabold shadow-glow">
            R
          </div>
          <h1 className="text-xl font-bold">RadarTcell Admin</h1>
          <p className="mt-1 text-sm text-ink-muted">Вход для администраторов</p>
        </div>
        <form onSubmit={onSubmit} className="space-y-3">
          <Input
            label="Логин"
            autoComplete="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
          <Input
            label="Пароль"
            type="password"
            autoComplete="current-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          <Button variant="primary" type="submit" fullWidth loading={loading}>
            Войти
          </Button>
        </form>
      </Card>
    </div>
  );
}

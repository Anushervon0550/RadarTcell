import { NavLink, Outlet, useNavigate } from 'react-router-dom';
import { cn } from '@radartcell/ui';
import { env } from '@/env';
import { useAuthStore } from '@/auth/store';

const groups = [
  {
    title: 'Обзор',
    links: [{ to: '/', label: 'Дашборд', icon: '◈', end: true }],
  },
  {
    title: 'Контент',
    links: [
      { to: '/technologies', label: 'Технологии', icon: '⚙' },
      { to: '/trends', label: 'Тренды', icon: '↗' },
      { to: '/tags', label: 'Теги', icon: '#' },
      { to: '/sdgs', label: 'ЦУР', icon: '◯' },
      { to: '/organizations', label: 'Организации', icon: '▣' },
      { to: '/metrics', label: 'Метрики', icon: '⋯' },
    ],
  },
  {
    title: 'Доступ',
    links: [{ to: '/users', label: 'Пользователи', icon: '⚿' }],
  },
];

export function AdminLayout() {
  const navigate = useNavigate();
  const username = useAuthStore((s) => s.username);
  const clear = useAuthStore((s) => s.clear);

  const logout = () => {
    clear();
    navigate('/login', { replace: true });
  };

  return (
    <div className="grid min-h-screen grid-cols-[240px_1fr] max-[900px]:grid-cols-1">
      <aside className="sticky top-0 h-screen overflow-y-auto border-r border-line bg-black/50 p-4 backdrop-blur max-[900px]:static max-[900px]:h-auto max-[900px]:border-b max-[900px]:border-r-0">
        <div className="mb-6 flex items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-brand-600 to-accent text-lg font-extrabold shadow-glow">
            R
          </div>
          <div className="min-w-0">
            <div className="text-sm font-bold tracking-tight">RadarTcell</div>
            <div className="text-[11px] uppercase tracking-[0.12em] text-brand-400">
              Admin
            </div>
          </div>
        </div>

        {groups.map((g) => (
          <div key={g.title}>
            <div className="mt-3 text-[11px] font-bold uppercase tracking-[0.12em] text-ink-muted-2">
              {g.title}
            </div>
            <nav className="mt-1 flex flex-col gap-1">
              {g.links.map((l) => (
                <NavLink
                  key={l.to}
                  to={l.to}
                  end={'end' in l ? (l as { end?: boolean }).end : undefined}
                  className={({ isActive }) =>
                    cn(
                      'group flex items-center gap-2.5 rounded-lg border border-transparent px-3 py-2 text-sm transition-all',
                      'hover:border-brand-600/45 hover:bg-brand-600/10',
                      isActive &&
                        'border-brand-600/45 bg-gradient-to-b from-brand-600/15 to-accent/10 text-white',
                    )
                  }
                >
                  <span aria-hidden className="text-[13px] text-ink-muted group-hover:text-brand-400">
                    {l.icon}
                  </span>
                  <span>{l.label}</span>
                </NavLink>
              ))}
            </nav>
          </div>
        ))}

        <div className="mt-6 border-t border-line pt-3 text-xs text-ink-muted">
          <div className="mb-2 truncate">Вошли как <b className="text-ink">{username ?? '—'}</b></div>
          <a
            href={env.PUBLIC_URL}
            target="_blank"
            rel="noopener noreferrer"
            className="block rounded-lg px-3 py-2 text-ink-muted hover:text-white"
          >
            Публичный сайт ↗
          </a>
          <button
            onClick={logout}
            type="button"
            className="mt-1 w-full rounded-lg px-3 py-2 text-left text-red-300 hover:bg-red-500/10"
          >
            Выйти
          </button>
        </div>
      </aside>

      <main className="min-w-0 px-8 py-7 max-[900px]:px-4 max-[900px]:py-4">
        <Outlet />
      </main>
    </div>
  );
}

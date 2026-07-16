import { NavLink, Outlet } from 'react-router-dom';
import { cn } from '@radartcell/ui';
import { env } from '@/env';

const publicLinks = [
  { to: '/', label: 'Радар', icon: '◎', end: true },
  { to: '/catalog', label: 'Каталог', icon: '▤' },
];

const entityLinks = [
  { to: '/trends', label: 'Тренды', icon: '↗' },
  { to: '/tags', label: 'Теги', icon: '#' },
  { to: '/sdgs', label: 'ЦУР', icon: '◯' },
  { to: '/organizations', label: 'Организации', icon: '▣' },
];

function MenuLink({
  to,
  label,
  icon,
  end,
}: {
  to: string;
  label: string;
  icon: string;
  end?: boolean;
}) {
  return (
    <NavLink
      to={to}
      end={end}
      className={({ isActive }) =>
        cn(
          'group flex items-center gap-2.5 rounded-lg border border-transparent px-3 py-2 text-sm transition-all',
          'hover:border-brand-600/45 hover:bg-gradient-to-b hover:from-brand-600/15 hover:to-accent/10',
          isActive &&
            'border-brand-600/45 bg-gradient-to-b from-brand-600/15 to-accent/10 text-white',
        )
      }
    >
      <span
        aria-hidden
        className="flex h-4 w-4 items-center justify-center text-[13px] text-ink-muted group-hover:text-brand-400"
      >
        {icon}
      </span>
      <span className="truncate">{label}</span>
    </NavLink>
  );
}

export function Layout() {
  return (
    <div className="grid min-h-screen grid-cols-[260px_1fr] max-[1080px]:grid-cols-1">
      <aside className="sticky top-0 h-screen overflow-y-auto border-r border-line bg-black/40 p-5 backdrop-blur max-[1080px]:static max-[1080px]:h-auto max-[1080px]:border-b max-[1080px]:border-r-0">
        <div className="mb-6 flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-brand-600 to-accent text-lg font-extrabold shadow-glow">
            R
          </div>
          <div className="min-w-0">
            <div className="text-sm font-bold tracking-tight">RadarTcell</div>
            <div className="text-[11px] text-ink-muted-2">Radar of Technologies</div>
          </div>
        </div>

        <div className="mt-4 text-[11px] font-bold uppercase tracking-[0.12em] text-ink-muted-2">
          Обзор
        </div>
        <nav className="mt-2 flex flex-col gap-1">
          {publicLinks.map((l) => (
            <MenuLink key={l.to} {...l} />
          ))}
        </nav>

        <div className="mt-4 text-[11px] font-bold uppercase tracking-[0.12em] text-ink-muted-2">
          Сущности
        </div>
        <nav className="mt-2 flex flex-col gap-1">
          {entityLinks.map((l) => (
            <MenuLink key={l.to} {...l} />
          ))}
        </nav>

        {env.ENABLE_SWAGGER && (
          <>
            <div className="mt-4 text-[11px] font-bold uppercase tracking-[0.12em] text-ink-muted-2">
              Системное
            </div>
            <nav className="mt-2 flex flex-col gap-1">
              <a
                href="/swagger/"
                target="_blank"
                rel="noopener noreferrer"
                className="rounded-lg px-3 py-2 text-sm text-ink-muted hover:text-white"
              >
                Swagger UI
              </a>
              <a
                href="/openapi.yaml"
                target="_blank"
                rel="noopener noreferrer"
                className="rounded-lg px-3 py-2 text-sm text-ink-muted hover:text-white"
              >
                OpenAPI
              </a>
              <a
                href="/healthz"
                target="_blank"
                rel="noopener noreferrer"
                className="rounded-lg px-3 py-2 text-sm text-ink-muted hover:text-white"
              >
                Health
              </a>
            </nav>
          </>
        )}
      </aside>

      <main className="min-w-0 px-8 py-7 max-[1080px]:px-4 max-[1080px]:py-4">
        <Outlet />
      </main>
    </div>
  );
}

import type { HomeTrendBlock } from '@radartcell/api';
import { cn } from '@radartcell/ui';
import { PALETTE } from './constants';

interface Props {
  trends: HomeTrendBlock[];
  activeSlug: string | null;
  onHover: (slug: string | null) => void;
  onSelect: (slug: string) => void;
  totalTechs: number;
  productCount: number;
}

export function TrendsLegend({
  trends,
  activeSlug,
  onHover,
  onSelect,
  totalTechs,
  productCount,
}: Props) {
  return (
    <aside className="pointer-events-auto absolute left-3 top-3 z-[5] max-h-[calc(100%-24px)] w-72 overflow-y-auto rounded-2xl border border-line bg-black/70 p-4 backdrop-blur-md max-[1100px]:static max-[1100px]:mb-3 max-[1100px]:w-full">
      <h3 className="mb-3 text-[11px] font-bold uppercase tracking-[0.1em] text-ink-muted">
        Тренды
      </h3>
      <div className="flex flex-col gap-2">
        {trends.length === 0 && (
          <div className="rounded-lg border border-dashed border-line px-3 py-4 text-center text-xs text-ink-muted">
            Нет данных
          </div>
        )}
        {trends.map((t, idx) => {
          const color = PALETTE[idx % PALETTE.length]!;
          const isActive = activeSlug === t.slug;
          return (
            <button
              key={t.slug}
              type="button"
              onMouseEnter={() => onHover(t.slug)}
              onMouseLeave={() => onHover(null)}
              onClick={() => onSelect(t.slug)}
              className={cn(
                'flex w-full items-center gap-2.5 rounded-xl border px-3 py-2 text-left text-sm transition-all',
                'hover:border-brand-600/50',
                isActive
                  ? 'border-brand-600/50 bg-gradient-to-b from-brand-600/20 to-accent/10'
                  : 'border-line bg-bg-panel',
              )}
            >
              <span
                className="h-3.5 w-3.5 shrink-0 rounded"
                style={{ background: color, boxShadow: `0 0 8px ${color}` }}
              />
              <span className="flex-1 truncate">{t.name}</span>
              <span className="text-xs text-ink-muted">{t.items.length}</span>
            </button>
          );
        })}
      </div>

      <div className="mt-4 grid grid-cols-2 gap-2">
        <SummaryCell label="Технологий" value={totalTechs} />
        <SummaryCell label="Трендов" value={trends.length} />
        <SummaryCell label="В продакшене" value={productCount} />
      </div>

      <div className="mt-4 border-t border-line pt-3">
        <div className="mb-2 text-[11px] font-bold uppercase tracking-[0.1em] text-ink-muted">
          Кольца
        </div>
        <div className="flex flex-col gap-1 text-xs text-ink-muted">
          <LegendDot color="#7c3aed" label="Production · TRL 8–9" />
          <LegendDot color="#3b82f6" label="Adopting · TRL 6–7" />
          <LegendDot color="#f59e0b" label="Experimenting · TRL 4–5" />
          <LegendDot color="#9ca3af" label="Envisioning · TRL 1–3" />
        </div>
      </div>
    </aside>
  );
}

function SummaryCell({ label, value }: { label: string; value: number }) {
  return (
    <div className="min-w-0 rounded-lg border border-line bg-bg-panel px-2.5 py-2">
      <div className="truncate text-[10px] uppercase tracking-[0.07em] text-ink-muted">{label}</div>
      <div className="mt-0.5 text-lg font-bold">{value}</div>
    </div>
  );
}

function LegendDot({ color, label }: { color: string; label: string }) {
  return (
    <div className="flex items-center gap-2">
      <span className="inline-block h-2 w-2 rounded-full" style={{ background: color }} />
      <span>{label}</span>
    </div>
  );
}

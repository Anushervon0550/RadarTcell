import { Link } from 'react-router-dom';
import { Button, cssUrl, Pill, safeUrl, Skeleton, StageChip, TrlChip, cn } from '@radartcell/ui';
import type { TechnologyCard } from '@radartcell/api';
import { OrgLogo } from '@/components/OrgLogo';
import { FALLBACK_COVER } from './constants';

interface Props {
  slug: string | null;
  data: TechnologyCard | null;
  loading: boolean;
  onClose: () => void;
}

export function TechInfoPanel({ slug, data, loading, onClose }: Props) {
  const collapsed = !slug;

  return (
    <aside
      aria-hidden={collapsed}
      className={cn(
        'pointer-events-auto absolute right-3 top-3 z-[5] flex max-h-[calc(100%-24px)] w-[360px] flex-col overflow-hidden rounded-2xl border border-line bg-[rgba(11,16,30,0.86)] shadow-elevated backdrop-blur-md transition-all duration-200 max-[1280px]:w-[300px] max-[1100px]:static max-[1100px]:w-full',
        collapsed && 'pointer-events-none translate-x-[110%] opacity-0 max-[1100px]:hidden',
      )}
    >
      {loading && !data ? <PanelSkeleton onClose={onClose} /> : null}
      {!loading && data && <PanelContent data={data} onClose={onClose} />}
    </aside>
  );
}

function PanelSkeleton({ onClose }: { onClose: () => void }) {
  return (
    <div className="flex-1 overflow-y-auto">
      <Skeleton className="h-32 w-full" />
      <div className="space-y-3 p-4">
        <Skeleton className="h-6 w-3/4" />
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-5/6" />
      </div>
      <button
        type="button"
        onClick={onClose}
        className="absolute right-2 top-2 h-8 w-8 rounded-lg border border-line bg-black/60 text-white hover:bg-brand-600/40"
        aria-label="Закрыть"
      >
        ×
      </button>
    </div>
  );
}

function MetricRow({ label, value }: { label: string; value?: number | null }) {
  const pct = Math.max(0, Math.min(100, Math.round((Number(value) || 0) * 100)));
  return (
    <div className="mt-2 flex items-center gap-2 text-xs">
      <span className="w-20 text-ink-muted">{label}</span>
      <span className="h-1.5 flex-1 overflow-hidden rounded-full bg-line-soft">
        <span
          className="block h-full bg-gradient-to-r from-brand-600 to-accent"
          style={{ width: `${pct}%` }}
        />
      </span>
      <span className="w-10 text-right text-ink-muted">{pct}%</span>
    </div>
  );
}

function PanelContent({ data, onClose }: { data: TechnologyCard; onClose: () => void }) {
  const cover = data.image_url ?? FALLBACK_COVER;
  const source = safeUrl(data.source_link);

  return (
    <div className="relative flex-1 overflow-y-auto">
      <div
        className="relative h-32 bg-cover bg-center"
        style={{ backgroundImage: cssUrl(cover) }}
      >
        <div className="absolute inset-0 bg-gradient-to-b from-transparent via-transparent to-[rgba(13,20,38,0.95)]" />
        <button
          type="button"
          onClick={onClose}
          aria-label="Закрыть"
          className="absolute right-2 top-2 z-[2] h-8 w-8 rounded-lg border border-line bg-black/60 text-white hover:bg-brand-600/40"
        >
          ×
        </button>
      </div>

      <div className="space-y-4 px-4 py-4">
        <div>
          <h2 className="text-base font-semibold leading-tight">{data.name}</h2>
          <div className="mt-2 flex flex-wrap items-center gap-1.5">
            <TrlChip trl={data.trl} />
            <StageChip trl={data.trl} />
            <Link to={`/trend/${encodeURIComponent(data.trend_slug)}`}>
              <Pill>{data.trend_name || data.trend_slug}</Pill>
            </Link>
          </div>
        </div>

        {data.description_short && (
          <section>
            <h4 className="mb-1 text-[11px] uppercase tracking-[0.08em] text-ink-muted">Кратко</h4>
            <p className="text-sm leading-relaxed text-[#d6dff0]">{data.description_short}</p>
          </section>
        )}

        {data.description_full && (
          <section>
            <h4 className="mb-1 text-[11px] uppercase tracking-[0.08em] text-ink-muted">Подробно</h4>
            <p className="text-sm leading-relaxed text-[#d6dff0]">{data.description_full}</p>
          </section>
        )}

        <section>
          <h4 className="mb-1 text-[11px] uppercase tracking-[0.08em] text-ink-muted">Метрики</h4>
          <MetricRow label="Зрелость" value={data.custom_metric_1} />
          <MetricRow label="Влияние" value={data.custom_metric_2} />
          <MetricRow label="Покрытие" value={data.custom_metric_3} />
          <MetricRow label="Стоимость" value={data.custom_metric_4} />
        </section>

        {data.tags?.length ? (
          <section>
            <h4 className="mb-1 text-[11px] uppercase tracking-[0.08em] text-ink-muted">Теги</h4>
            <div className="flex flex-wrap gap-1.5">
              {data.tags.map((tag) => (
                <Link key={tag.id} to={`/tag/${encodeURIComponent(tag.slug)}`}>
                  <Pill variant="accent">{tag.title || tag.slug}</Pill>
                </Link>
              ))}
            </div>
          </section>
        ) : null}

        {data.sdgs?.length ? (
          <section>
            <h4 className="mb-1 text-[11px] uppercase tracking-[0.08em] text-ink-muted">ЦУР</h4>
            <div className="flex flex-wrap gap-1.5">
              {data.sdgs.map((s) => (
                <Link key={s.id} to={`/sdg/${encodeURIComponent(s.code)}`}>
                  <Pill>
                    {s.code} · {s.title}
                  </Pill>
                </Link>
              ))}
            </div>
          </section>
        ) : null}

        {data.organizations?.length ? (
          <section>
            <h4 className="mb-1 text-[11px] uppercase tracking-[0.08em] text-ink-muted">Организации</h4>
            <div className="flex flex-col gap-2">
              {data.organizations.map((org) => (
                <Link
                  key={org.id}
                  to={`/organization/${encodeURIComponent(org.slug)}`}
                  className="flex items-center gap-2.5 text-sm text-ink hover:text-white"
                >
                  <OrgLogo
                    name={org.name}
                    logoUrl={org.logo_url}
                    className="h-9 w-9 rounded-lg text-xs"
                  />
                  <span>
                    <div className="font-medium">{org.name}</div>
                    <div className="text-xs text-ink-muted">
                      {org.headquarters || org.website || ''}
                    </div>
                  </span>
                </Link>
              ))}
            </div>
          </section>
        ) : null}

        <div className="flex flex-wrap gap-2 pt-2">
          <Link to={`/technology/${encodeURIComponent(data.slug)}`}>
            <Button variant="primary" size="sm">
              Открыть полностью
            </Button>
          </Link>
          {source && (
            <a href={source} target="_blank" rel="noopener noreferrer">
              <Button variant="ghost" size="sm">
                Источник ↗
              </Button>
            </a>
          )}
        </div>
      </div>
    </div>
  );
}

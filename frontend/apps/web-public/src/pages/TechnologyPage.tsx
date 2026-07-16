import { Link, useNavigate, useParams } from 'react-router-dom';
import {
  Button,
  cssUrl,
  EmptyState,
  Modal,
  Pill,
  Skeleton,
  StageChip,
  TrlChip,
  safeUrl,
  stageFromTrl,
  stageLabel,
} from '@radartcell/ui';
import { useTechnology } from '@/api/queries';
import { OrgLogo } from '@/components/OrgLogo';
import { FALLBACK_COVER } from '@/pages/radar/constants';

export function TechnologyPage() {
  const { slug } = useParams<{ slug: string }>();
  const navigate = useNavigate();
  const { data, isLoading, error } = useTechnology(slug);

  const close = () => navigate(-1);

  return (
    <Modal open onClose={close} size="lg">
      {isLoading && (
        <div className="p-6">
          <Skeleton className="h-56 w-full rounded-t-2xl" />
          <div className="space-y-3 p-4">
            <Skeleton className="h-8 w-2/3" />
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-5/6" />
          </div>
        </div>
      )}

      {error && (
        <div className="p-6">
          <EmptyState title="Ошибка" description={(error as Error).message} />
        </div>
      )}

      {data && <TechDetail data={data} onClose={close} />}
    </Modal>
  );
}

function TechDetail({
  data,
  onClose,
}: {
  data: NonNullable<ReturnType<typeof useTechnology>['data']>;
  onClose: () => void;
}) {
  const cover = data.image_url ?? FALLBACK_COVER;
  const stage = stageFromTrl(data.trl);
  const completion = Math.max(
    0,
    Math.min(100, Math.round(((data.custom_metric_1 as number) || 0) * 100)),
  );
  const source = safeUrl(data.source_link);

  return (
    <>
      <button
        type="button"
        onClick={onClose}
        aria-label="Закрыть"
        className="absolute right-3 top-3 z-10 h-9 w-9 rounded-lg border border-line bg-black/70 text-white hover:bg-brand-600/40"
      >
        ×
      </button>
      <div
        className="relative h-60 bg-cover bg-center"
        style={{ backgroundImage: cssUrl(cover) }}
      >
        <div className="absolute inset-0 bg-gradient-to-b from-transparent via-transparent to-[rgba(13,20,38,0.95)]" />
      </div>

      <div className="px-6 pb-8 pt-4 md:px-8">
        <div className="-mt-14 space-y-3">
          <h2 className="text-3xl font-bold leading-tight">{data.name}</h2>
          <div className="flex flex-wrap items-center gap-1.5">
            <TrlChip trl={data.trl} />
            <StageChip trl={data.trl} />
            <Link to={`/trend/${encodeURIComponent(data.trend_slug)}`}>
              <Pill variant="accent">{data.trend_name || data.trend_slug}</Pill>
            </Link>
            {source && (
              <a href={source} target="_blank" rel="noopener noreferrer">
                <Pill>Источник ↗</Pill>
              </a>
            )}
          </div>
        </div>

        {data.description_short && (
          <Section title="Краткое описание">
            <p className="leading-relaxed text-[#d6dff0]">{data.description_short}</p>
          </Section>
        )}
        {data.description_full && (
          <Section title="Подробно">
            <p className="whitespace-pre-line leading-relaxed text-[#d6dff0]">
              {data.description_full}
            </p>
          </Section>
        )}

        <Section title="Метрики">
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
            <MetricCell label="Зрелость" value={data.custom_metric_1} />
            <MetricCell label="Влияние" value={data.custom_metric_2} />
            <MetricCell label="Покрытие" value={data.custom_metric_3} />
            <MetricCell label="Стоимость" value={data.custom_metric_4} />
          </div>
        </Section>

        <Section title="Стадия развития">
          <p className="mb-3">
            Готовность: <strong>{completion}%</strong> · Уровень зрелости:{' '}
            <strong>TRL {data.trl}</strong> · Класс: <strong>{stageLabel(stage)}</strong>
          </p>
          <div className="h-2 overflow-hidden rounded-full bg-line-soft">
            <span
              className="block h-full bg-gradient-to-r from-brand-600 to-accent"
              style={{ width: `${Math.max(10, Math.min(100, ((data.trl || 0) / 9) * 100))}%` }}
            />
          </div>
        </Section>

        {data.tags?.length ? (
          <Section title="Теги">
            <div className="flex flex-wrap gap-1.5">
              {data.tags.map((t) => (
                <Link key={t.id} to={`/tag/${encodeURIComponent(t.slug)}`}>
                  <Pill variant="accent">{t.title || t.slug}</Pill>
                </Link>
              ))}
            </div>
          </Section>
        ) : null}

        {data.sdgs?.length ? (
          <Section title="Цели устойчивого развития">
            <div className="flex flex-wrap gap-1.5">
              {data.sdgs.map((s) => (
                <Link key={s.id} to={`/sdg/${encodeURIComponent(s.code)}`}>
                  <Pill>
                    {s.code} · {s.title}
                  </Pill>
                </Link>
              ))}
            </div>
          </Section>
        ) : null}

        {data.organizations?.length ? (
          <Section title="Организации">
            <div className="grid gap-3 sm:grid-cols-2">
              {data.organizations.map((o) => (
                <Link
                  key={o.id}
                  to={`/organization/${encodeURIComponent(o.slug)}`}
                  className="flex items-center gap-3 rounded-xl border border-line bg-bg-panel/50 p-3 transition-all hover:border-brand-600/50"
                >
                  <OrgLogo
                    name={o.name}
                    logoUrl={o.logo_url}
                    className="h-10 w-10 rounded-md text-sm"
                  />
                  <span>
                    <div className="font-semibold">{o.name}</div>
                    <div className="text-xs text-ink-muted">
                      {o.headquarters || o.website || ''}
                    </div>
                  </span>
                </Link>
              ))}
            </div>
          </Section>
        ) : null}

        <div className="mt-6 flex justify-end">
          <Button onClick={onClose}>Закрыть</Button>
        </div>
      </div>
    </>
  );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <section className="mt-6">
      <h3 className="mb-2 text-[13px] uppercase tracking-[0.08em] text-ink-muted">{title}</h3>
      {children}
    </section>
  );
}

function MetricCell({ label, value }: { label: string; value?: number | null }) {
  const pct = Math.max(0, Math.min(100, Math.round((Number(value) || 0) * 100)));
  return (
    <div className="rounded-xl border border-line bg-bg-panel p-3">
      <div className="text-xs text-ink-muted">{label}</div>
      <div className="mt-1 text-lg font-bold">{pct}%</div>
      <div className="mt-2 h-2 overflow-hidden rounded-full bg-line-soft">
        <span
          className="block h-full bg-gradient-to-r from-brand-600 to-accent"
          style={{ width: `${pct}%` }}
        />
      </div>
    </div>
  );
}

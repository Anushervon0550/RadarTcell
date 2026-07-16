import { useMemo, useRef, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { Button, PageHeader, Spinner } from '@radartcell/ui';
import { useHome, useTechnology } from '@/api/queries';
import { RadarChart, type RadarChartHandle } from './RadarChart';
import { RadarControls } from './RadarControls';
import { TrendsLegend } from './TrendsLegend';
import { TechInfoPanel } from './TechInfoPanel';
import { PALETTE } from './constants';
import { env } from '@/env';

export function RadarPage() {
  const params = useParams<{ trend?: string }>();
  const navigate = useNavigate();
  const [selectedSlug, setSelectedSlug] = useState<string | null>(null);
  const [hoverTrend, setHoverTrend] = useState<string | null>(null);
  const radarRef = useRef<RadarChartHandle>(null);

  const { data: home, isLoading, error } = useHome(200);
  const { data: tech, isFetching: techLoading } = useTechnology(selectedSlug ?? undefined);

  const allTrends = useMemo(
    () => (home?.trends ?? []).filter((t) => t.items && t.items.length > 0),
    [home],
  );

  const focused = params.trend ? allTrends.find((t) => t.slug === params.trend) ?? null : null;
  const trends = focused ? [focused] : allTrends;

  // Toggle focus: clicking the already-focused trend returns to the full radar.
  const toggleTrend = (slug: string) => {
    if (focused?.slug === slug) {
      navigate('/');
    } else {
      navigate(`/radar/${encodeURIComponent(slug)}`);
    }
  };
  const flat = useMemo(() => trends.flatMap((t) => t.items), [trends]);

  const totalCount = flat.length;
  const productCount = flat.filter((t) => t.stage === 'product').length;

  const focusedColor = focused
    ? PALETTE[allTrends.findIndex((x) => x.slug === focused.slug) % PALETTE.length]!
    : null;

  if (isLoading) {
    return (
      <div className="flex h-[70vh] items-center justify-center">
        <Spinner />
      </div>
    );
  }
  if (error) {
    return (
      <div className="rounded-card border border-red-500/40 bg-red-500/10 p-6 text-red-200">
        Ошибка загрузки: {(error as Error).message}
      </div>
    );
  }

  return (
    <>
      <PageHeader
        title={focused ? focused.name : 'Радар технологий'}
        description={
          focused
            ? 'Отдельный радар тренда. Колесо — масштаб, перетаскивание — пан. Клик по точке откроет инфо справа.'
            : 'Клик по названию тренда (или элементу слева) — откроет отдельный радар этого тренда.'
        }
        actions={
          focused ? (
            <Link to="/">
              <Button variant="ghost">← Все тренды</Button>
            </Link>
          ) : null
        }
        meta={<span>Локаль: {env.DEFAULT_LOCALE.toUpperCase()}</span>}
      />

      <div className="relative min-h-[480px] overflow-hidden rounded-2xl border border-line shadow-elevated max-[1100px]:h-auto"
           style={{ height: 'calc(100vh - 160px)' }}
      >
        <TrendsLegend
          trends={allTrends}
          activeSlug={focused?.slug ?? hoverTrend}
          onHover={setHoverTrend}
          onSelect={toggleTrend}
          totalTechs={totalCount}
          productCount={productCount}
        />

        <div className="absolute inset-0 max-[1100px]:relative max-[1100px]:aspect-square max-[1100px]:h-auto">
          <RadarChart
            ref={radarRef}
            trends={trends}
            flat={flat}
            fullList={allTrends}
            focusedSlug={focused?.slug ?? null}
            focusedColor={focusedColor}
            selectedSlug={selectedSlug}
            activeTrend={hoverTrend}
            onSelect={setSelectedSlug}
            onSectorClick={toggleTrend}
          />
        </div>

        <TechInfoPanel
          slug={selectedSlug}
          data={tech ?? null}
          loading={techLoading}
          onClose={() => setSelectedSlug(null)}
        />

        <RadarControls
          onZoomIn={() => radarRef.current?.zoomIn()}
          onZoomOut={() => radarRef.current?.zoomOut()}
          onReset={() => radarRef.current?.reset()}
        />
      </div>
    </>
  );
}

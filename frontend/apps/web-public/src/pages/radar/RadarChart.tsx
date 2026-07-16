import { forwardRef, memo, useImperativeHandle, useMemo, useRef } from 'react';
import type { HomeTechItem, HomeTrendBlock } from '@radartcell/api';
import {
  buildRadarConfig,
  computeBlips,
  computeSectors,
  type RadarConfig,
  type Sector,
  type Blip,
} from './radarGeometry';
import { useRadarPanZoom } from './useRadarPanZoom';

export interface RadarChartHandle {
  zoomIn: () => void;
  zoomOut: () => void;
  reset: () => void;
}

interface Props {
  trends: HomeTrendBlock[];
  flat: HomeTechItem[];
  fullList: HomeTrendBlock[];
  focusedSlug: string | null;
  focusedColor: string | null;
  selectedSlug: string | null;
  activeTrend: string | null;
  onSelect: (slug: string) => void;
  onSectorClick?: (trendSlug: string) => void;
}

const STAGE_ORDER = ['PRODUCTION', 'ADOPTING', 'EXPERIMENTING', 'ENVISIONING'];

function buildSectorPath(a0: number, a1: number, cfg: RadarConfig): string {
  const { cx, cy, ringInner, ringOuter } = cfg;
  const isFull = a1 - a0 >= Math.PI * 2 - 1e-3;
  const x0 = cx + Math.cos(a0) * ringOuter;
  const y0 = cy + Math.sin(a0) * ringOuter;
  const x1 = cx + Math.cos(a1) * ringOuter;
  const y1 = cy + Math.sin(a1) * ringOuter;
  const xi0 = cx + Math.cos(a0) * ringInner;
  const yi0 = cy + Math.sin(a0) * ringInner;
  const xi1 = cx + Math.cos(a1) * ringInner;
  const yi1 = cy + Math.sin(a1) * ringInner;

  if (isFull) {
    return (
      `M ${cx - ringOuter},${cy} ` +
      `A ${ringOuter},${ringOuter} 0 1 0 ${cx + ringOuter},${cy} ` +
      `A ${ringOuter},${ringOuter} 0 1 0 ${cx - ringOuter},${cy} ` +
      `M ${cx - ringInner},${cy} ` +
      `A ${ringInner},${ringInner} 0 1 1 ${cx + ringInner},${cy} ` +
      `A ${ringInner},${ringInner} 0 1 1 ${cx - ringInner},${cy}`
    );
  }
  const large = a1 - a0 > Math.PI ? 1 : 0;
  return (
    `M ${xi0.toFixed(1)},${yi0.toFixed(1)} ` +
    `L ${x0.toFixed(1)},${y0.toFixed(1)} ` +
    `A ${ringOuter},${ringOuter} 0 ${large} 1 ${x1.toFixed(1)},${y1.toFixed(1)} ` +
    `L ${xi1.toFixed(1)},${yi1.toFixed(1)} ` +
    `A ${ringInner},${ringInner} 0 ${large} 0 ${xi0.toFixed(1)},${yi0.toFixed(1)} Z`
  );
}

function buildTrendArcPath(sec: Sector, cfg: RadarConfig): string {
  const { cx, cy, trendArcR } = cfg;
  const isFull = sec.a1 - sec.a0 >= Math.PI * 2 - 1e-3;
  const arcPad = 0.04;
  const aa0 = sec.a0 + arcPad;
  const aa1 = sec.a1 - arcPad;
  const isBottom = Math.sin(sec.aMid) > 0;

  if (isFull) {
    const r = trendArcR;
    return (
      `M ${cx - r},${cy} A ${r},${r} 0 1 1 ${cx + r},${cy} ` +
      `A ${r},${r} 0 1 1 ${cx - r},${cy}`
    );
  }
  if (isBottom) {
    const r = trendArcR + 20;
    const sx = cx + Math.cos(aa1) * r;
    const sy = cy + Math.sin(aa1) * r;
    const ex = cx + Math.cos(aa0) * r;
    const ey = cy + Math.sin(aa0) * r;
    const large = aa1 - aa0 > Math.PI ? 1 : 0;
    return `M ${sx.toFixed(1)},${sy.toFixed(1)} A ${r},${r} 0 ${large} 0 ${ex.toFixed(1)},${ey.toFixed(1)}`;
  }
  const sx = cx + Math.cos(aa0) * trendArcR;
  const sy = cy + Math.sin(aa0) * trendArcR;
  const ex = cx + Math.cos(aa1) * trendArcR;
  const ey = cy + Math.sin(aa1) * trendArcR;
  const large = aa1 - aa0 > Math.PI ? 1 : 0;
  return `M ${sx.toFixed(1)},${sy.toFixed(1)} A ${trendArcR},${trendArcR} 0 ${large} 1 ${ex.toFixed(1)},${ey.toFixed(1)}`;
}

function BlipNode({
  blip,
  cfg,
  selected,
  dimmed,
  onSelect,
}: {
  blip: Blip;
  cfg: RadarConfig;
  selected: boolean;
  dimmed: boolean;
  onSelect: (slug: string) => void;
}) {
  const { x, y, angle, color, tech } = blip;
  const rayEnd = cfg.labelTextRadius - 6;
  const rxe = cfg.cx + Math.cos(angle) * rayEnd;
  const rye = cfg.cy + Math.sin(angle) * rayEnd;
  const tx = cfg.cx + Math.cos(angle) * cfg.labelTextRadius;
  const ty = cfg.cy + Math.sin(angle) * cfg.labelTextRadius;
  const deg = (angle * 180) / Math.PI;
  const flip = Math.cos(angle) < 0;
  const rotate = flip ? deg + 180 : deg;
  const anchor = flip ? 'end' : 'start';

  return (
    <g
      className={[
        'radar-dot cursor-pointer transition-opacity',
        dimmed && 'opacity-20',
        selected && 'radar-dot--selected',
      ]
        .filter(Boolean)
        .join(' ')}
      style={{ color }}
      onClick={(e) => {
        e.stopPropagation();
        onSelect(tech.slug);
      }}
    >
      <line
        x1={x.toFixed(1)}
        y1={y.toFixed(1)}
        x2={rxe.toFixed(1)}
        y2={rye.toFixed(1)}
        stroke={color}
        strokeOpacity={0.55}
        strokeWidth={1}
      />
      <circle
        cx={x.toFixed(1)}
        cy={y.toFixed(1)}
        r={11}
        fill={color}
        fillOpacity={0.3}
        filter="url(#dotGlow)"
      />
      <circle
        cx={x.toFixed(1)}
        cy={y.toFixed(1)}
        r={selected ? 7 : 4.5}
        fill={selected ? '#fff' : color}
        stroke={selected ? color : undefined}
        strokeWidth={selected ? 3 : undefined}
        filter="url(#dotGlow)"
      />
      <text
        x={tx.toFixed(1)}
        y={ty.toFixed(1)}
        fill={color}
        fontSize={12}
        fontWeight={selected ? 700 : 500}
        textAnchor={anchor}
        dominantBaseline="middle"
        transform={`rotate(${rotate.toFixed(2)} ${tx.toFixed(1)} ${ty.toFixed(1)})`}
        style={{ textShadow: `0 0 6px ${color}88` }}
      >
        {tech.name}
      </text>
      <title>{`${blip.no}. ${tech.name} · TRL ${tech.trl}`}</title>
    </g>
  );
}

export const RadarChart = memo(
  forwardRef<RadarChartHandle, Props>(function RadarChart(
    {
      trends,
      flat,
      fullList,
      focusedSlug,
      focusedColor,
      selectedSlug,
      activeTrend,
      onSelect,
      onSectorClick,
    }: Props,
    ref,
  ) {
  const isSingle = trends.length === 1;
  const cfg = useMemo(() => buildRadarConfig(flat, isSingle), [flat, isSingle]);
  const byTrend = useMemo(() => {
    const acc: Record<string, HomeTechItem[]> = {};
    // Group by the trend each item actually belongs to. The API nests items
    // under their trend, and flat items don't carry a trend_slug, so grouping
    // must come from the trend blocks themselves.
    for (const trend of trends) {
      acc[trend.slug] = trend.items ?? [];
    }
    return acc;
  }, [trends]);

  const sectors = useMemo(
    () => computeSectors(trends, byTrend, focusedSlug, fullList),
    [trends, byTrend, focusedSlug, fullList],
  );
  const blips = useMemo(() => computeBlips(sectors, byTrend, cfg), [sectors, byTrend, cfg]);

  const svgRef = useRef<SVGSVGElement>(null);
  const rootRef = useRef<SVGGElement>(null);
  const { zoomBy, reset } = useRadarPanZoom(svgRef, rootRef, { viewSize: cfg.view });

  useImperativeHandle(
    ref,
    () => ({
      zoomIn: () => zoomBy(1.25),
      zoomOut: () => zoomBy(1 / 1.25),
      reset,
    }),
    [zoomBy, reset],
  );

  const accent = focusedColor ?? '#7c3aed';
  const accentSoft = focusedColor ?? '#a78bfa';

  const dimmedTrend = activeTrend && !focusedSlug ? activeTrend : null;

  return (
    <svg
      ref={svgRef}
      viewBox={`0 0 ${cfg.view} ${cfg.view}`}
      preserveAspectRatio="xMidYMid meet"
      role="img"
      aria-labelledby="radarTitle radarDesc"
      className="block h-full w-full cursor-grab select-none touch-none animate-fade-in"
    >
      <title id="radarTitle">Радар технологий RadarTcell</title>
      <desc id="radarDesc">
        Круговая визуализация: сектора — тренды, точки — технологии. Кольца — стадии TRL.
      </desc>

      <defs>
        <radialGradient id="centerGlow" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor={accent} stopOpacity={0.22} />
          <stop offset="100%" stopColor={accent} stopOpacity={0} />
        </radialGradient>
        <filter id="dotGlow" x="-100%" y="-100%" width="300%" height="300%">
          <feGaussianBlur stdDeviation="2.5" result="coloredBlur" />
          <feMerge>
            <feMergeNode in="coloredBlur" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>
        <radialGradient id="scanRay" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor={accentSoft} stopOpacity={0.25} />
          <stop offset="60%" stopColor={accentSoft} stopOpacity={0.05} />
          <stop offset="100%" stopColor={accentSoft} stopOpacity={0} />
        </radialGradient>
      </defs>

      <g ref={rootRef}>
        <circle cx={cfg.cx} cy={cfg.cy} r={cfg.ringOuter} fill="url(#centerGlow)" />

        <g>
          <path
            d={`M ${cfg.cx},${cfg.cy} L ${cfg.cx},${(cfg.cy - cfg.ringOuter).toFixed(1)} A ${cfg.ringOuter},${cfg.ringOuter} 0 0 1 ${(cfg.cx + Math.cos(-Math.PI / 2 + Math.PI / 4) * cfg.ringOuter).toFixed(1)},${(cfg.cy + Math.sin(-Math.PI / 2 + Math.PI / 4) * cfg.ringOuter).toFixed(1)} Z`}
            fill="url(#scanRay)"
            opacity={0.55}
          >
            <animateTransform
              attributeName="transform"
              type="rotate"
              from={`0 ${cfg.cx} ${cfg.cy}`}
              to={`360 ${cfg.cx} ${cfg.cy}`}
              dur="8s"
              repeatCount="indefinite"
            />
          </path>
        </g>

        {sectors.map((sec, i) => {
          const arcId = `trendArc_${i}`;
          const isDim = dimmedTrend != null && sec.trend.slug !== dimmedTrend;
          return (
            <g key={sec.trend.slug} style={{ opacity: isDim ? 0.25 : 1 }}>
              <path
                d={buildSectorPath(sec.a0, sec.a1, cfg)}
                fill="transparent"
                onClick={() => onSectorClick?.(sec.trend.slug)}
                style={{ cursor: 'pointer' }}
              />
              {!isSingle && (
                <>
                  <path id={arcId} d={buildTrendArcPath(sec, cfg)} fill="none" stroke="none" />
                  <text
                    fill={sec.color}
                    fontSize={18}
                    fontWeight={800}
                    letterSpacing="0.06em"
                    style={{
                      textShadow: `0 0 10px ${sec.color}dd`,
                      textTransform: 'uppercase',
                      cursor: 'pointer',
                    }}
                    onClick={() => onSectorClick?.(sec.trend.slug)}
                  >
                    <textPath href={`#${arcId}`} startOffset="50%" textAnchor="middle">
                      {sec.trend.name}
                    </textPath>
                  </text>
                </>
              )}
            </g>
          );
        })}

        {cfg.ringR.map((r, i) => (
          <circle
            key={`ring-${i}`}
            cx={cfg.cx}
            cy={cfg.cy}
            r={r.toFixed(1)}
            fill="none"
            stroke="#1a2440"
            strokeWidth={1}
            strokeDasharray="2 4"
          />
        ))}
        {STAGE_ORDER.map((name, i) => {
          const rMid = (cfg.ringR[i]! + cfg.ringR[i + 1]!) / 2;
          return (
            <text
              key={name}
              x={cfg.cx.toFixed(1)}
              y={(cfg.cy - rMid).toFixed(1)}
              fill="#3e537e"
              fontSize={10}
              fontWeight={700}
              letterSpacing="0.18em"
              textAnchor="middle"
              dominantBaseline="middle"
            >
              {name}
            </text>
          );
        })}

        {blips.map((b) => (
          <BlipNode
            key={b.tech.slug}
            blip={b}
            cfg={cfg}
            selected={selectedSlug === b.tech.slug}
            dimmed={dimmedTrend != null && b.trendSlug !== dimmedTrend}
            onSelect={onSelect}
          />
        ))}

        <g>
          <circle
            cx={cfg.cx}
            cy={cfg.cy}
            r={cfg.ringInner - 6}
            fill={accent}
            fillOpacity={0.14}
            stroke={accent}
            strokeOpacity={0.6}
            strokeWidth={1}
          />
          <text
            x={cfg.cx}
            y={cfg.cy}
            fill="#fff"
            fontSize={12}
            fontWeight={700}
            letterSpacing="0.22em"
            textAnchor="middle"
            dominantBaseline="middle"
            style={{ textShadow: `0 0 12px ${accentSoft}` }}
          >
            RADARTCELL
          </text>
        </g>
      </g>
    </svg>
    );
  }),
);

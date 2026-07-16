import type { HomeTechItem, HomeTrendBlock } from '@radartcell/api';
import { PALETTE } from './constants';

export interface RadarConfig {
  size: number;
  pad: number;
  ringInner: number;
  ringOuter: number;
  cx: number;
  cy: number;
  view: number;
  ringR: number[];
  labelTextRadius: number;
  trendArcR: number;
}

export interface Sector {
  a0: number;
  a1: number;
  aMid: number;
  color: string;
  trend: HomeTrendBlock;
}

export interface Blip {
  tech: HomeTechItem;
  trendSlug: string;
  trendName: string;
  color: string;
  angle: number;
  x: number;
  y: number;
  ringIndex: number;
  no: number;
}

const RING_APPROX_CHAR = 7.0;

export function buildRadarConfig(flat: HomeTechItem[], isSingle: boolean): RadarConfig {
  const ringInner = 60;
  const ringOuter = 380;
  const ringR = [ringInner];
  for (let i = 1; i <= 4; i++) ringR.push(ringInner + (i / 4) * (ringOuter - ringInner));

  let maxNameLen = 0;
  for (const tech of flat) maxNameLen = Math.max(maxNameLen, (tech.name || '').length);
  const labelTextRadius = ringOuter + 14;
  const trendArcR = labelTextRadius + maxNameLen * RING_APPROX_CHAR + 28;

  // Actual outward extent of the drawn content, so the viewBox wraps it tightly
  // instead of leaving huge empty margins around a small radar.
  const blipLabelExtent = labelTextRadius + maxNameLen * 7.5 + 20; // dot leader labels
  const trendArcExtent = trendArcR + 46; // curved trend names (only when >1 trend)
  const contentRadius = isSingle
    ? blipLabelExtent
    : Math.max(blipLabelExtent, trendArcExtent);

  const margin = 24;
  const view = (contentRadius + margin) * 2;
  const cx = view / 2;
  const cy = view / 2;
  const size = ringOuter * 2;
  const pad = margin;

  return { size, pad, view, cx, cy, ringInner, ringOuter, ringR, labelTextRadius, trendArcR };
}

export function ringIndex(trl: number): number {
  const n = Math.max(1, Math.min(9, Number(trl) || 5));
  if (n >= 8) return 0; // PRODUCTION (inner)
  if (n >= 6) return 1; // ADOPTING
  if (n >= 4) return 2; // EXPERIMENTING
  return 3; // ENVISIONING (outer)
}

export function computeSectors(
  trends: HomeTrendBlock[],
  byTrend: Record<string, HomeTechItem[]>,
  focusedSlug: string | null,
  fullList: HomeTrendBlock[],
): Sector[] {
  const MIN_WEIGHT = 1;
  const weights = trends.map((t) => Math.max((byTrend[t.slug] || []).length, MIN_WEIGHT));
  const total = weights.reduce((s, w) => s + w, 0) || 1;

  const sectors: Sector[] = [];
  let acc = -Math.PI / 2;
  trends.forEach((trend, idx) => {
    const span = (weights[idx]! / total) * Math.PI * 2;
    const a0 = acc;
    const a1 = acc + span;
    acc = a1;
    const colorIdx = focusedSlug ? fullList.findIndex((x) => x.slug === trend.slug) : idx;
    const color = PALETTE[(colorIdx >= 0 ? colorIdx : 0) % PALETTE.length]!;
    sectors.push({ a0, a1, aMid: (a0 + a1) / 2, color, trend });
  });
  return sectors;
}

export function computeBlips(
  sectors: Sector[],
  byTrend: Record<string, HomeTechItem[]>,
  cfg: RadarConfig,
): Blip[] {
  const blips: Blip[] = [];
  let blipNo = 0;

  sectors.forEach((sec) => {
    const techs = byTrend[sec.trend.slug] || [];
    const localCount = Math.max(techs.length, 1);
    const seg = sec.a1 - sec.a0;

    const LABEL_LINE_PX = 16;
    const minStep = LABEL_LINE_PX / cfg.labelTextRadius;
    const step =
      localCount > 1
        ? Math.min(seg / localCount, Math.max(minStep, seg / Math.max(localCount, 6)))
        : 0;
    const totalSpan = step * (localCount - 1);
    const aStart = (sec.a0 + sec.a1) / 2 - totalSpan / 2;

    const ordered = techs.slice().sort((p, q) => (Number(p.trl) || 0) - (Number(q.trl) || 0));

    ordered.forEach((tech, localIdx) => {
      const angle = localCount === 1 ? (sec.a0 + sec.a1) / 2 : aStart + localIdx * step;
      const ri = ringIndex(tech.trl);
      const rInner = cfg.ringR[ri]!;
      const rOuter = cfg.ringR[ri + 1]!;
      const r = (rInner + rOuter) / 2;
      const x = cfg.cx + Math.cos(angle) * r;
      const y = cfg.cy + Math.sin(angle) * r;
      blipNo += 1;
      blips.push({
        tech,
        trendSlug: sec.trend.slug,
        trendName: sec.trend.name,
        color: sec.color,
        angle,
        x,
        y,
        ringIndex: ri,
        no: blipNo,
      });
    });
  });

  return blips;
}

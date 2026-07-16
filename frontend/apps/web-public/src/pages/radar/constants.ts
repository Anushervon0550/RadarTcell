// Ordered so that adjacent trends get maximally contrasting hues
// (violet → cyan → lime → orange → pink → blue → yellow → purple).
export const PALETTE = [
  '#a78bfa', // violet
  '#22d3ee', // cyan
  '#a3e635', // lime
  '#fb923c', // orange
  '#f472b6', // pink
  '#60a5fa', // blue
  '#facc15', // yellow
  '#c084fc', // purple
] as const;

export const FALLBACK_COVER =
  'https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=1200&q=80';

export function hashCode(str: string): number {
  let h = 0;
  for (let i = 0; i < str.length; i++) h = ((h << 5) - h + str.charCodeAt(i)) | 0;
  return h;
}

export function colorForTrend(slug: string, fallbackIdx = 0): string {
  if (!slug) return PALETTE[fallbackIdx % PALETTE.length]!;
  return PALETTE[Math.abs(hashCode(slug)) % PALETTE.length]!;
}

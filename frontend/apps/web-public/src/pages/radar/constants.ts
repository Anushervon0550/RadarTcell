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

// Pool of themed cover images used when a trend has no explicit image_url.
export const TREND_COVERS = [
  'https://images.unsplash.com/photo-1677442136019-21780ecad995?w=1200&q=80', // AI
  'https://images.unsplash.com/photo-1518770660439-4636190af475?w=1200&q=80', // networks/circuits
  'https://images.unsplash.com/photo-1558002038-1055907df827?w=1200&q=80', // IoT
  'https://images.unsplash.com/photo-1550751827-4bd374c3f58b?w=1200&q=80', // cybersecurity
  'https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=1200&q=80', // cloud/space
  'https://images.unsplash.com/photo-1556742044-3c52d6e88c62?w=1200&q=80', // fintech
  'https://images.unsplash.com/photo-1526628953301-3e589a6a8b74?w=1200&q=80', // data
  'https://images.unsplash.com/photo-1620712943543-bcc4688e7485?w=1200&q=80', // abstract tech
] as const;

// Explicit covers for the well-known seeded trends.
const TREND_COVER_BY_SLUG: Record<string, string> = {
  ai: TREND_COVERS[0],
  networks: TREND_COVERS[1],
  iot: TREND_COVERS[2],
  cyber: TREND_COVERS[3],
  'cloud-edge': TREND_COVERS[4],
  fintech: TREND_COVERS[5],
};

/** Pick a stable, distinct cover image for a trend by its slug. */
export function coverForTrend(slug: string): string {
  if (slug && TREND_COVER_BY_SLUG[slug]) return TREND_COVER_BY_SLUG[slug]!;
  if (!slug) return FALLBACK_COVER;
  return TREND_COVERS[Math.abs(hashCode(slug)) % TREND_COVERS.length]!;
}

export function hashCode(str: string): number {
  let h = 0;
  for (let i = 0; i < str.length; i++) h = ((h << 5) - h + str.charCodeAt(i)) | 0;
  return h;
}

export function colorForTrend(slug: string, fallbackIdx = 0): string {
  if (!slug) return PALETTE[fallbackIdx % PALETTE.length]!;
  return PALETTE[Math.abs(hashCode(slug)) % PALETTE.length]!;
}

// ============ PUBLIC APP STATE ============

export const S = {
  view: 'explore',
  locale: localStorage.getItem('locale') || 'ru',
  techs: [], total: 0, page: 1, limit: 20,
  trends: [], sdgs: [], tags: [], orgs: [], metrics: [],
  filters: {
    search: '', trlMin: 1, trlMax: 9,
    trends: new Set(), sdgs: new Set(), tags: new Set(), orgs: new Set(),
  },
  sort: { by: '', order: 'asc' },
  zoom: 1, pan: { x: 0, y: 0 }, selSlug: null,
  trendColor: new Map(),
};

const PALETTE = [
  '#c084fc', '#22d3ee', '#34d399', '#f472b6', '#fbbf24', '#60a5fa',
  '#a78bfa', '#fb7185', '#4ade80', '#facc15', '#818cf8', '#f97316',
];

export function colorForTrend(slug) {
  if (!S.trendColor.has(slug)) {
    S.trendColor.set(slug, PALETTE[S.trendColor.size % PALETTE.length]);
  }
  return S.trendColor.get(slug);
}

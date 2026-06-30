// ============ DATA LOADING ============
import { apiGet } from '../core/api.js';
import { toast } from '../core/util.js';
import { S, colorForTrend } from './state.js';
import { renderFilters } from './filters.js';
import { renderRadar } from './radar.js';
import { renderList } from './list.js';

export async function loadCatalogs() {
  try {
    const lo = '?locale=' + S.locale;
    const [tr, sd, tg, og, me] = await Promise.all([
      apiGet('/trends' + lo).catch(() => []),
      apiGet('/sdgs' + lo).catch(() => []),
      apiGet('/tags' + lo).catch(() => []),
      apiGet('/organizations' + lo).catch(() => []),
      apiGet('/metrics' + lo).catch(() => []),
    ]);
    S.trends = tr || []; S.sdgs = sd || []; S.tags = tg || [];
    S.orgs = og || []; S.metrics = me || [];
    S.trends.forEach((t) => colorForTrend(t.slug));
    renderFilters();
  } catch (e) {
    console.error(e);
    toast('Каталог недоступен', 'err');
  }
}

export async function loadTechs() {
  const f = S.filters;
  const q = new URLSearchParams();
  q.set('limit', String(S.limit));
  q.set('page', String(S.page));
  q.set('locale', S.locale);
  if (f.search) q.set('search', f.search);
  if (f.trlMin > 1) q.set('trl_min', f.trlMin);
  if (f.trlMax < 9) q.set('trl_max', f.trlMax);
  if (f.trends.size === 1) q.set('trend_id', [...f.trends][0]);
  if (f.sdgs.size === 1) q.set('sdg_id', [...f.sdgs][0]);
  if (f.tags.size === 1) q.set('tag_id', [...f.tags][0]);
  if (f.orgs.size === 1) q.set('organization_id', [...f.orgs][0]);
  if (S.sort.by) q.set('sort_by', S.sort.by);
  if (S.sort.order) q.set('order', S.sort.order);

  try {
    const r = await apiGet('/technologies?' + q);
    S.techs = r.items || [];
    S.total = r.total || 0;
    document.getElementById('ld').classList.add('hide');
    document.getElementById('meta').textContent = S.total + ' технологий';
    if (S.view === 'explore') renderRadar();
    if (S.view === 'list') renderList();
  } catch (e) {
    console.error(e);
    toast('Не удалось загрузить технологии', 'err');
    S.techs = [];
  }
}

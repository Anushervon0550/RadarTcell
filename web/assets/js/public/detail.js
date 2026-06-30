// ============ TECHNOLOGY DETAIL PANEL ============
import { apiGet } from '../core/api.js';
import { esc, escAt } from '../core/util.js';
import { S, colorForTrend } from './state.js';
import { renderRadar } from './radar.js';
import { openTrend, openSDG, openTag, openOrg } from './views.js';

export async function openTechDetail(slug) {
  S.selSlug = slug;
  if (S.view === 'explore') renderRadar();
  const d = document.getElementById('detail');
  const b = document.getElementById('dBody');
  d.classList.add('on');
  b.innerHTML = '<div style="padding:40px;color:var(--t3);font-family:var(--fm);font-size:11px;letter-spacing:2px">ЗАГРУЗКА…</div>';
  try {
    const t = await apiGet('/technologies/' + encodeURIComponent(slug) + '?locale=' + S.locale);
    b.innerHTML = detailTpl(t);
    b.querySelectorAll('[data-trend-slug]').forEach((e) => { e.onclick = () => { closeDetail(); openTrend(e.dataset.trendSlug); }; });
    b.querySelectorAll('[data-tag-slug]').forEach((e) => { e.onclick = () => { closeDetail(); openTag(e.dataset.tagSlug); }; });
    b.querySelectorAll('[data-sdg-code]').forEach((e) => { e.onclick = () => { closeDetail(); openSDG(e.dataset.sdgCode); }; });
    b.querySelectorAll('[data-org-slug]').forEach((e) => { e.onclick = () => { closeDetail(); openOrg(e.dataset.orgSlug); }; });
    const img = b.querySelector('img.img');
    if (img) img.addEventListener('error', () => { img.style.display = 'none'; });
  } catch {
    b.innerHTML = '<div style="padding:40px;color:var(--t3)">Не удалось загрузить.</div>';
  }
}

export function closeDetail() {
  document.getElementById('detail').classList.remove('on');
  S.selSlug = null;
  if (S.view === 'explore') renderRadar();
}

function detailTpl(t) {
  const trl = Math.max(1, Math.min(9, t.trl || 0));
  const trlBar = Array.from({ length: 9 }, (_, i) => `<span class="${i < trl ? 'on' : ''}"></span>`).join('');
  const tags = (t.tags || []).map((x) => `<span class="pill" data-tag-slug="${esc(x.slug)}">${esc(x.title)}</span>`).join('');
  const sdgs = (t.sdgs || []).map((x) => `<span class="pill" data-sdg-code="${esc(x.code)}">${esc(x.code)}</span>`).join('');
  const orgs = (t.organizations || []).map((x) => `<span class="pill" data-org-slug="${esc(x.slug)}">${esc(x.name)}</span>`).join('');
  const cm = (t.custom_metrics || []).map((m) => `<div class="row"><span class="k">${esc(m.field_key || m.metric_id)}</span><span class="v">${m.value != null ? (+m.value).toFixed(2) : '—'}</span></div>`).join('');
  return `
    ${t.image_url ? `<img src="${escAt(t.image_url)}" class="img" alt="${escAt(t.name || '')}"/>` : ''}
    <div class="idx">№ ${t.index || '—'}</div>
    <h1>${esc(t.name || '')}</h1>
    ${t.trend_name ? `<div class="badge" data-trend-slug="${esc(t.trend_slug)}" style="cursor:pointer;color:${colorForTrend(t.trend_slug)}">${esc(t.trend_name)}</div>` : ''}
    <p class="desc">${esc(t.description_full || t.description_short || 'Описание отсутствует.')}</p>
    <div class="row"><span class="k">TRL</span><span class="v">${trl}/9</span></div>
    <div class="trl-bar">${trlBar}</div>
    ${cm ? `<h3>Метрики</h3>${cm}` : ''}
    ${t.source_link ? `<a class="src" href="${escAt(t.source_link)}" target="_blank" rel="noopener">Источник →</a>` : ''}
    ${tags ? `<h3>Теги</h3><div class="pills">${tags}</div>` : ''}
    ${sdgs ? `<h3>SDG</h3><div class="pills">${sdgs}</div>` : ''}
    ${orgs ? `<h3>Организации</h3><div class="pills">${orgs}</div>` : ''}
  `;
}

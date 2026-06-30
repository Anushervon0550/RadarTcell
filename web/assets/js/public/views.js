// ============ CATALOG VIEWS + ROUTER ============
import { $$, esc } from '../core/util.js';
import { S, colorForTrend } from './state.js';
import { loadTechs } from './data.js';
import { renderFilters } from './filters.js';
import { renderRadar } from './radar.js';
import { renderList } from './list.js';

const CATALOG_VIEWS = ['trends', 'sdgs', 'tags', 'orgs', 'metrics'];

export function renderCatalog() {
  document.getElementById('grTrends').innerHTML =
    S.trends.map((t) => cardTpl(t.name, t.slug, `${t.technologies_count} технологий`, colorForTrend(t.slug))).join('')
    || '<div style="color:var(--t3)">пусто</div>';
  document.getElementById('grSdgs').innerHTML =
    S.sdgs.map((s) => cardTpl(`${s.code} — ${s.title}`, s.code, `${s.technologies_count} технологий`)).join('')
    || '<div style="color:var(--t3)">пусто</div>';
  document.getElementById('grTags').innerHTML =
    S.tags.map((t) => cardTpl(t.title, t.slug, t.category || '')).join('')
    || '<div style="color:var(--t3)">пусто</div>';
  document.getElementById('grOrgs').innerHTML =
    S.orgs.map((o) => cardTpl(o.name, o.slug, `${o.technologies_count} технологий`)).join('')
    || '<div style="color:var(--t3)">пусто</div>';
  document.getElementById('grMetrics').innerHTML =
    S.metrics.map((m) => cardTpl(m.name, m.id, (m.description || '') + ' · ' + m.type)).join('')
    || '<div style="color:var(--t3)">пусто</div>';

  $$('#v-trends .card').forEach((c) => { c.onclick = () => openTrend(c.dataset.id); });
  $$('#v-sdgs .card').forEach((c) => { c.onclick = () => openSDG(c.dataset.id); });
  $$('#v-tags .card').forEach((c) => { c.onclick = () => openTag(c.dataset.id); });
  $$('#v-orgs .card').forEach((c) => { c.onclick = () => openOrg(c.dataset.id); });
}

function cardTpl(title, id, sub, color) {
  return `<div class="card" data-id="${esc(id)}">
    <div class="top">${color ? `<div class="num" style="color:${color}">●</div>` : '<div class="num">→</div>'}</div>
    <div class="nm">${esc(title)}</div>
    <div class="ds">${esc(sub || '')}</div>
  </div>`;
}

/* ---- drill-downs: switch to filtered list ---- */
export async function openTrend(slug) {
  const t = S.trends.find((x) => x.slug === slug);
  if (!t) return;
  S.filters.trends.clear(); S.filters.trends.add(t.id);
  switchView('list', `Тренд: ${t.name}`);
  document.getElementById('listSub').textContent = `Технологии тренда «${t.name}»`;
  S.page = 1; renderFilters(); await loadTechs();
}
export async function openSDG(code) {
  const s = S.sdgs.find((x) => x.code === code);
  if (!s) return;
  S.filters.sdgs.clear(); S.filters.sdgs.add(s.id);
  switchView('list', `SDG: ${s.code}`);
  document.getElementById('listSub').textContent = s.title;
  S.page = 1; renderFilters(); await loadTechs();
}
export async function openTag(slug) {
  const t = S.tags.find((x) => x.slug === slug);
  if (!t) return;
  S.filters.tags.clear(); S.filters.tags.add(t.id);
  switchView('list', `Тег: ${t.title}`);
  document.getElementById('listSub').textContent = '';
  S.page = 1; renderFilters(); await loadTechs();
}
export async function openOrg(slug) {
  const o = S.orgs.find((x) => x.slug === slug);
  if (!o) return;
  S.filters.orgs.clear(); S.filters.orgs.add(o.id);
  switchView('list', `Организация: ${o.name}`);
  document.getElementById('listSub').textContent = '';
  S.page = 1; renderFilters(); await loadTechs();
}

/* ---- view router ---- */
export function switchView(v, subtitle) {
  S.view = v;
  $$('.sb-icon[data-view]').forEach((b) => b.classList.toggle('active', b.dataset.view === v));
  $$('.view').forEach((x) => x.classList.toggle('on', x.id === 'v-' + v));
  document.getElementById('viewName').textContent = (subtitle || v).toUpperCase();
  document.getElementById('drawer').classList.remove('on');
  if (v === 'explore') { if (!S.techs.length) loadTechs(); else renderRadar(); }
  if (v === 'list') { if (!S.techs.length) loadTechs(); else renderList(); }
  if (CATALOG_VIEWS.includes(v)) renderCatalog();
}

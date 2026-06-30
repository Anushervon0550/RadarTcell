// ============ LIST VIEW ============
import { $$, esc } from '../core/util.js';
import { S, colorForTrend } from './state.js';
import { loadTechs } from './data.js';
import { openTechDetail } from './detail.js';

export function renderList() {
  const grid = document.getElementById('listGrid');
  document.getElementById('listMeta').textContent = `${S.total} всего · стр. ${S.page}`;
  if (!S.techs.length) {
    grid.innerHTML = '<div style="color:var(--t3);padding:40px;text-align:center">Нет технологий</div>';
    document.getElementById('listPg').innerHTML = '';
    return;
  }
  grid.innerHTML = S.techs.map((t) => {
    const color = colorForTrend(t.trend_slug);
    const tag = t.trend_name
      ? `<div class="tag" style="color:${color};border-color:${color}33;background:${color}1a">${esc(t.trend_name)}</div>`
      : '';
    return `
    <div class="card" data-slug="${esc(t.slug)}">
      <div class="top">
        <div class="num">№ ${t.index || '—'}</div>
        ${tag}
      </div>
      <div class="nm">${esc(t.name)}</div>
      <div class="ds">${esc(t.description_short || '')}</div>
      <div class="foot"><span>${esc(t.slug)}</span><span class="trl">TRL ${t.trl}</span></div>
    </div>`;
  }).join('');
  grid.querySelectorAll('.card').forEach((c) => { c.onclick = () => openTechDetail(c.dataset.slug); });
  renderPag();
}

function renderPag() {
  const total = S.total, lim = S.limit, pgs = Math.max(1, Math.ceil(total / lim));
  const el = document.getElementById('listPg');
  if (pgs <= 1) { el.innerHTML = ''; return; }
  const cur = S.page;
  const pages = [];
  const add = (p) => pages.push(p);
  add(1);
  if (cur > 3) add('...');
  for (let i = Math.max(2, cur - 1); i <= Math.min(pgs - 1, cur + 1); i++) add(i);
  if (cur < pgs - 2) add('...');
  if (pgs > 1) add(pgs);
  el.innerHTML = `<button ${cur === 1 ? 'disabled' : ''} data-p="${cur - 1}">‹</button>`
    + pages.map((p) => (p === '...' ? '<span class="info">…</span>' : `<button class="${p === cur ? 'on' : ''}" data-p="${p}">${p}</button>`)).join('')
    + `<button ${cur === pgs ? 'disabled' : ''} data-p="${cur + 1}">›</button>`
    + `<span class="info">из ${pgs}</span>`;
  $$('button[data-p]', el).forEach((b) => { b.onclick = () => { S.page = +b.dataset.p; loadTechs(); }; });
}

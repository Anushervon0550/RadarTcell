// ============ FILTERS (drawer) ============
import { $, $$, esc } from '../core/util.js';
import { S } from './state.js';
import { loadTechs } from './data.js';

export function renderFilters() {
  fillChips('#g-trends .chip-list', S.trends, (t) => ({ id: t.id, label: t.name, n: t.technologies_count, key: 'trends' }));
  fillChips('#g-sdgs .chip-list', S.sdgs, (s) => ({ id: s.id, label: s.code, n: s.technologies_count, key: 'sdgs' }));
  fillChips('#g-tags .chip-list', S.tags, (t) => ({ id: t.id, label: t.title, n: null, key: 'tags' }));
  fillChips('#g-orgs .chip-list', S.orgs, (o) => ({ id: o.id, label: o.name, n: o.technologies_count, key: 'orgs' }));
}

function fillChips(sel, items, mk) {
  const el = $(sel);
  if (!el) return;
  el.innerHTML = '';
  if (!items || !items.length) {
    el.innerHTML = '<span style="font-size:11px;color:var(--t3)">нет данных</span>';
    return;
  }
  items.forEach((it) => {
    const m = mk(it);
    const c = document.createElement('div');
    c.className = 'chip';
    c.dataset.id = m.id;
    c.innerHTML = `<span>${esc(m.label)}</span>${m.n != null ? `<span class="n">${m.n}</span>` : ''}`;
    if (S.filters[m.key].has(m.id)) c.classList.add('on');
    c.onclick = () => {
      const set = S.filters[m.key];
      if (set.has(m.id)) set.delete(m.id); else set.add(m.id);
      c.classList.toggle('on');
      S.page = 1;
      loadTechs();
    };
    el.appendChild(c);
  });
}

/** Слайдеры TRL, кнопки сброса групп, сортировка списка. */
export function initFilters() {
  const a = document.getElementById('trlMin');
  const b = document.getElementById('trlMax');
  const av = document.getElementById('trlMinV');
  const bv = document.getElementById('trlMaxV');
  let trlTimer;
  function upd(ev) {
    let mi = +a.value, mx = +b.value;
    if (mi > mx) {
      if (ev && ev.target === a) b.value = mi; else a.value = mx;
      mi = +a.value; mx = +b.value;
    }
    av.textContent = mi; bv.textContent = mx;
    S.filters.trlMin = mi; S.filters.trlMax = mx;
    clearTimeout(trlTimer);
    trlTimer = setTimeout(() => { S.page = 1; loadTechs(); }, 200);
  }
  a.addEventListener('input', upd);
  b.addEventListener('input', upd);

  $$('[data-clr]').forEach((el) => {
    el.onclick = () => {
      const k = el.dataset.clr;
      if (k === 'trl') {
        S.filters.trlMin = 1; S.filters.trlMax = 9;
        a.value = 1; b.value = 9; av.textContent = '1'; bv.textContent = '9';
      } else {
        const map = { trend: 'trends', sdg: 'sdgs', tag: 'tags', org: 'orgs' };
        S.filters[map[k]].clear();
        $$(`#g-${k}s .chip.on`).forEach((c) => c.classList.remove('on'));
      }
      S.page = 1;
      loadTechs();
    };
  });

  const sortBy = document.getElementById('sortBy');
  const sortOrder = document.getElementById('sortOrder');
  const pageSize = document.getElementById('pageSize');
  sortBy.onchange = (e) => { S.sort.by = e.target.value; S.page = 1; loadTechs(); };
  sortOrder.onchange = (e) => { S.sort.order = e.target.value; S.page = 1; loadTechs(); };
  pageSize.onchange = (e) => { S.limit = +e.target.value; S.page = 1; loadTechs(); };
}

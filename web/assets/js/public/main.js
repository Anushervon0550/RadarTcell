// ============ PUBLIC APP ENTRY ============
import { $$, debounce } from '../core/util.js';
import { S } from './state.js';
import { loadCatalogs, loadTechs } from './data.js';
import { initFilters } from './filters.js';
import { initRadarControls } from './radar.js';
import { renderCatalog, switchView } from './views.js';
import { closeDetail } from './detail.js';

/* ---- sidebar navigation ---- */
$$('.sb-icon[data-view]').forEach((b) => { b.onclick = () => switchView(b.dataset.view); });

/* ---- search ---- */
const runSearch = debounce((value) => {
  S.filters.search = value.trim();
  S.page = 1;
  loadTechs();
}, 280);
document.getElementById('q').oninput = (e) => runSearch(e.target.value);

/* ---- keyboard shortcuts ---- */
document.addEventListener('keydown', (e) => {
  if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault();
    document.getElementById('q').focus();
  }
  if (e.key === 'Escape') {
    closeDetail();
    document.getElementById('drawer').classList.remove('on');
  }
});

/* ---- detail close + filters toggle ---- */
document.getElementById('dClose').onclick = closeDetail;
document.getElementById('btn-filters').onclick = () =>
  document.getElementById('drawer').classList.toggle('on');

/* ---- locale switcher ---- */
$$('#loc button').forEach((b) => {
  b.onclick = async () => {
    $$('#loc button').forEach((x) => x.classList.toggle('on', x === b));
    S.locale = b.dataset.l;
    localStorage.setItem('locale', S.locale);
    await loadCatalogs();
    await loadTechs();
    if (['trends', 'sdgs', 'tags', 'orgs', 'metrics'].includes(S.view)) renderCatalog();
  };
});
$$('#loc button').forEach((b) => b.classList.toggle('on', b.dataset.l === S.locale));

/* ---- init ---- */
initFilters();
initRadarControls();
(async function init() {
  await loadCatalogs();
  await loadTechs();
  switchView('explore');
})();

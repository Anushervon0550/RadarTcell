'use strict';

/* ============================ CONFIG / STATE ============================ */
const API_BASE = (window.__API_BASE__ || window.location.origin).replace(/\/$/, '');
const app = document.getElementById('app');
const statusEl = document.getElementById('status');
const apiBaseView = document.getElementById('apiBaseView');
const modalEl = document.getElementById('modal');
const modalWindow = document.getElementById('modalWindow');
apiBaseView.textContent = API_BASE;

const PALETTE = ['#7c3aed', '#3b82f6', '#22c55e', '#f59e0b', '#ef4444', '#e879f9', '#14b8a6', '#f97316'];
const FALLBACK_COVER = 'https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=1200&q=80';

const state = {
  token: localStorage.getItem('rt_admin_token') || '',
  locale: 'ru',
  catalog: { trends: [], tags: [], sdgs: [], organizations: [], metrics: [] },
  home: null,
  techCache: new Map(), // slug -> карточка
  filters: {
    search: '', trend_id: '', tag_id: '', sdg_id: '',
    organization_id: '', trl_min: '', trl_max: ''
  },
  activeTrend: '',     // подсветка тренда на радаре
  selectedSlug: '',    // выбранная технология (показывается в правой панели)
};

/* ============================ ADMIN CRUD CONFIG ============================ */
const adminConfigs = {
  trends: {
    title: 'Тренды',
    key: 'slug',
    list: '/api/admin/trends',
    create: '/api/admin/trends',
    update: '/api/admin/trends/{key}',
    remove: '/api/admin/trends/{key}',
    columns: ['slug', 'name', 'order_index'],
    fields: [
      { name: 'slug', label: 'Slug', type: 'text', requiredOnCreate: true },
      { name: 'name', label: 'Name', type: 'text', required: true },
      { name: 'order_index', label: 'Order', type: 'number', required: true },
      { name: 'description', label: 'Description', type: 'textarea' },
      { name: 'image_url', label: 'Image URL', type: 'text' },
    ],
  },
  tags: {
    title: 'Теги',
    key: 'slug',
    list: '/api/admin/tags',
    create: '/api/admin/tags',
    update: '/api/admin/tags/{key}',
    remove: '/api/admin/tags/{key}',
    columns: ['slug', 'title', 'category'],
    fields: [
      { name: 'slug', label: 'Slug', type: 'text', requiredOnCreate: true },
      { name: 'title', label: 'Title', type: 'text', required: true },
      { name: 'category', label: 'Category', type: 'text', required: true },
      { name: 'description', label: 'Description', type: 'textarea' },
    ],
  },
  organizations: {
    title: 'Организации',
    key: 'slug',
    list: '/api/admin/organizations',
    create: '/api/admin/organizations',
    update: '/api/admin/organizations/{key}',
    remove: '/api/admin/organizations/{key}',
    columns: ['slug', 'name', 'website'],
    fields: [
      { name: 'slug', label: 'Slug', type: 'text', requiredOnCreate: true },
      { name: 'name', label: 'Name', type: 'text', required: true },
      { name: 'logo_url', label: 'Logo URL', type: 'text' },
      { name: 'website', label: 'Website', type: 'text' },
      { name: 'headquarters', label: 'Headquarters', type: 'text' },
      { name: 'description', label: 'Description', type: 'textarea' },
    ],
  },
  metrics: {
    title: 'Метрики',
    key: 'id',
    list: '/api/admin/metrics',
    create: '/api/admin/metrics',
    update: '/api/admin/metrics/{key}',
    remove: '/api/admin/metrics/{key}',
    columns: ['id', 'name', 'type', 'field_key', 'orderable'],
    fields: [
      { name: 'name', label: 'Name', type: 'text', required: true },
      { name: 'type', label: 'Type', type: 'select', options: ['bubble', 'bar', 'distance'], required: true },
      { name: 'field_key', label: 'Field Key', type: 'text' },
      { name: 'orderable', label: 'Orderable', type: 'checkbox', required: true },
      { name: 'description', label: 'Description', type: 'textarea' },
    ],
  },
  sdgs: {
    title: 'ЦУР',
    key: 'code',
    list: '/api/admin/sdgs',
    create: '/api/admin/sdgs',
    update: '/api/admin/sdgs/{key}',
    remove: '/api/admin/sdgs/{key}',
    columns: ['code', 'title'],
    fields: [
      { name: 'code', label: 'Code', type: 'text', requiredOnCreate: true },
      { name: 'title', label: 'Title', type: 'text', required: true },
      { name: 'icon', label: 'Icon URL', type: 'text' },
      { name: 'description', label: 'Description', type: 'textarea' },
    ],
  },
  technologies: {
    title: 'Технологии',
    key: 'slug',
    list: '/api/admin/technologies?include_deleted=true&limit=200',
    create: '/api/admin/technologies',
    update: '/api/admin/technologies/{key}',
    remove: '/api/admin/technologies/{key}',
    restore: '/api/admin/technologies/{key}/restore',
    columns: ['slug', 'name', 'trend_slug', 'trl', 'deleted_at'],
    fields: [
      { name: 'slug', label: 'Slug', type: 'text', requiredOnCreate: true },
      { name: 'index', label: 'Index', type: 'number', required: true },
      { name: 'name', label: 'Name', type: 'text', required: true },
      { name: 'trend_slug', label: 'Trend Slug', type: 'text', required: true },
      { name: 'trl', label: 'TRL (1-9)', type: 'number', required: true },
      { name: 'source_link', label: 'Source Link', type: 'text' },
      { name: 'image_url', label: 'Image URL', type: 'text' },
      { name: 'description_short', label: 'Description Short', type: 'textarea' },
      { name: 'description_full', label: 'Description Full', type: 'textarea' },
      { name: 'tag_slugs', label: 'Tag slugs (через запятую)', type: 'text' },
      { name: 'sdg_codes', label: 'SDG codes (через запятую)', type: 'text' },
      { name: 'organization_slugs', label: 'Org slugs (через запятую)', type: 'text' },
      { name: 'custom_metric_1', label: 'Custom metric 1', type: 'number' },
      { name: 'custom_metric_2', label: 'Custom metric 2', type: 'number' },
      { name: 'custom_metric_3', label: 'Custom metric 3', type: 'number' },
      { name: 'custom_metric_4', label: 'Custom metric 4', type: 'number' },
    ],
  },
};

/* ============================ HELPERS ============================ */
function setStatus(message, ok) {
  statusEl.textContent = message;
  statusEl.className = 'status ' + (ok ? 'ok' : 'err');
  statusEl.style.display = 'block';
  window.clearTimeout(setStatus.timer);
  setStatus.timer = window.setTimeout(() => { statusEl.style.display = 'none'; }, 3500);
}

function escapeHtml(value) {
  return String(value == null ? '' : value).replace(/[&<>"']/g, (ch) => ({
    '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;'
  }[ch]));
}

function attr(value) { return escapeHtml(value); }

function hashCode(str) {
  let h = 0;
  for (let i = 0; i < str.length; i++) h = ((h << 5) - h + str.charCodeAt(i)) | 0;
  return h;
}

function colorForTrend(slug) {
  if (!slug) return PALETTE[0];
  return PALETTE[Math.abs(hashCode(slug)) % PALETTE.length];
}

function stageLabel(stage) {
  return ({ idea: 'Идея', prototype: 'Прототип', product: 'Продукт' })[stage] || stage || '—';
}

async function apiRequest(path, method = 'GET', body) {
  const headers = { 'Content-Type': 'application/json' };
  if (state.token) headers.Authorization = 'Bearer ' + state.token;
  const options = { method, headers };
  if (body !== undefined) options.body = JSON.stringify(body);

  const response = await fetch(API_BASE + path, options);
  if (response.status === 204) return null;
  const text = await response.text();
  let data = null;
  if (text) {
    try { data = JSON.parse(text); } catch (_e) { data = { error: text }; }
  }
  if (!response.ok) throw new Error((data && data.error) ? data.error : ('HTTP ' + response.status));
  return data;
}

async function preloadCatalog() {
  if (state.catalog.trends.length > 0) return;
  const [trends, tags, sdgs, orgs, metrics] = await Promise.all([
    apiRequest('/api/trends?locale=' + state.locale),
    apiRequest('/api/tags?locale=' + state.locale),
    apiRequest('/api/sdgs?locale=' + state.locale),
    apiRequest('/api/organizations?locale=' + state.locale),
    apiRequest('/api/metrics?locale=' + state.locale),
  ]);
  state.catalog.trends = trends || [];
  state.catalog.tags = tags || [];
  state.catalog.sdgs = sdgs || [];
  state.catalog.organizations = orgs || [];
  state.catalog.metrics = metrics || [];
}

function invalidateCaches() {
  state.home = null;
  state.techCache.clear();
  state.catalog = { trends: [], tags: [], sdgs: [], organizations: [], metrics: [] };
}

function parseRoute() {
  const raw = window.location.hash.replace(/^#\/?/, '');
  const parts = raw.split('/').filter(Boolean);
  return { raw, first: parts[0] || 'explore', second: parts[1] || '', third: parts[2] || '', parts };
}

function markMenu() {
  const current = '#/' + parseRoute().raw;
  document.querySelectorAll('.menu a').forEach((link) => {
    link.classList.toggle('active', current === link.getAttribute('href'));
  });
}

function renderLoading(title) {
  app.innerHTML =
    '<h1 class="page-header">' + escapeHtml(title || 'Загрузка') + '</h1>' +
    '<div class="card"><div class="spinner"></div><div class="tagline" style="text-align:center">Загружаем данные…</div></div>';
}

function renderError(title, err) {
  app.innerHTML = '<div class="card"><h2 class="section-title">' + escapeHtml(title) + '</h2>' +
    '<div style="color:#fecaca">' + escapeHtml(err.message || err) + '</div>' +
    '<div style="margin-top:12px"><button class="btn" onclick="router()">Повторить</button></div></div>';
}

/* ============================ MODAL ============================ */
function openModal(html) {
  modalWindow.innerHTML = html;
  modalEl.classList.add('open');
  modalEl.setAttribute('aria-hidden', 'false');
  document.body.style.overflow = 'hidden';
}
function closeModal() {
  modalEl.classList.remove('open');
  modalEl.setAttribute('aria-hidden', 'true');
  modalWindow.innerHTML = '';
  document.body.style.overflow = '';
  if ((window.location.hash || '').startsWith('#/technology/')) {
    history.replaceState(null, '', '#/explore');
  }
}
modalEl.addEventListener('click', (e) => {
  if (e.target.classList.contains('modal-backdrop') || e.target.closest('[data-modal-close]')) {
    closeModal();
  }
});
window.addEventListener('keydown', (e) => {
  if (e.key === 'Escape' && modalEl.classList.contains('open')) closeModal();
});

/* ============================ EXPLORE / RADAR ============================ */
async function renderExplore() {
  // Кеш: при возврате с карточки переисполнения избегаем.
  if (!state.home) {
    renderLoading('Радар технологий');
    state.home = await apiRequest('/api/home?limit=200&locale=' + state.locale);
  }
  const data = state.home;

  const trends = (data && data.trends) ? data.trends.filter((t) => (t.items || []).length > 0) : [];
  // Сглаживаем и присваиваем глобальный номер блипа в порядке отображения трендов.
  const flat = [];
  let blipNo = 0;
  trends.forEach((trend, idx) => {
    const trendColor = PALETTE[idx % PALETTE.length];
    (trend.items || []).forEach((item) => {
      blipNo += 1;
      flat.push({
        ...item,
        trend_slug: trend.slug,
        trend_name: trend.name,
        _color: trendColor,
        _no: blipNo,
      });
    });
  });
  state._blipMap = new Map(flat.map((b) => [b.slug, b]));

  const totalCount = flat.length;
  const trendCount = trends.length;
  const productCount = flat.filter((t) => t.stage === 'product').length;

  const summaryHtml =
    '<div class="radar-summary">' +
      '<div class="summary-cell"><div class="label">Технологий</div><div class="value">' + totalCount + '</div></div>' +
      '<div class="summary-cell"><div class="label">Трендов</div><div class="value">' + trendCount + '</div></div>' +
      '<div class="summary-cell"><div class="label">В продакшене</div><div class="value">' + productCount + '</div></div>' +
    '</div>';

  // Левая панель: компактная легенда трендов с цветами и подсчётом.
  const trendsHtml = trends.map((t, idx) => {
    const color = PALETTE[idx % PALETTE.length];
    const active = state.activeTrend === t.slug ? ' active' : '';
    return '<div class="legend-item' + active + '" data-trend="' + attr(t.slug) + '">' +
      '<span class="legend-swatch" style="background:' + color + ';box-shadow:0 0 8px ' + color + '"></span>' +
      '<span class="legend-name">' + escapeHtml(t.name) + '</span>' +
      '<span class="legend-count">' + (t.items || []).length + '</span>' +
    '</div>';
  }).join('');

  app.innerHTML =
    '<div class="page-header">' +
      '<div><h1>Радар технологий</h1>' +
      '<div class="sub">Перетащите радар мышкой, колесо — масштаб. Клик по точке или номеру откроет инфо справа.</div></div>' +
      '<div class="meta">Локаль: ' + escapeHtml(state.locale.toUpperCase()) + '</div>' +
    '</div>' +
    '<div class="explore-screen">' +
      '<aside class="trends-panel">' +
        '<h3>Тренды</h3>' +
        '<div class="radar-legend">' + (trendsHtml || '<div class="empty">Нет данных</div>') + '</div>' +
        '<div style="margin-top:14px">' + summaryHtml + '</div>' +
        '<div class="card-help" style="margin-top:14px">' +
          '<div class="help-title">Кольца</div>' +
          '<div class="help-list">' +
            '<div><span class="help-dot" style="background:#7c3aed"></span>Production · TRL 8–9</div>' +
            '<div><span class="help-dot" style="background:#3b82f6"></span>Adopting · TRL 6–7</div>' +
            '<div><span class="help-dot" style="background:#f59e0b"></span>Experimenting · TRL 4–5</div>' +
            '<div><span class="help-dot" style="background:#9ca3af"></span>Envisioning · TRL 1–3</div>' +
          '</div>' +
        '</div>' +
      '</aside>' +
      '<div class="radar-stage">' +
        '<div class="radar-card" id="radarCard"></div>' +
      '</div>' +
      '<div class="radar-zoom-indicator" id="radarZoomLabel">100%</div>' +
      '<div class="radar-controls">' +
        '<button type="button" data-zoom="in"  title="Приблизить">+</button>' +
        '<button type="button" data-zoom="out" title="Отдалить">−</button>' +
        '<button type="button" data-zoom="reset" title="Сбросить" style="font-size:14px">⟲</button>' +
      '</div>' +
      '<aside class="tech-panel collapsed" id="techPanel"></aside>' +
    '</div>';

  drawRadar(trends, flat);
  renderTechPanel();

  // Клики по элементам легенды слева
  app.querySelectorAll('.legend-item').forEach((el) => {
    el.onclick = () => toggleTrendHighlight(el.getAttribute('data-trend'));
  });
}

function drawRadar(trends, flat) {
  const card = document.getElementById('radarCard');
  if (!card) return;

  const SIZE = 920;
  // Большой PAD — нужно много места на длинные радиальные подписи.
  const PAD = 220;
  const VIEW = SIZE + PAD * 2;
  const cx = VIEW / 2;
  const cy = VIEW / 2;

  // Кольца стадий: от внешнего к внутреннему (как в EY Radar / ThoughtWorks).
  const RING_NAMES = ['ENVISIONING', 'EXPERIMENTING', 'ADOPTING', 'PRODUCTION'];
  // ENVISIONING = TRL 1..3, EXPERIMENTING = 4..5, ADOPTING = 6..7, PRODUCTION = 8..9
  const ringOuter = 380; // внешний радиус всех технологий
  const ringInner = 60;  // внутренний центральный круг
  // 4 кольца: r0=inner < r1 < r2 < r3 < r4=outer
  const ringR = [ringInner];
  for (let i = 1; i <= 4; i++) ringR.push(ringInner + (i / 4) * (ringOuter - ringInner));

  const totalTrends = Math.max(trends.length, 1);

  // Группируем по тренду.
  const byTrend = {};
  for (const tech of flat) {
    const k = tech.trend_slug || 'other';
    (byTrend[k] = byTrend[k] || []).push(tech);
  }

  // Маппинг TRL -> ring index (0..3), 0 = центральный «PRODUCTION», 3 = внешний «ENVISIONING»
  function ringIndex(trl) {
    const n = Math.max(1, Math.min(9, Number(trl) || 5));
    if (n >= 8) return 0;          // PRODUCTION (внутри)
    if (n >= 6) return 1;          // ADOPTING
    if (n >= 4) return 2;          // EXPERIMENTING
    return 3;                      // ENVISIONING (снаружи)
  }
  // Перевернём порядок для отображения (внутри хочется PRODUCTION, снаружи ENVISIONING).
  // ringR построены так: ringR[0] = inner; ringR[4] = outer. Кольцо k: между ringR[k] и ringR[k+1].
  // ringIndex 0 -> кольцо PRODUCTION между ringR[0]..ringR[1] (внутреннее).

  // ---------- кольца ----------
  const ringStrokes = [];
  for (let i = 0; i <= 4; i++) {
    ringStrokes.push(
      '<circle cx="' + cx + '" cy="' + cy + '" r="' + ringR[i].toFixed(1) +
      '" fill="none" stroke="#22304d" stroke-width="1" />'
    );
  }
  // Подписи стадий вдоль вертикальной оси (как в EY).
  const stageOrder = ['PRODUCTION', 'ADOPTING', 'EXPERIMENTING', 'ENVISIONING'];
  for (let i = 0; i < 4; i++) {
    const rMid = (ringR[i] + ringR[i + 1]) / 2;
    ringStrokes.push(
      '<text x="' + cx.toFixed(1) + '" y="' + (cy - rMid).toFixed(1) +
      '" fill="#3e537e" font-size="9" text-anchor="middle" dominant-baseline="middle"' +
      ' style="letter-spacing: 0.18em;font-weight:700">' + stageOrder[i] + '</text>'
    );
  }

  // ---------- сектора (фон + подписи трендов по дуге) ----------
  const sectorParts = [];
  trends.forEach((t, i) => {
    const a0 = (i / totalTrends) * Math.PI * 2 - Math.PI / 2;
    const a1 = ((i + 1) / totalTrends) * Math.PI * 2 - Math.PI / 2;
    const color = PALETTE[i % PALETTE.length];

    // Радиальная линия-разделитель сектора.
    const lx0 = cx + Math.cos(a0) * ringInner;
    const ly0 = cy + Math.sin(a0) * ringInner;
    const lx1 = cx + Math.cos(a0) * ringOuter;
    const ly1 = cy + Math.sin(a0) * ringOuter;
    sectorParts.push(
      '<line x1="' + lx0.toFixed(1) + '" y1="' + ly0.toFixed(1) +
      '" x2="' + lx1.toFixed(1) + '" y2="' + ly1.toFixed(1) +
      '" stroke="#22304d" stroke-width="1" />'
    );

    // Полупрозрачная цветовая зона — лёгкая, не отвлекает.
    const x0 = cx + Math.cos(a0) * ringOuter, y0 = cy + Math.sin(a0) * ringOuter;
    const x1 = cx + Math.cos(a1) * ringOuter, y1 = cy + Math.sin(a1) * ringOuter;
    const xi0 = cx + Math.cos(a0) * ringInner, yi0 = cy + Math.sin(a0) * ringInner;
    const xi1 = cx + Math.cos(a1) * ringInner, yi1 = cy + Math.sin(a1) * ringInner;
    const arcD =
      'M ' + xi0.toFixed(1) + ',' + yi0.toFixed(1) +
      ' L ' + x0.toFixed(1) + ',' + y0.toFixed(1) +
      ' A ' + ringOuter + ',' + ringOuter + ' 0 0 1 ' + x1.toFixed(1) + ',' + y1.toFixed(1) +
      ' L ' + xi1.toFixed(1) + ',' + yi1.toFixed(1) +
      ' A ' + ringInner + ',' + ringInner + ' 0 0 0 ' + xi0.toFixed(1) + ',' + yi0.toFixed(1) +
      ' Z';
    sectorParts.push(
      '<path class="radar-sector" data-trend="' + attr(t.slug) +
      '" d="' + arcD + '" fill="' + color + '" fill-opacity="0.06" />'
    );

    // Подпись тренда по дуге над внешним кольцом.
    const labelR = ringOuter + 18;
    const labelArcId = 'arc-' + i;
    // Чтобы текст не был «вверх ногами» в нижней половине круга — разворачиваем дугу.
    const aMid = (a0 + a1) / 2;
    const isUpper = Math.sin(aMid) < -0.1; // верхняя половина
    let pa0 = a0, pa1 = a1, sweep = 1;
    if (!isUpper) {
      // нижняя половина: рисуем дугу против часовой и текст пойдёт правильно
      pa0 = a1; pa1 = a0; sweep = 0;
    }
    const px0 = cx + Math.cos(pa0) * labelR;
    const py0 = cy + Math.sin(pa0) * labelR;
    const px1 = cx + Math.cos(pa1) * labelR;
    const py1 = cy + Math.sin(pa1) * labelR;
    sectorParts.push(
      '<defs><path id="' + labelArcId + '" d="' +
      'M ' + px0.toFixed(1) + ',' + py0.toFixed(1) +
      ' A ' + labelR + ',' + labelR + ' 0 0 ' + sweep + ' ' + px1.toFixed(1) + ',' + py1.toFixed(1) +
      '"/></defs>'
    );
    sectorParts.push(
      '<text class="radar-label" data-trend="' + attr(t.slug) + '" fill="' + color +
      '" font-size="14" font-weight="700" letter-spacing="0.15em"' +
      ' style="text-shadow: 0 1px 4px rgba(0,0,0,0.85);text-transform:uppercase">' +
        '<textPath href="#' + labelArcId + '" startOffset="50%" text-anchor="middle">' +
          escapeHtml(t.name) +
        '</textPath>' +
      '</text>'
    );
  });

  // ---------- блипы ----------
  // Распределяем технологии равномерно ВНУТРИ сектора по углу,
  // и кладём в радиус, соответствующий ring index.
  // Подпись — длинная, идёт радиально наружу от точки.
  const blipParts = [];
  const RAY_GAP_OUT = 18;  // отступ от точки до начала текста
  const TEXT_OFFSET = 4;   // дополнительный отступ для текста

  trends.forEach((trend, ti) => {
    const techs = byTrend[trend.slug] || [];
    const seg = (Math.PI * 2) / totalTrends;
    const a0 = (ti / totalTrends) * Math.PI * 2 - Math.PI / 2;
    const localCount = Math.max(techs.length, 1);
    const color = PALETTE[ti % PALETTE.length];

    techs.forEach((tech, localIdx) => {
      const angle = a0 + ((localIdx + 0.5) / localCount) * seg;

      const ri = ringIndex(tech.trl);
      const rInner = ringR[ri];
      const rOuter = ringR[ri + 1];
      // Лёгкий jitter, чтобы блипы внутри одного кольца не лежали идеально на одной окружности.
      const jitter = ((Math.cos(localIdx * 5.3 + ti * 2.7) + 1) / 2) * 0.5 + 0.25; // 0.25..0.75
      const r = rInner + jitter * (rOuter - rInner);

      const x = cx + Math.cos(angle) * r;
      const y = cy + Math.sin(angle) * r;

      // Луч от точки наружу.
      const rayEnd = ringOuter + 6;
      const rxe = cx + Math.cos(angle) * rayEnd;
      const rye = cy + Math.sin(angle) * rayEnd;

      // Положение текста: за внешним кольцом + ещё чуть-чуть.
      // Если подпись тренда идёт по дуге labelR (=ringOuter+18), название технологии
      // надо поставить ещё дальше — через TEXT_OFFSET.
      const textR = rayEnd + TEXT_OFFSET;
      const tx = cx + Math.cos(angle) * textR;
      const ty = cy + Math.sin(angle) * textR;
      // Поворот текста: радиально, наружу. Если sin>=0 (нижняя половина) — переворачиваем.
      const deg = (angle * 180) / Math.PI;
      // Стандартно textPath по углу: angle°, но в правой половине это нормально, в левой
      // получается «вверх ногами». Поэтому если cos<0 — поворачиваем на 180° и делаем text-anchor=end.
      const flip = Math.cos(angle) < 0;
      const rotate = flip ? deg + 180 : deg;
      const anchor = flip ? 'end' : 'start';
      const num = tech._no != null ? tech._no : '';

      blipParts.push(
        '<g class="radar-dot" data-slug="' + attr(tech.slug) +
        '" data-trend="' + attr(tech.trend_slug) + '" style="color:' + color + '">' +
          // тонкая радиальная линия от точки до края (как «луч»)
          '<line class="dot-ray" x1="' + x.toFixed(1) + '" y1="' + y.toFixed(1) +
          '" x2="' + rxe.toFixed(1) + '" y2="' + rye.toFixed(1) +
          '" stroke="' + color + '" stroke-opacity="0.45" stroke-width="1" />' +
          // glow вокруг точки
          '<circle class="dot-glow" cx="' + x.toFixed(1) + '" cy="' + y.toFixed(1) +
          '" r="9" fill="' + color + '" fill-opacity="0.25" />' +
          // ядро
          '<circle class="dot-core" cx="' + x.toFixed(1) + '" cy="' + y.toFixed(1) +
          '" r="4.5" fill="' + color + '" />' +
          // подпись радиально
          '<text class="dot-label" x="' + tx.toFixed(1) + '" y="' + ty.toFixed(1) +
          '" fill="#cbd5ec" font-size="11" text-anchor="' + anchor + '"' +
          ' dominant-baseline="middle"' +
          ' transform="rotate(' + rotate.toFixed(2) + ' ' + tx.toFixed(1) + ' ' + ty.toFixed(1) + ')">' +
          escapeHtml(tech.name) +
          '</text>' +
          '<title>' + (num ? num + '. ' : '') + escapeHtml(tech.name) + ' · TRL ' + escapeHtml(tech.trl) + '</title>' +
        '</g>'
      );
    });
  });

  // ---------- центр ----------
  const centerLabel =
    '<g class="radar-center">' +
      '<circle cx="' + cx + '" cy="' + cy + '" r="' + (ringInner - 6) +
      '" fill="rgba(124,58,237,0.10)" stroke="#3a4670" stroke-width="1" />' +
      '<text x="' + cx + '" y="' + cy + '" fill="#cbd5ec" font-size="12" font-weight="700"' +
      ' letter-spacing="0.22em" text-anchor="middle" dominant-baseline="middle">RADARTCELL</text>' +
    '</g>';

  // ---------- сборка SVG ----------
  card.insertAdjacentHTML('afterbegin',
    '<svg viewBox="0 0 ' + VIEW + ' ' + VIEW + '" preserveAspectRatio="xMidYMid meet">' +
      '<defs>' +
        '<radialGradient id="centerGlow" cx="50%" cy="50%" r="50%">' +
          '<stop offset="0%" stop-color="#7c3aed" stop-opacity="0.18" />' +
          '<stop offset="100%" stop-color="#7c3aed" stop-opacity="0" />' +
        '</radialGradient>' +
      '</defs>' +
      '<g id="radarRoot" transform="translate(0,0) scale(1)">' +
        '<circle cx="' + cx + '" cy="' + cy + '" r="' + ringOuter + '" fill="url(#centerGlow)" />' +
        sectorParts.join('') +
        ringStrokes.join('') +
        blipParts.join('') +
        centerLabel +
      '</g>' +
    '</svg>'
  );

  applyRadarHighlight(card);
  attachRadarInteractions(card, VIEW);
}

/* ------- pan/zoom + click selection ------- */
function attachRadarInteractions(card, viewSize) {
  const svg = card.querySelector('svg');
  const root = card.querySelector('#radarRoot');
  // Кнопки и индикатор зума живут вне card (на уровне сцены) — ищем у общего родителя.
  const stage = card.closest('.explore-screen') || document;
  const zoomLabel = stage.querySelector('#radarZoomLabel');
  if (!svg || !root) return;

  // Состояние трансформации в координатах SVG (0..viewSize).
  const t = { scale: 1, tx: 0, ty: 0 };
  const MIN = 0.5;
  const MAX = 5;

  function applyT() {
    root.setAttribute('transform', 'translate(' + t.tx + ',' + t.ty + ') scale(' + t.scale + ')');
    if (zoomLabel) zoomLabel.textContent = Math.round(t.scale * 100) + '%';
  }

  function clientToSvg(clientX, clientY) {
    const rect = svg.getBoundingClientRect();
    return {
      x: ((clientX - rect.left) / rect.width) * viewSize,
      y: ((clientY - rect.top) / rect.height) * viewSize,
    };
  }

  function zoomAt(svgX, svgY, factor) {
    const next = Math.max(MIN, Math.min(MAX, t.scale * factor));
    if (next === t.scale) return;
    const wx = (svgX - t.tx) / t.scale;
    const wy = (svgY - t.ty) / t.scale;
    t.scale = next;
    t.tx = svgX - wx * t.scale;
    t.ty = svgY - wy * t.scale;
    applyT();
  }

  function reset() { t.scale = 1; t.tx = 0; t.ty = 0; applyT(); }

  // Колесо — зум.
  svg.addEventListener('wheel', (e) => {
    e.preventDefault();
    const p = clientToSvg(e.clientX, e.clientY);
    const factor = e.deltaY < 0 ? 1.15 : 1 / 1.15;
    zoomAt(p.x, p.y, factor);
  }, { passive: false });

  // Pan: используем document-level события, без pointer capture, чтобы
  // не блокировать обычные click'и по точкам.
  let drag = null;
  let movedPx = 0;

  function onMove(e) {
    if (!drag) return;
    const dx = e.clientX - drag.startX;
    const dy = e.clientY - drag.startY;
    movedPx = Math.max(movedPx, Math.abs(dx) + Math.abs(dy));
    if (movedPx < 4) return; // мёртвая зона, чтобы не было дёрганий и не съедало click
    const rect = svg.getBoundingClientRect();
    const k = viewSize / rect.width;
    t.tx = drag.startTx + dx * k;
    t.ty = drag.startTy + dy * k;
    applyT();
  }
  function onUp() {
    if (!drag) return;
    drag = null;
    svg.classList.remove('dragging');
    document.removeEventListener('mousemove', onMove);
    document.removeEventListener('mouseup', onUp);
  }

  svg.addEventListener('mousedown', (e) => {
    if (e.button !== 0 && e.button !== 1) return;
    drag = { startX: e.clientX, startY: e.clientY, startTx: t.tx, startTy: t.ty };
    movedPx = 0;
    svg.classList.add('dragging');
    document.addEventListener('mousemove', onMove);
    document.addEventListener('mouseup', onUp);
  });

  // Клик по точке/сектору.
  svg.addEventListener('click', (e) => {
    if (movedPx > 4) return; // это был drag
    const dot = e.target.closest('.radar-dot');
    if (dot) {
      e.stopPropagation();
      selectTechnology(dot.getAttribute('data-slug'));
      return;
    }
    const sector = e.target.closest('.radar-sector, .radar-sector-arc, .radar-label');
    if (sector) toggleTrendHighlight(sector.getAttribute('data-trend'));
  });

  // Двойной клик — сбросить.
  svg.addEventListener('dblclick', (e) => {
    e.preventDefault();
    reset();
  });

  // Кнопки масштаба (живут на сцене, рядом с радаром).
  stage.querySelectorAll('[data-zoom]').forEach((btn) => {
    btn.addEventListener('click', (e) => {
      e.stopPropagation();
      const center = { x: viewSize / 2, y: viewSize / 2 };
      const action = btn.getAttribute('data-zoom');
      if (action === 'in')   zoomAt(center.x, center.y, 1.25);
      if (action === 'out')  zoomAt(center.x, center.y, 1 / 1.25);
      if (action === 'reset') reset();
    });
  });
}

function applyRadarHighlight(card) {
  const active = state.activeTrend;
  card.querySelectorAll('.radar-dot').forEach((el) => {
    const isDimmed = !!active && el.getAttribute('data-trend') !== active;
    const isSelected = state.selectedSlug && el.getAttribute('data-slug') === state.selectedSlug;
    el.classList.toggle('dimmed', isDimmed && !isSelected);
    el.classList.toggle('selected', !!isSelected);
  });
  card.querySelectorAll('.radar-sector, .radar-sector-arc, .radar-label').forEach((el) => {
    const same = !active || el.getAttribute('data-trend') === active;
    el.setAttribute('opacity', same ? '1' : '0.3');
  });
}

function toggleTrendHighlight(slug) {
  state.activeTrend = state.activeTrend === slug ? '' : slug;
  // Обновляем DOM «на месте», без полной перерисовки страницы.
  const card = document.getElementById('radarCard');
  if (card) applyRadarHighlight(card);
  document.querySelectorAll('.legend-item').forEach((el) => {
    el.classList.toggle('active', el.getAttribute('data-trend') === state.activeTrend);
  });
}

/* ------- right-side info panel ------- */
async function selectTechnology(slug) {
  state.selectedSlug = slug;
  const card = document.getElementById('radarCard');
  if (card) applyRadarHighlight(card);
  await renderTechPanel();
}

function clearTechnologySelection() {
  state.selectedSlug = '';
  const card = document.getElementById('radarCard');
  if (card) applyRadarHighlight(card);
  renderTechPanel();
}

async function renderTechPanel() {
  const panel = document.getElementById('techPanel');
  if (!panel) return;

  const slug = state.selectedSlug;
  if (!slug) {
    panel.classList.add('collapsed');
    panel.innerHTML = '';
    return;
  }
  panel.classList.remove('collapsed');

  // Loading state
  panel.innerHTML =
    '<div style="padding:30px;text-align:center"><div class="spinner"></div></div>';

  let tech;
  try {
    if (state.techCache.has(slug)) {
      tech = state.techCache.get(slug);
    } else {
      tech = await apiRequest('/api/technologies/' + encodeURIComponent(slug) + '?locale=' + state.locale);
      state.techCache.set(slug, tech);
    }
  } catch (err) {
    panel.innerHTML =
      '<div style="padding:30px;text-align:center;color:#fecaca">Ошибка: ' +
      escapeHtml(err.message || err) + '</div>';
    return;
  }

  if (state.selectedSlug !== slug) return; // пользователь успел выбрать другое

  const cover = tech.image_url || FALLBACK_COVER;
  const stage = stageFromTRL(tech.trl);

  const tagsHtml = (tech.tags || []).map((x) =>
    '<a class="pill accent" href="#/tag/' + encodeURIComponent(x.slug) + '">' + escapeHtml(x.title || x.slug) + '</a>'
  ).join('') || '<span class="tagline">—</span>';

  const sdgsHtml = (tech.sdgs || []).map((x) =>
    '<a class="pill" href="#/sdg/' + encodeURIComponent(x.code) + '">' + escapeHtml(x.code) + '</a>'
  ).join('') || '<span class="tagline">—</span>';

  const orgsHtml = (tech.organizations || []).slice(0, 6).map((x) => {
    return '<a class="pill" href="#/organization/' + encodeURIComponent(x.slug) + '">' + escapeHtml(x.name) + '</a>';
  }).join('') || '<span class="tagline">—</span>';

  const metricRow = (label, value) => {
    const pct = Math.max(0, Math.min(100, Math.round((Number(value) || 0) * 100)));
    return '<div class="metric-row">' +
      '<span class="label">' + escapeHtml(label) + '</span>' +
      '<span class="bar"><span style="width:' + pct + '%"></span></span>' +
      '<span class="pct">' + pct + '%</span>' +
    '</div>';
  };

  panel.innerHTML =
    '<div class="panel-scroll">' +
      '<div class="panel-hero" style="background-image:url(' + attr(cover) + ')">' +
        '<button class="panel-close" type="button" aria-label="Закрыть" data-close-panel>×</button>' +
      '</div>' +
      '<div class="panel-body">' +
        '<h2>' + escapeHtml(tech.name) + '</h2>' +
        '<div class="meta-row" style="margin-top:6px">' +
          '<span class="trl-chip">TRL ' + escapeHtml(tech.trl) + '</span>' +
          '<span class="stage-chip stage-' + stage + '">' + stageLabel(stage) + '</span>' +
          '<a class="pill" href="#/trend/' + encodeURIComponent(tech.trend_slug) + '">' + escapeHtml(tech.trend_name || tech.trend_slug) + '</a>' +
        '</div>' +

        (tech.description_short ? '<div class="section"><h4>Кратко</h4><p>' + escapeHtml(tech.description_short) + '</p></div>' : '') +
        (tech.description_full  ? '<div class="section"><h4>Подробно</h4><p>' + escapeHtml(tech.description_full) + '</p></div>' : '') +

        '<div class="section"><h4>Метрики</h4>' +
          metricRow('Зрелость',  tech.custom_metric_1) +
          metricRow('Влияние',   tech.custom_metric_2) +
          metricRow('Покрытие',  tech.custom_metric_3) +
          metricRow('Стоимость', tech.custom_metric_4) +
        '</div>' +

        '<div class="section"><h4>Теги</h4><div>' + tagsHtml + '</div></div>' +
        '<div class="section"><h4>ЦУР</h4><div>' + sdgsHtml + '</div></div>' +
        '<div class="section"><h4>Организации</h4><div>' + orgsHtml + '</div></div>' +

        '<div class="panel-actions">' +
          '<a class="btn primary sm" href="#/technology/' + encodeURIComponent(tech.slug) + '">Открыть полностью</a>' +
          (tech.source_link ? '<a class="btn ghost sm" href="' + attr(tech.source_link) + '" target="_blank" rel="noreferrer">Источник ↗</a>' : '') +
        '</div>' +
      '</div>' +
    '</div>';

  panel.querySelector('[data-close-panel]').addEventListener('click', clearTechnologySelection);
}

/* ============================ TECH DETAILS (MODAL) ============================ */
async function renderTechDetails(slug) {
  // Если фон ещё не отрисован (прямой заход по ссылке), отрисуем радар.
  if (!app.firstChild || !app.querySelector('.radar-card')) {
    try { await renderExplore(); } catch (_e) { /* ignore */ }
  }
  let tech;
  try {
    if (state.techCache.has(slug)) {
      tech = state.techCache.get(slug);
    } else {
      tech = await apiRequest('/api/technologies/' + encodeURIComponent(slug) + '?locale=' + state.locale);
      state.techCache.set(slug, tech);
    }
  } catch (err) {
    setStatus('Не удалось загрузить технологию: ' + (err.message || err), false);
    return;
  }

  const cover = tech.image_url || FALLBACK_COVER;
  const stage = stageFromTRL(tech.trl);
  const completion = Math.max(0, Math.min(100, Math.round(((tech.custom_metric_1 || 0)) * 100)));

  const tags = (tech.tags || []).map((x) =>
    '<a class="pill accent" href="#/tag/' + encodeURIComponent(x.slug) + '">' + escapeHtml(x.title || x.slug) + '</a>'
  ).join('') || '<span class="tagline">не указаны</span>';

  const sdgs = (tech.sdgs || []).map((x) =>
    '<a class="pill" href="#/sdg/' + encodeURIComponent(x.code) + '">' + escapeHtml(x.code) + ' · ' + escapeHtml(x.title || '') + '</a>'
  ).join('') || '<span class="tagline">не указаны</span>';

  const orgs = (tech.organizations || []).map((x) => {
    const logo = x.logo_url ? 'background-image:url(' + attr(x.logo_url) + ')' : '';
    return '<a class="org-row" href="#/organization/' + encodeURIComponent(x.slug) + '" style="text-decoration:none;color:inherit">' +
      '<span class="org-logo" style="' + logo + '"></span>' +
      '<span><div style="font-weight:600">' + escapeHtml(x.name) + '</div>' +
      '<div class="tagline">' + escapeHtml(x.headquarters || x.website || '') + '</div></span>' +
      '</a>';
  }).join('') || '<span class="tagline">нет связанных</span>';

  const metricCells = [
    ['Зрелость',   tech.custom_metric_1],
    ['Влияние',    tech.custom_metric_2],
    ['Покрытие',   tech.custom_metric_3],
    ['Стоимость',  tech.custom_metric_4],
  ].map(([label, value]) => {
    const pct = Math.round((Number(value) || 0) * 100);
    return '<div class="metric-cell"><div class="label">' + label + '</div>' +
      '<div class="value">' + pct + '%</div>' +
      '<div class="bar"><span style="width:' + pct + '%"></span></div></div>';
  }).join('');

  const html =
    '<button class="modal-close" data-modal-close aria-label="Закрыть">×</button>' +
    '<div class="detail-hero" style="background-image:url(' + attr(cover) + ')"></div>' +
    '<div class="detail-body">' +
      '<div class="detail-title-row">' +
        '<h2 class="detail-title">' + escapeHtml(tech.name) + '</h2>' +
        '<div class="detail-row-meta">' +
          '<span class="trl-chip">TRL ' + escapeHtml(tech.trl) + '</span>' +
          '<span class="stage-chip stage-' + stage + '">' + stageLabel(stage) + '</span>' +
          '<a class="pill accent" href="#/trend/' + encodeURIComponent(tech.trend_slug) + '">' + escapeHtml(tech.trend_name || tech.trend_slug) + '</a>' +
          (tech.source_link ? '<a class="pill" href="' + attr(tech.source_link) + '" target="_blank" rel="noreferrer">Источник ↗</a>' : '') +
        '</div>' +
      '</div>' +

      (tech.description_short ? '<div class="detail-section"><h3>Краткое описание</h3><p>' + escapeHtml(tech.description_short) + '</p></div>' : '') +
      (tech.description_full ? '<div class="detail-section"><h3>Подробно</h3><p>' + escapeHtml(tech.description_full) + '</p></div>' : '') +

      '<div class="detail-section"><h3>Метрики</h3><div class="detail-grid">' + metricCells + '</div></div>' +

      '<div class="detail-section"><h3>Стадия развития</h3>' +
        '<p style="margin:0">Готовность: <strong>' + completion + '%</strong> · ' +
        'Уровень зрелости: <strong>TRL ' + escapeHtml(tech.trl) + '</strong> · ' +
        'Класс: <strong>' + stageLabel(stage) + '</strong></p>' +
        '<div class="bar" style="margin-top:10px"><span style="width:' + Math.max(10, Math.round((tech.trl/9)*100)) + '%"></span></div>' +
      '</div>' +

      '<div class="detail-section"><h3>Теги</h3><div>' + tags + '</div></div>' +
      '<div class="detail-section"><h3>Цели устойчивого развития</h3><div>' + sdgs + '</div></div>' +
      '<div class="detail-section"><h3>Организации</h3><div class="detail-grid">' + orgs + '</div></div>' +

    '</div>';

  openModal(html);
}

function stageFromTRL(trl) {
  const n = Number(trl) || 0;
  if (n <= 3) return 'idea';
  if (n <= 6) return 'prototype';
  return 'product';
}

/* ============================ CATALOG ============================ */
function buildFiltersQuery() {
  const q = new URLSearchParams();
  Object.keys(state.filters).forEach((k) => {
    const v = state.filters[k];
    if (v !== '' && v !== null && v !== undefined) q.set(k, v);
  });
  q.set('limit', '200');
  q.set('locale', state.locale);
  return q;
}

async function renderCatalog() {
  renderLoading('Каталог технологий');
  await preloadCatalog();
  const data = await apiRequest('/api/technologies?' + buildFiltersQuery().toString());
  const rows = (data && data.items) ? data.items : [];

  const trendOptions = ['<option value="">Все тренды</option>']
    .concat(state.catalog.trends.map((t) => '<option value="' + attr(t.id) + '">' + escapeHtml(t.name) + '</option>')).join('');
  const tagOptions = ['<option value="">Все теги</option>']
    .concat(state.catalog.tags.map((t) => '<option value="' + attr(t.id) + '">' + escapeHtml(t.title || t.slug) + '</option>')).join('');
  const sdgOptions = ['<option value="">Все ЦУР</option>']
    .concat(state.catalog.sdgs.map((s) => '<option value="' + attr(s.id) + '">' + escapeHtml(s.code) + '</option>')).join('');
  const orgOptions = ['<option value="">Все организации</option>']
    .concat(state.catalog.organizations.map((o) => '<option value="' + attr(o.id) + '">' + escapeHtml(o.name) + '</option>')).join('');

  const cards = rows.map((it) => techCardHtml(it)).join('');

  app.innerHTML =
    '<div class="page-header">' +
      '<div><h1>Каталог технологий</h1>' +
        '<div class="sub">Найдено: ' + rows.length + ' из ' + (data.total || rows.length) + '</div></div>' +
    '</div>' +
    '<div class="card" style="margin-bottom:18px"><div class="toolbar">' +
      '<input class="field" id="fSearch" placeholder="Поиск по названию" value="' + attr(state.filters.search) + '">' +
      '<select id="fTrend">' + trendOptions + '</select>' +
      '<select id="fTag">' + tagOptions + '</select>' +
      '<select id="fSdg">' + sdgOptions + '</select>' +
      '<select id="fOrg">' + orgOptions + '</select>' +
      '<input class="field" id="fTrlMin" type="number" min="1" max="9" placeholder="TRL мин" value="' + attr(state.filters.trl_min) + '" style="min-width:90px">' +
      '<input class="field" id="fTrlMax" type="number" min="1" max="9" placeholder="TRL макс" value="' + attr(state.filters.trl_max) + '" style="min-width:90px">' +
      '<button class="btn primary" id="applyFilters">Применить</button>' +
      '<button class="btn ghost" id="resetFilters">Сброс</button>' +
    '</div></div>' +
    '<div class="grid cards" id="catalogGrid">' + (cards || '<div class="empty">Ничего не найдено по текущим фильтрам.</div>') + '</div>';

  document.getElementById('fTrend').value = state.filters.trend_id;
  document.getElementById('fTag').value = state.filters.tag_id;
  document.getElementById('fSdg').value = state.filters.sdg_id;
  document.getElementById('fOrg').value = state.filters.organization_id;

  document.getElementById('applyFilters').onclick = () => {
    state.filters.search = document.getElementById('fSearch').value.trim();
    state.filters.trend_id = document.getElementById('fTrend').value;
    state.filters.tag_id = document.getElementById('fTag').value;
    state.filters.sdg_id = document.getElementById('fSdg').value;
    state.filters.organization_id = document.getElementById('fOrg').value;
    state.filters.trl_min = document.getElementById('fTrlMin').value;
    state.filters.trl_max = document.getElementById('fTrlMax').value;
    router();
  };
  document.getElementById('resetFilters').onclick = () => {
    state.filters = { search: '', trend_id: '', tag_id: '', sdg_id: '', organization_id: '', trl_min: '', trl_max: '' };
    router();
  };

  bindCardClicks();
}

function techCardHtml(it) {
  const stage = stageFromTRL(it.trl);
  const cover = it.image_url || FALLBACK_COVER;
  return '<div class="card tech-card" data-slug="' + attr(it.slug) + '">' +
    '<div class="cover" style="background-image:url(' + attr(cover) + ')"></div>' +
    '<div class="body">' +
      '<h3>' + escapeHtml(it.name) + '</h3>' +
      '<div class="meta-row">' +
        '<span class="trl-chip">TRL ' + escapeHtml(it.trl) + '</span>' +
        '<span class="stage-chip stage-' + stage + '">' + stageLabel(stage) + '</span>' +
        '<span class="pill">' + escapeHtml(it.trend_name || it.trend_slug) + '</span>' +
      '</div>' +
      '<p class="tagline">' + escapeHtml(it.description_short || '') + '</p>' +
    '</div>' +
  '</div>';
}

function bindCardClicks() {
  app.querySelectorAll('.tech-card').forEach((el) => {
    el.addEventListener('click', () => {
      const slug = el.getAttribute('data-slug');
      window.location.hash = '#/technology/' + encodeURIComponent(slug);
    });
  });
}

/* ============================ ENTITY LISTS ============================ */
async function renderTrends() {
  renderLoading('Тренды');

  // Public API не возвращает image_url у трендов. Если есть админ-токен, запросим
  // дополнительно админ-список ради картинок и описаний.
  const publicList = await apiRequest('/api/trends?locale=' + state.locale);
  let adminBySlug = null;
  if (state.token) {
    try {
      const admin = await apiRequest('/api/admin/trends');
      adminBySlug = {};
      (Array.isArray(admin) ? admin : (admin.items || [])).forEach((t) => { adminBySlug[t.slug] = t; });
    } catch (_e) { /* ignore */ }
  }

  const items = Array.isArray(publicList) ? publicList : (publicList.items || []);
  const cards = items.map((t) => {
    const adm = adminBySlug && adminBySlug[t.slug] ? adminBySlug[t.slug] : {};
    const img = adm.image_url || FALLBACK_COVER;
    const desc = adm.description || ('Технологий: ' + (t.technologies_count || 0));
    return '<div class="card trend-card" onclick="location.hash=\'#/trend/' + encodeURIComponent(t.slug) + '\'">' +
      '<div class="cover" style="background-image:url(' + attr(img) + ')"></div>' +
      '<div class="body">' +
        '<h3>' + escapeHtml(t.name) + '</h3>' +
        '<div class="tagline">' + escapeHtml(desc) + '</div>' +
        '<div class="meta-row" style="margin-top:6px"><span class="pill">Технологий: ' + (t.technologies_count || 0) + '</span></div>' +
      '</div></div>';
  }).join('');
  app.innerHTML = '<div class="page-header"><h1>Тренды</h1></div>' +
    '<div class="grid cards">' + (cards || '<div class="empty">Нет данных</div>') + '</div>';
}

async function renderOrganizations() {
  renderLoading('Организации');
  const list = await apiRequest('/api/organizations?locale=' + state.locale);
  const cards = (list || []).map((o) => {
    const logo = o.logo_url ? 'background-image:url(' + attr(o.logo_url) + ')' : '';
    return '<div class="card org-card" onclick="location.hash=\'#/organization/' + encodeURIComponent(o.slug) + '\'">' +
      '<div class="logo-box" style="' + logo + '"></div>' +
      '<div><h3 style="margin-bottom:4px">' + escapeHtml(o.name) + '</h3>' +
      '<div class="tagline">' + escapeHtml(o.headquarters || '') + '</div>' +
      '<div class="tagline">Технологий: ' + (o.technologies_count || 0) + '</div></div>' +
      '</div>';
  }).join('');
  app.innerHTML = '<div class="page-header"><h1>Организации</h1></div>' +
    '<div class="grid cards">' + (cards || '<div class="empty">Нет данных</div>') + '</div>';
}

async function renderTags() {
  renderLoading('Теги');
  const list = await apiRequest('/api/tags?locale=' + state.locale);
  const cards = (list || []).map((t) => {
    return '<div class="card" style="cursor:pointer" onclick="location.hash=\'#/tag/' + encodeURIComponent(t.slug) + '\'">' +
      '<h3>' + escapeHtml(t.title) + '</h3>' +
      '<div class="tagline">' + escapeHtml(t.category || '') + '</div>' +
      (t.description ? '<p class="tagline">' + escapeHtml(t.description) + '</p>' : '') +
    '</div>';
  }).join('');
  app.innerHTML = '<div class="page-header"><h1>Теги</h1></div>' +
    '<div class="grid cards">' + (cards || '<div class="empty">Нет данных</div>') + '</div>';
}

async function renderSDGs() {
  renderLoading('Цели устойчивого развития');
  const list = await apiRequest('/api/sdgs?locale=' + state.locale);
  const cards = (list || []).map((s) => {
    return '<div class="card" style="cursor:pointer" onclick="location.hash=\'#/sdg/' + encodeURIComponent(s.code) + '\'">' +
      '<h3>' + escapeHtml(s.code) + '</h3>' +
      '<div class="tagline" style="margin-bottom:6px">' + escapeHtml(s.title) + '</div>' +
      '<div class="tagline">Технологий: ' + (s.technologies_count || 0) + '</div>' +
    '</div>';
  }).join('');
  app.innerHTML = '<div class="page-header"><h1>Цели устойчивого развития</h1></div>' +
    '<div class="grid cards">' + (cards || '<div class="empty">Нет данных</div>') + '</div>';
}

async function renderTechByEntity(type, value) {
  renderLoading('Технологии: ' + value);
  const map = {
    trend: '/api/trends/' + encodeURIComponent(value) + '/technologies?limit=200&locale=' + state.locale,
    tag: '/api/tags/' + encodeURIComponent(value) + '/technologies?limit=200&locale=' + state.locale,
    sdg: '/api/sdgs/' + encodeURIComponent(value) + '/technologies?limit=200&locale=' + state.locale,
    organization: '/api/organizations/' + encodeURIComponent(value) + '/technologies?limit=200&locale=' + state.locale,
  };
  const data = await apiRequest(map[type]);
  const rows = (data && data.items) ? data.items : [];
  const cards = rows.map((it) => techCardHtml(it)).join('');
  const titles = { trend: 'Тренд', tag: 'Тег', sdg: 'ЦУР', organization: 'Организация' };

  app.innerHTML =
    '<div class="page-header">' +
      '<div><h1>' + escapeHtml(titles[type]) + ': ' + escapeHtml(value) + '</h1>' +
      '<div class="sub">Технологий найдено: ' + rows.length + '</div></div>' +
      '<a class="btn ghost" href="#/explore">← К радару</a>' +
    '</div>' +
    '<div class="grid cards">' + (cards || '<div class="empty">Нет связанных технологий</div>') + '</div>';

  bindCardClicks();
}

/* ============================ ADMIN ============================ */
function renderAdminLogin() {
  app.innerHTML = '<div class="page-header"><h1>Вход в админку</h1></div>' +
    '<div class="card" style="max-width:520px">' +
      '<div class="form-row"><label>Логин</label><input class="field" id="adminUser" placeholder="username" value="admin"></div>' +
      '<div class="form-row" style="margin-top:10px"><label>Пароль</label><input class="field" id="adminPass" placeholder="password" type="password" value="admin123"></div>' +
      '<div class="form-actions">' +
        '<button class="btn primary" id="adminLoginBtn">Войти</button>' +
        (state.token ? '<button class="btn danger" id="adminLogoutBtn">Выйти</button>' : '') +
      '</div>' +
      '<div class="tagline" id="adminSessionInfo" style="margin-top:8px">Токен: ' + (state.token ? 'сохранён' : 'не найден') + '</div>' +
    '</div>';

  document.getElementById('adminLoginBtn').onclick = async () => {
    try {
      const username = document.getElementById('adminUser').value.trim();
      const password = document.getElementById('adminPass').value;
      const data = await apiRequest('/api/admin/login', 'POST', { username, password });
      state.token = data.token || '';
      localStorage.setItem('rt_admin_token', state.token);
      setStatus('Вход выполнен', true);
      window.location.hash = '#/admin/technologies';
    } catch (err) {
      setStatus('Ошибка входа: ' + (err.message || err), false);
    }
  };
  const logout = document.getElementById('adminLogoutBtn');
  if (logout) {
    logout.onclick = () => {
      state.token = '';
      localStorage.removeItem('rt_admin_token');
      setStatus('Выход выполнен', true);
      router();
    };
  }
}

async function renderAdminEntity(name) {
  if (!state.token) { window.location.hash = '#/admin/login'; return; }
  const cfg = adminConfigs[name];
  if (!cfg) { app.innerHTML = '<div class="card">Неизвестный раздел</div>'; return; }

  renderLoading(cfg.title);
  let data;
  try { data = await apiRequest(cfg.list); }
  catch (err) {
    if ((err.message || '').includes('401')) { window.location.hash = '#/admin/login'; return; }
    throw err;
  }
  const list = Array.isArray(data) ? data : (data.items || []);

  const cols = cfg.columns;
  const tableHead = cols.map((c) => '<th>' + escapeHtml(c) + '</th>').join('');
  const rows = list.map((item) => {
    const cellHtml = cols.map((c) => '<td>' + escapeHtml(item[c]) + '</td>').join('');
    const keyValue = item[cfg.key];
    const actions = '<td>' +
      '<button class="btn sm" data-action="edit" data-key="' + attr(keyValue) + '">✎</button> ' +
      '<button class="btn sm danger" data-action="delete" data-key="' + attr(keyValue) + '">×</button> ' +
      (cfg.restore && item.deleted_at ? '<button class="btn sm warn" data-action="restore" data-key="' + attr(keyValue) + '">↻</button>' : '') +
    '</td>';
    return '<tr>' + cellHtml + actions + '</tr>';
  }).join('');

  app.innerHTML =
    '<div class="page-header"><h1>Админ · ' + escapeHtml(cfg.title) + '</h1></div>' +
    '<div class="split">' +
      '<div class="card table-wrap">' +
        '<table><thead><tr>' + tableHead + '<th>Действия</th></tr></thead><tbody>' + (rows || '') + '</tbody></table>' +
      '</div>' +
      '<div class="card">' +
        '<h3 id="formTitle" style="margin-top:0">Создание</h3>' +
        '<form id="adminForm" class="grid" style="grid-template-columns:1fr"></form>' +
      '</div>' +
    '</div>';

  const form = document.getElementById('adminForm');
  renderAdminForm(cfg, form, null, true);

  app.querySelectorAll('[data-action="edit"]').forEach((btn) => {
    btn.onclick = () => {
      const key = btn.getAttribute('data-key');
      const item = list.find((x) => String(x[cfg.key]) === String(key));
      renderAdminForm(cfg, form, item, false);
      document.getElementById('formTitle').textContent = 'Редактирование: ' + key;
    };
  });

  app.querySelectorAll('[data-action="delete"]').forEach((btn) => {
    btn.onclick = async () => {
      const key = btn.getAttribute('data-key');
      if (!window.confirm('Удалить ' + key + '?')) return;
      try {
        const path = cfg.remove.replace('{key}', encodeURIComponent(key));
        await apiRequest(path, 'DELETE');
        setStatus('Удалено: ' + key, true);
        invalidateCaches();
        router();
      } catch (err) {
        setStatus('Ошибка удаления: ' + (err.message || err), false);
      }
    };
  });

  if (cfg.restore) {
    app.querySelectorAll('[data-action="restore"]').forEach((btn) => {
      btn.onclick = async () => {
        const key = btn.getAttribute('data-key');
        try {
          const path = cfg.restore.replace('{key}', encodeURIComponent(key));
          await apiRequest(path, 'PUT');
          setStatus('Восстановлено: ' + key, true);
          invalidateCaches();
          router();
        } catch (err) {
          setStatus('Ошибка восстановления: ' + (err.message || err), false);
        }
      };
    });
  }
}

function renderAdminForm(cfg, formEl, item, isCreate) {
  const html = cfg.fields.map((f) => {
    const id = 'f_' + f.name;
    const value = item ? item[f.name] : '';
    let input;
    if (f.type === 'textarea') {
      input = '<textarea id="' + id + '">' + escapeHtml(value || '') + '</textarea>';
    } else if (f.type === 'select') {
      const opts = (f.options || []).map((o) => '<option value="' + attr(o) + (value === o ? '" selected' : '"') + '>' + escapeHtml(o) + '</option>').join('');
      input = '<select id="' + id + '">' + opts + '</select>';
    } else if (f.type === 'checkbox') {
      input = '<input id="' + id + '" type="checkbox"' + (value ? ' checked' : '') + '>';
    } else {
      const t = f.type || 'text';
      input = '<input id="' + id + '" class="field" type="' + t + '" value="' + attr(value == null ? '' : value) + '">';
    }
    return '<div class="form-row"><label>' + escapeHtml(f.label) + '</label>' + input + '</div>';
  }).join('');

  const buttons = isCreate
    ? '<button class="btn success" type="submit">Создать</button>'
    : '<button class="btn primary" type="submit">Сохранить</button>' +
      ' <button class="btn ghost" type="button" id="cancelEdit">Отмена</button>';

  formEl.innerHTML = html + '<div class="form-actions">' + buttons + '</div>';

  formEl.onsubmit = async (e) => {
    e.preventDefault();
    const payload = {};
    cfg.fields.forEach((f) => {
      const el = document.getElementById('f_' + f.name);
      if (!el) return;
      let v;
      if (f.type === 'checkbox') v = el.checked;
      else if (f.type === 'number') {
        const raw = el.value.trim();
        v = raw === '' ? null : Number(raw);
        if (v != null && Number.isNaN(v)) v = null;
      } else {
        v = el.value;
      }
      // turn comma-separated into array if name ends with _slugs/_codes
      if (typeof v === 'string' && /(_slugs|_codes)$/.test(f.name)) {
        const arr = v.split(',').map((x) => x.trim()).filter(Boolean);
        if (arr.length === 0) return;
        payload[f.name] = arr;
        return;
      }
      if (v === '' || v === null || v === undefined) return;
      payload[f.name] = v;
    });
    try {
      if (isCreate) {
        await apiRequest(cfg.create, 'POST', payload);
        setStatus('Создано', true);
      } else {
        const key = item[cfg.key];
        const path = cfg.update.replace('{key}', encodeURIComponent(key));
        await apiRequest(path, 'PUT', payload);
        setStatus('Сохранено: ' + key, true);
      }
      invalidateCaches();
      router();
    } catch (err) {
      setStatus('Ошибка: ' + (err.message || err), false);
    }
  };

  const cancel = document.getElementById('cancelEdit');
  if (cancel) cancel.onclick = () => {
    renderAdminForm(cfg, formEl, null, true);
    document.getElementById('formTitle').textContent = 'Создание';
  };
}

/* ============================ ROUTER ============================ */
async function router() {
  markMenu();
  try {
    const r = parseRoute();
    // close modal when navigating away from technology
    if (r.first !== 'technology' && modalEl.classList.contains('open')) {
      modalEl.classList.remove('open');
      modalWindow.innerHTML = '';
      document.body.style.overflow = '';
    }

    switch (r.first) {
      case '':
      case 'explore':
        return renderExplore();
      case 'catalog':
        return renderCatalog();
      case 'trends':
        return renderTrends();
      case 'tags':
        return renderTags();
      case 'sdgs':
        return renderSDGs();
      case 'organizations':
        return renderOrganizations();
      case 'metrics':
        return renderEntityListGeneric('Метрики', state.catalog.metrics.length ? state.catalog.metrics : await apiRequest('/api/metrics?locale=' + state.locale));
      case 'technology':
        return renderTechDetails(r.second);
      case 'trend':
        return renderTechByEntity('trend', r.second);
      case 'tag':
        return renderTechByEntity('tag', r.second);
      case 'sdg':
        return renderTechByEntity('sdg', r.second);
      case 'organization':
        return renderTechByEntity('organization', r.second);
      case 'admin':
        if (r.second === 'login' || r.second === '') return renderAdminLogin();
        return renderAdminEntity(r.second);
      default:
        return renderExplore();
    }
  } catch (err) {
    renderError('Не удалось загрузить страницу', err);
  }
}

function renderEntityListGeneric(title, list) {
  const items = Array.isArray(list) ? list : [];
  const cards = items.map((m) => '<div class="card"><h3>' + escapeHtml(m.name || m.title || m.slug || m.code) + '</h3>' +
    (m.type ? '<div class="tagline">type: ' + escapeHtml(m.type) + '</div>' : '') +
    (m.description ? '<p class="tagline">' + escapeHtml(m.description) + '</p>' : '') +
  '</div>').join('');
  app.innerHTML = '<div class="page-header"><h1>' + escapeHtml(title) + '</h1></div>' +
    '<div class="grid cards">' + (cards || '<div class="empty">Нет данных</div>') + '</div>';
}

window.addEventListener('hashchange', router);
window.addEventListener('DOMContentLoaded', router);
window.router = router;

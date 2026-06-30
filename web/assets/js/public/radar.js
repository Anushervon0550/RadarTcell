// ============ RADAR (explore view) ============
import { ns } from '../core/util.js';
import { S, colorForTrend } from './state.js';
import { openTechDetail } from './detail.js';

const R_OUTER = 380, R_INNER = 60;

function ringRadius(trl) {
  const t = Math.max(1, Math.min(9, trl || 1));
  return R_INNER + ((9 - t) / 8) * (R_OUTER - R_INNER);
}

export function renderRadar() {
  const svg = document.getElementById('rsvg');
  if (!svg) return;
  svg.innerHTML = '';

  for (let trl = 1; trl <= 9; trl++) {
    const r = ringRadius(trl);
    const c = ns('circle');
    c.setAttribute('cx', 0); c.setAttribute('cy', 0); c.setAttribute('r', r);
    c.setAttribute('class', 'ring' + ([3, 6, 9].includes(trl) ? ' major' : ''));
    svg.appendChild(c);
  }

  [
    { trl: 9, l: 'TRL 9 — Внедрено' },
    { trl: 6, l: 'TRL 6 — Прототип' },
    { trl: 3, l: 'TRL 3 — Концепция' },
    { trl: 1, l: 'TRL 1 — Исследование' },
  ].forEach((rd) => {
    const r = ringRadius(rd.trl);
    const t = ns('text');
    t.setAttribute('x', 0); t.setAttribute('y', -r - 4);
    t.setAttribute('text-anchor', 'middle'); t.setAttribute('class', 'ring-label');
    t.textContent = rd.l;
    svg.appendChild(t);
  });

  const tg = new Map();
  S.techs.forEach((t) => {
    if (!tg.has(t.trend_slug)) tg.set(t.trend_slug, { name: t.trend_name || t.trend_slug, slug: t.trend_slug, items: [] });
    tg.get(t.trend_slug).items.push(t);
  });

  if (tg.size > 1) {
    tg.forEach((g) => {
      let sx = 0, sy = 0;
      g.items.forEach((t) => { sx += Math.cos(t.angle); sy += Math.sin(t.angle); });
      g.angle = Math.atan2(sy, sx);
    });
    const sorted = [...tg.values()].sort((a, b) => a.angle - b.angle);
    for (let i = 0; i < sorted.length; i++) {
      const a = sorted[i].angle, b = sorted[(i + 1) % sorted.length].angle;
      let mid = (a + b) / 2;
      if (b < a) mid = (a + b + Math.PI * 2) / 2;
      const x = Math.cos(mid) * R_OUTER, y = Math.sin(mid) * R_OUTER;
      const ln = ns('line');
      ln.setAttribute('x1', 0); ln.setAttribute('y1', 0);
      ln.setAttribute('x2', x); ln.setAttribute('y2', y);
      ln.setAttribute('class', 'sector-line');
      svg.appendChild(ln);
    }
    sorted.forEach((g) => {
      const r = R_OUTER + 22, x = Math.cos(g.angle) * r, y = Math.sin(g.angle) * r;
      const t = ns('text');
      t.setAttribute('x', x); t.setAttribute('y', y);
      t.setAttribute('text-anchor', 'middle'); t.setAttribute('dominant-baseline', 'middle');
      t.setAttribute('class', 'trend-label'); t.setAttribute('fill', colorForTrend(g.slug));
      t.textContent = (g.name || '').slice(0, 22);
      svg.appendChild(t);
    });
  }

  S.techs.forEach((t) => {
    const r = ringRadius(t.trl), x = Math.cos(t.angle) * r, y = Math.sin(t.angle) * r;
    const g = ns('g');
    g.setAttribute('class', 'dot' + (S.selSlug === t.slug ? ' sel' : ''));
    g.setAttribute('transform', `translate(${x.toFixed(2)},${y.toFixed(2)})`);
    g.dataset.slug = t.slug;
    const c = ns('circle');
    c.setAttribute('r', 5); c.setAttribute('fill', colorForTrend(t.trend_slug));
    g.appendChild(c);
    if (t.index) {
      const tx = ns('text');
      tx.setAttribute('y', -9); tx.setAttribute('text-anchor', 'middle');
      tx.setAttribute('font-family', 'DM Mono'); tx.setAttribute('font-size', '9');
      tx.setAttribute('fill', '#b9a8d9'); tx.setAttribute('opacity', '.85');
      tx.textContent = t.index;
      g.appendChild(tx);
    }
    g.addEventListener('mouseenter', (e) => showTip(e, t));
    g.addEventListener('mousemove', moveTip);
    g.addEventListener('mouseleave', hideTip);
    g.addEventListener('click', () => openTechDetail(t.slug));
    svg.appendChild(g);
  });
}

/* ---- tooltip ---- */
function showTip(e, t) {
  const tt = document.getElementById('tt');
  tt.textContent = `${t.name} · TRL ${t.trl}`;
  moveTip(e);
  tt.classList.add('on');
}
function moveTip(e) {
  const tt = document.getElementById('tt');
  tt.style.left = e.clientX + 'px';
  tt.style.top = e.clientY + 'px';
}
function hideTip() { document.getElementById('tt').classList.remove('on'); }

/* ---- zoom / pan ---- */
export function applyTransform() {
  const svgEl = document.getElementById('rsvg');
  svgEl.style.transform = `translate(${S.pan.x}px,${S.pan.y}px) scale(${S.zoom})`;
  document.getElementById('zv').textContent = Math.round(S.zoom * 100) + '%';
}

export function initRadarControls() {
  const wrap = document.getElementById('rwrap');

  document.getElementById('zi').onclick = () => { S.zoom = Math.min(4, S.zoom * 1.25); applyTransform(); };
  document.getElementById('zo').onclick = () => { S.zoom = Math.max(0.4, S.zoom / 1.25); applyTransform(); };
  document.getElementById('zr').onclick = () => { S.zoom = 1; S.pan = { x: 0, y: 0 }; applyTransform(); };
  document.getElementById('zf').onclick = () => {
    if (document.fullscreenElement) document.exitFullscreen();
    else document.documentElement.requestFullscreen();
  };

  wrap.addEventListener('wheel', (e) => {
    if (S.view !== 'explore') return;
    e.preventDefault();
    const d = e.deltaY < 0 ? 1.1 : 1 / 1.1;
    S.zoom = Math.max(0.4, Math.min(4, S.zoom * d));
    applyTransform();
  }, { passive: false });

  let dragging = false, dx = 0, dy = 0, opx = 0, opy = 0;
  wrap.addEventListener('mousedown', (e) => {
    if (e.target.closest('.dot')) return;
    dragging = true; dx = e.clientX; dy = e.clientY; opx = S.pan.x; opy = S.pan.y;
    wrap.classList.add('dragging');
  });
  window.addEventListener('mousemove', (e) => {
    if (!dragging) return;
    S.pan.x = opx + (e.clientX - dx);
    S.pan.y = opy + (e.clientY - dy);
    applyTransform();
  });
  window.addEventListener('mouseup', () => { dragging = false; wrap.classList.remove('dragging'); });

  applyTransform();
}

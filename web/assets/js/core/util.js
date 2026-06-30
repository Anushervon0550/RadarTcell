// ============ SHARED UTILITIES ============

/** DOM-селекторы. */
export const $ = (sel, root = document) => root.querySelector(sel);
export const $$ = (sel, root = document) => Array.from(root.querySelectorAll(sel));

/** Создать SVG-элемент в нужном namespace. */
export const ns = (name) => document.createElementNS('http://www.w3.org/2000/svg', name);

/** Экранирование для безопасной вставки в HTML. */
export function esc(s) {
  return String(s == null ? '' : s).replace(/[&<>"]/g, (c) => ({
    '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;',
  }[c]));
}

/** Экранирование для значения внутри одинарных кавычек атрибута. */
export function escAt(s) { return esc(s).replace(/'/g, '&#39;'); }

/** Тост-уведомление. Требует элемент #toast в DOM. */
let toastTimer;
export function toast(msg, type = '') {
  const el = document.getElementById('toast');
  if (!el) return;
  el.className = 'toast on ' + type;
  el.textContent = msg;
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => el.classList.remove('on'), 2400);
}

/** debounce — откладывает вызов до паузы в событиях. */
export function debounce(fn, wait = 250) {
  let t;
  return (...args) => {
    clearTimeout(t);
    t = setTimeout(() => fn(...args), wait);
  };
}

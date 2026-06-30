// ============ ADMIN CRUD ============
import { apiGet, apiPost, apiPut, apiDel } from '../core/api.js';
import { $$, esc, toast } from '../core/util.js';
import { AS, ENT } from './state.js';

const EDIT_ICON = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M12 20h9M16.5 3.5a2.1 2.1 0 1 1 3 3L7 19l-4 1 1-4 12.5-12.5z"/></svg>';
const DEL_ICON = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M3 6h18M8 6V4h8v2M19 6l-1 14H6L5 6M10 11v6M14 11v6"/></svg>';

export async function adminLoad() {
  const ent = ENT[AS.entity];
  const table = document.getElementById('aTable');
  try {
    const r = await apiGet(ent.list);
    const items = Array.isArray(r) ? r : (r.items || r.data || []);
    AS.items = items;
    document.getElementById('aMeta').textContent = items.length + ' записей';

    const head = '<tr>' + ent.cols.map((c) => `<th>${c[1]}</th>`).join('') + '<th></th></tr>';
    const rows = items.map((it) => '<tr>' + ent.cols.map((c) => `<td>${esc(it[c[0]] ?? '')}</td>`).join('')
      + `<td class="act">${ent.noEdit ? '' : `<button class="iconbtn" data-edit="${esc(it[ent.idKey])}" aria-label="Редактировать">${EDIT_ICON}</button>`}`
      + `<button class="iconbtn danger" data-del="${esc(it[ent.idKey])}" aria-label="Удалить">${DEL_ICON}</button></td></tr>`).join('');

    table.innerHTML = `<table class="tbl"><thead>${head}</thead><tbody>${rows
      || `<tr><td colspan="${ent.cols.length + 1}" style="color:var(--t3);text-align:center;padding:30px">Пусто</td></tr>`}</tbody></table>`;

    $$('[data-edit]').forEach((b) => { b.onclick = () => openForm(b.dataset.edit); });
    $$('[data-del]').forEach((b) => { b.onclick = () => delEntity(b.dataset.del); });
  } catch (e) {
    table.innerHTML = '<div style="padding:40px;color:var(--danger)">' + esc(e.message) + '</div>';
  }
}

export function openForm(id) {
  const ent = ENT[AS.entity];
  const item = id ? AS.items.find((x) => String(x[ent.idKey]) === String(id)) : null;
  document.getElementById('mfTitle').textContent = item ? `Редактирование · ${ent.title}` : `Новая запись · ${ent.title}`;
  document.getElementById('mfSub').textContent = item ? String(item[ent.idKey]) : '';
  const f = document.getElementById('mfForm');
  f.innerHTML = ent.fields.map((fd) => {
    const v = item ? (item[fd.k] ?? '') : '';
    const lbl = esc(fd.label || fd.k);
    if (fd.type === 'textarea') return `<div class="field"><label>${lbl}</label><textarea name="${fd.k}">${esc(v)}</textarea></div>`;
    if (fd.type === 'bool') return `<div class="field"><label>${lbl}</label><select name="${fd.k}"><option value="true" ${v ? 'selected' : ''}>true</option><option value="false" ${!v ? 'selected' : ''}>false</option></select></div>`;
    if (fd.type === 'number') return `<div class="field"><label>${lbl}</label><input type="number" step="any" name="${fd.k}" value="${esc(v)}"/></div>`;
    if (fd.type === 'password') return `<div class="field"><label>${lbl}</label><input type="password" name="${fd.k}"/></div>`;
    if (fd.type === 'csv') { const arr = Array.isArray(v) ? v.join(',') : (v || ''); return `<div class="field"><label>${lbl}</label><input name="${fd.k}" value="${esc(arr)}"/></div>`; }
    return `<div class="field"><label>${lbl} ${fd.req ? '<span style="color:var(--accent)">*</span>' : ''}</label><input name="${fd.k}" value="${esc(v)}"/></div>`;
  }).join('');
  document.getElementById('mfErr').style.display = 'none';
  f.dataset.editId = item ? item[ent.idKey] : '';
  openModal('mForm');
}

export async function saveForm() {
  const ent = ENT[AS.entity];
  const f = document.getElementById('mfForm');
  const err = document.getElementById('mfErr');
  const id = f.dataset.editId;
  const data = {};
  ent.fields.forEach((fd) => {
    const el = f.elements[fd.k];
    if (!el) return;
    let v = el.value;
    if (fd.type === 'number') v = v === '' ? null : Number(v);
    else if (fd.type === 'bool') v = v === 'true';
    else if (fd.type === 'csv') v = v ? v.split(',').map((s) => s.trim()).filter(Boolean) : [];
    if (v !== '' && v !== null) data[fd.k] = v;
    else if (fd.type === 'number') data[fd.k] = null;
  });
  try {
    if (id) await apiPut(`/admin/${AS.entity}/${encodeURIComponent(id)}`, data);
    else await apiPost(`/admin/${AS.entity}`, data);
    closeModal('mForm');
    toast('Сохранено', 'ok');
    adminLoad();
  } catch (e) {
    err.textContent = e.message || 'Ошибка';
    err.style.display = 'block';
  }
}

export async function delEntity(id) {
  if (!confirm('Удалить запись «' + id + '»?')) return;
  try {
    await apiDel(`/admin/${AS.entity}/${encodeURIComponent(id)}`);
    toast('Удалено', 'ok');
    adminLoad();
  } catch (e) {
    toast(e.message || 'Ошибка', 'err');
  }
}

/* ---- modal helpers ---- */
export function openModal(mid) { document.getElementById(mid).classList.add('on'); }
export function closeModal(mid) { document.getElementById(mid).classList.remove('on'); }

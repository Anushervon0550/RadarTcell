// ============ ADMIN APP ENTRY ============
import { $$, toast } from '../core/util.js';
import { AS } from './state.js';
import { refreshMe, login, logout } from './auth.js';
import { adminLoad, openForm, saveForm, closeModal } from './crud.js';

const gate = document.getElementById('gate');
const shell = document.getElementById('admin');

function showGate() {
  gate.style.display = 'flex';
  shell.style.display = 'none';
}
function showShell() {
  gate.style.display = 'none';
  shell.style.display = 'flex';
  document.getElementById('mePill').textContent = '👤 ' + (AS.me?.username || '');
  adminLoad();
}

/* ---- login gate ---- */
document.getElementById('liGo').onclick = doLogin;
document.getElementById('liPass').addEventListener('keydown', (e) => { if (e.key === 'Enter') doLogin(); });

async function doLogin() {
  const u = document.getElementById('liUser').value.trim();
  const p = document.getElementById('liPass').value;
  const err = document.getElementById('liErr');
  err.style.display = 'none';
  try {
    await login(u, p);
    showShell();
    toast('Вход выполнен', 'ok');
  } catch (e) {
    err.textContent = e.message || 'Ошибка входа';
    err.style.display = 'block';
  }
}

/* ---- logout ---- */
document.getElementById('mePill').onclick = () => {
  if (confirm('Выйти из админки?')) { logout(); showGate(); }
};

/* ---- session expired (401) ---- */
window.addEventListener('auth:expired', () => {
  logout();
  showGate();
  toast('Сессия истекла, войдите заново', 'err');
});

/* ---- entity switch + create ---- */
document.getElementById('aEntity').onchange = (e) => { AS.entity = e.target.value; adminLoad(); };
document.getElementById('aNew').onclick = () => openForm(null);

/* ---- modal wiring ---- */
document.getElementById('mfSave').onclick = saveForm;
$$('[data-mc]').forEach((b) => { b.onclick = () => closeModal(b.dataset.mc); });
$$('.modal').forEach((m) => m.addEventListener('click', (e) => { if (e.target === m) m.classList.remove('on'); }));
document.addEventListener('keydown', (e) => {
  if (e.key === 'Escape') $$('.modal.on').forEach((m) => m.classList.remove('on'));
});

/* ---- init ---- */
(async function init() {
  const ok = await refreshMe();
  if (ok) showShell(); else showGate();
}());

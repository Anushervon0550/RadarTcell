// ============ API CLIENT ============
// Тонкая обёртка над fetch: авторизация по JWT, JSON-сериализация,
// единая обработка ошибок и глобальная реакция на 401 (протухший токен).

const API = '/api';
const TOKEN_KEY = 'jwt';

export function getToken() { return localStorage.getItem(TOKEN_KEY) || ''; }
export function setToken(t) { localStorage.setItem(TOKEN_KEY, t || ''); }
export function clearToken() { localStorage.removeItem(TOKEN_KEY); }

function authHeader() {
  const t = getToken();
  return t ? { Authorization: 'Bearer ' + t } : {};
}

export async function api(path, opt = {}) {
  const headers = { Accept: 'application/json', ...(opt.headers || {}), ...authHeader() };
  if (opt.body && !(opt.body instanceof FormData)) {
    headers['Content-Type'] = 'application/json';
    opt = { ...opt, body: JSON.stringify(opt.body) };
  }
  const res = await fetch(API + path, { ...opt, headers });
  const text = await res.text();
  let data;
  try { data = text ? JSON.parse(text) : null; } catch { data = text; }

  if (!res.ok) {
    if (res.status === 401) {
      clearToken();
      window.dispatchEvent(new CustomEvent('auth:expired'));
    }
    const err = new Error((data && data.error) || res.statusText);
    err.status = res.status;
    err.data = data;
    throw err;
  }
  return data;
}

export const apiGet = (p) => api(p);
export const apiPost = (p, b) => api(p, { method: 'POST', body: b });
export const apiPut = (p, b) => api(p, { method: 'PUT', body: b });
export const apiDel = (p) => api(p, { method: 'DELETE' });

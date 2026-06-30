// ============ ADMIN AUTH ============
import { apiGet, apiPost, getToken, setToken, clearToken } from '../core/api.js';
import { AS } from './state.js';

/** Проверяет токен и подгружает профиль. Возвращает true, если авторизован. */
export async function refreshMe() {
  if (!getToken()) { AS.me = null; return false; }
  try {
    AS.me = await apiGet('/admin/me');
    return true;
  } catch {
    clearToken();
    AS.me = null;
    return false;
  }
}

export async function login(username, password) {
  const r = await apiPost('/admin/login', { username, password });
  setToken(r.token || r.access_token || '');
  AS.me = await apiGet('/admin/me');
  return AS.me;
}

export function logout() {
  clearToken();
  AS.me = null;
}

/** Env exposed to the browser. Only VITE_* prefixed values are bundled. */
interface AppEnv {
  API_BASE: string;
  ADMIN_URL: string;
  DEFAULT_LOCALE: string;
  ENABLE_SWAGGER: boolean;
}

const raw = import.meta.env;

export const env: AppEnv = {
  API_BASE: (raw.VITE_API_BASE ?? '').replace(/\/$/, '') || window.location.origin,
  ADMIN_URL: raw.VITE_ADMIN_URL ?? '/admin',
  DEFAULT_LOCALE: raw.VITE_DEFAULT_LOCALE ?? 'ru',
  ENABLE_SWAGGER: raw.VITE_ENABLE_SWAGGER === 'true',
};

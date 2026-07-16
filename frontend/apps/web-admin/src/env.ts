interface AppEnv {
  API_BASE: string;
  PUBLIC_URL: string;
  DEFAULT_LOCALE: string;
}

const raw = import.meta.env;

export const env: AppEnv = {
  API_BASE: (raw.VITE_API_BASE ?? '').replace(/\/$/, '') || window.location.origin,
  PUBLIC_URL: raw.VITE_PUBLIC_URL ?? '/',
  DEFAULT_LOCALE: raw.VITE_DEFAULT_LOCALE ?? 'ru',
};

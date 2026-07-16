import { ApiClient } from '@radartcell/api';
import { env } from '@/env';
import { useAuthStore } from '@/auth/store';

export const api = new ApiClient({
  baseUrl: env.API_BASE,
  defaultLocale: env.DEFAULT_LOCALE,
  getToken: () => useAuthStore.getState().token,
  onUnauthorized: () => {
    useAuthStore.getState().clear();
    if (!location.pathname.startsWith('/login')) {
      const next = encodeURIComponent(location.pathname + location.search);
      location.assign(`/login?next=${next}`);
    }
  },
});

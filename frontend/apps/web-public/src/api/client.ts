import { ApiClient } from '@radartcell/api';
import { env } from '@/env';

export const api = new ApiClient({
  baseUrl: env.API_BASE,
  defaultLocale: env.DEFAULT_LOCALE,
});

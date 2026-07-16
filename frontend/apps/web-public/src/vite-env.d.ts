/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE?: string;
  readonly VITE_ADMIN_URL?: string;
  readonly VITE_DEFAULT_LOCALE?: string;
  readonly VITE_ENABLE_SWAGGER?: string;
  readonly VITE_API_PROXY_TARGET?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}

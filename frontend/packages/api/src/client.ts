import type { ApiErrorResponse } from './types';

export class ApiError extends Error {
  status: number;
  code?: string;
  details?: unknown;

  constructor(message: string, status: number, code?: string, details?: unknown) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.code = code;
    this.details = details;
  }

  get isUnauthorized(): boolean {
    return this.status === 401;
  }

  get isForbidden(): boolean {
    return this.status === 403;
  }

  get isNotFound(): boolean {
    return this.status === 404;
  }
}

export interface ApiClientOptions {
  baseUrl: string;
  getToken?: () => string | null | undefined;
  onUnauthorized?: () => void;
  fetchImpl?: typeof fetch;
  defaultLocale?: string;
}

export interface RequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
  body?: unknown;
  headers?: Record<string, string>;
  signal?: AbortSignal;
  query?: Record<string, string | number | boolean | null | undefined | string[]>;
}

/**
 * Serialize query params. Arrays produce repeated keys (`highlight=a&highlight=b`),
 * which matches Go chi router expectations.
 */
export function buildQuery(query: RequestOptions['query']): string {
  if (!query) return '';
  const usp = new URLSearchParams();
  for (const [key, value] of Object.entries(query)) {
    if (value === undefined || value === null || value === '') continue;
    if (Array.isArray(value)) {
      for (const v of value) {
        if (v !== undefined && v !== null && v !== '') usp.append(key, String(v));
      }
    } else {
      usp.set(key, String(value));
    }
  }
  const s = usp.toString();
  return s ? `?${s}` : '';
}

export class ApiClient {
  readonly baseUrl: string;
  private readonly getToken?: () => string | null | undefined;
  private readonly onUnauthorized?: () => void;
  private readonly fetchImpl: typeof fetch;
  readonly defaultLocale: string;

  constructor(opts: ApiClientOptions) {
    this.baseUrl = opts.baseUrl.replace(/\/$/, '');
    this.getToken = opts.getToken;
    this.onUnauthorized = opts.onUnauthorized;
    this.fetchImpl = opts.fetchImpl ?? globalThis.fetch.bind(globalThis);
    this.defaultLocale = opts.defaultLocale ?? 'ru';
  }

  async request<T>(path: string, options: RequestOptions = {}): Promise<T> {
    const { method = 'GET', body, headers = {}, signal, query } = options;

    const finalHeaders: Record<string, string> = {
      Accept: 'application/json',
      ...headers,
    };

    if (body !== undefined) {
      finalHeaders['Content-Type'] = 'application/json';
    }

    const token = this.getToken?.();
    if (token) {
      finalHeaders.Authorization = `Bearer ${token}`;
    }

    const url = this.baseUrl + path + buildQuery(query);

    let response: Response;
    try {
      response = await this.fetchImpl(url, {
        method,
        headers: finalHeaders,
        body: body !== undefined ? JSON.stringify(body) : undefined,
        signal,
        credentials: 'same-origin',
      });
    } catch (err) {
      if ((err as Error).name === 'AbortError') throw err;
      throw new ApiError('Network error', 0);
    }

    if (response.status === 204) return null as T;

    const text = await response.text();
    let data: unknown = null;
    if (text) {
      try {
        data = JSON.parse(text);
      } catch {
        data = { error: text };
      }
    }

    if (!response.ok) {
      if (response.status === 401) this.onUnauthorized?.();
      const err = (data ?? {}) as Partial<ApiErrorResponse>;
      throw new ApiError(
        err.error ?? `HTTP ${response.status}`,
        response.status,
        err.code,
        err.details,
      );
    }

    return data as T;
  }

  get<T>(path: string, query?: RequestOptions['query'], signal?: AbortSignal): Promise<T> {
    return this.request<T>(path, { method: 'GET', query, signal });
  }

  post<T>(path: string, body?: unknown, query?: RequestOptions['query']): Promise<T> {
    return this.request<T>(path, { method: 'POST', body, query });
  }

  put<T>(path: string, body?: unknown, query?: RequestOptions['query']): Promise<T> {
    return this.request<T>(path, { method: 'PUT', body, query });
  }

  patch<T>(path: string, body?: unknown, query?: RequestOptions['query']): Promise<T> {
    return this.request<T>(path, { method: 'PATCH', body, query });
  }

  delete<T>(path: string, query?: RequestOptions['query']): Promise<T> {
    return this.request<T>(path, { method: 'DELETE', query });
  }
}

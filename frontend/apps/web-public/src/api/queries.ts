import { useQuery } from '@tanstack/react-query';
import type {
  HomeResponse,
  MetricDefinition,
  Organization,
  SDG,
  Tag,
  TechnologyCard,
  TechnologyListParams,
  TechnologyListResponse,
  Trend,
} from '@radartcell/api';
import { api } from './client';
import { env } from '@/env';

const withLocale = (params: Record<string, unknown> = {}) => ({
  locale: env.DEFAULT_LOCALE,
  ...params,
});

export const catalogKeys = {
  trends: ['catalog', 'trends'] as const,
  tags: ['catalog', 'tags'] as const,
  sdgs: ['catalog', 'sdgs'] as const,
  organizations: ['catalog', 'organizations'] as const,
  metrics: ['catalog', 'metrics'] as const,
};

export function useTrends() {
  return useQuery({
    queryKey: catalogKeys.trends,
    queryFn: () => api.get<Trend[]>('/api/trends', withLocale()),
    staleTime: 5 * 60_000,
  });
}

export function useTags() {
  return useQuery({
    queryKey: catalogKeys.tags,
    queryFn: () => api.get<Tag[]>('/api/tags', withLocale()),
    staleTime: 5 * 60_000,
  });
}

export function useSDGs() {
  return useQuery({
    queryKey: catalogKeys.sdgs,
    queryFn: () => api.get<SDG[]>('/api/sdgs', withLocale()),
    staleTime: 5 * 60_000,
  });
}

export function useOrganizations() {
  return useQuery({
    queryKey: catalogKeys.organizations,
    queryFn: () => api.get<Organization[]>('/api/organizations', withLocale()),
    staleTime: 5 * 60_000,
  });
}

export function useMetrics() {
  return useQuery({
    queryKey: catalogKeys.metrics,
    queryFn: () => api.get<MetricDefinition[]>('/api/metrics', withLocale()),
    staleTime: 5 * 60_000,
  });
}

export function useHome(limit = 200) {
  return useQuery({
    queryKey: ['home', limit] as const,
    queryFn: () => api.get<HomeResponse>('/api/home', withLocale({ limit })),
    staleTime: 60_000,
  });
}

export function useTechnologies(params: TechnologyListParams) {
  return useQuery({
    queryKey: ['technologies', params] as const,
    queryFn: () =>
      api.get<TechnologyListResponse>(
        '/api/technologies',
        withLocale(params as Record<string, unknown>),
      ),
    staleTime: 30_000,
  });
}

export function useTechnology(slug: string | undefined) {
  return useQuery({
    queryKey: ['technology', slug] as const,
    queryFn: () =>
      api.get<TechnologyCard>(`/api/technologies/${encodeURIComponent(slug!)}`, withLocale()),
    enabled: !!slug,
    staleTime: 60_000,
  });
}

export function useTechnologiesByEntity(
  kind: 'trend' | 'tag' | 'sdg' | 'organization',
  value: string | undefined,
  limit = 200,
) {
  const paths: Record<typeof kind, (v: string) => string> = {
    trend: (v) => `/api/trends/${encodeURIComponent(v)}/technologies`,
    tag: (v) => `/api/tags/${encodeURIComponent(v)}/technologies`,
    sdg: (v) => `/api/sdgs/${encodeURIComponent(v)}/technologies`,
    organization: (v) => `/api/organizations/${encodeURIComponent(v)}/technologies`,
  };

  return useQuery({
    queryKey: ['technologies-by-entity', kind, value, limit] as const,
    queryFn: () =>
      api.get<TechnologyListResponse>(paths[kind](value!), withLocale({ limit })),
    enabled: !!value,
    staleTime: 60_000,
  });
}

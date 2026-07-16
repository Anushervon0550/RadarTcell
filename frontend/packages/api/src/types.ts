/**
 * Shared API types for RadarTcell.
 * Mirrors Go domain types from internal/domain and internal/httpapi.
 * Keep in sync with docs/openapi.yaml.
 */

// ============================== CATALOG =====================================

export interface Trend {
  id: string;
  slug: string;
  name: string;
  technologies_count: number;
}

export interface SDG {
  id: string;
  code: string;
  title: string;
  technologies_count: number;
}

export interface Tag {
  id: string;
  slug: string;
  title: string;
  category?: string | null;
  description?: string | null;
}

export interface Organization {
  id: string;
  slug: string;
  name: string;
  logo_url?: string | null;
  description?: string | null;
  website?: string | null;
  headquarters?: string | null;
  technologies_count: number;
}

export type MetricType = 'distance' | 'bubble' | 'bar';

export interface MetricDefinition {
  id: string;
  name: string;
  type: MetricType;
  description?: string | null;
  orderable: boolean;
  field_key?: string | null;
}

// ============================ TECHNOLOGIES ==================================

export type Stage = 'idea' | 'prototype' | 'product';

export interface TechnologyMetricValue {
  metric_id: string;
  field_key?: string | null;
  value?: number | null;
}

export interface TechnologyBase {
  id: string;
  slug: string;
  index: number;
  name: string;
  description_short?: string | null;
  trl: number;
  trend_id: string;
  trend_slug: string;
  trend_name: string;
  custom_metric_1?: number | null;
  custom_metric_2?: number | null;
  custom_metric_3?: number | null;
  custom_metric_4?: number | null;
  image_url?: string | null;
  angle: number;
  radius: number;
}

export interface TechnologyListItem extends TechnologyBase {
  custom_metrics?: TechnologyMetricValue[];
  custom_metric_1_norm: number;
  custom_metric_2_norm: number;
  custom_metric_3_norm: number;
  custom_metric_4_norm: number;
}

export interface TechnologyListResponse {
  page: number;
  limit: number;
  total: number;
  next_cursor?: string;
  items: TechnologyListItem[];
}

export interface TechnologyCard extends TechnologyBase {
  custom_metrics?: TechnologyMetricValue[];
  description_full?: string | null;
  source_link?: string | null;
  tags: Tag[];
  sdgs: SDG[];
  organizations: Organization[];
}

// =============================== HOME =======================================

export interface HomeTechItem {
  id: string;
  slug: string;
  name: string;
  description_short?: string | null;
  image_url?: string | null;
  trl: number;
  stage: Stage;
  completion: number;
  angle: number;
  radius: number;
}

export interface HomeTrendBlock {
  slug: string;
  name: string;
  items: HomeTechItem[];
}

export interface HomeResponse {
  page: number;
  limit: number;
  total: number;
  locale?: string;
  trends: HomeTrendBlock[];
}

// =============================== ADMIN ======================================

export interface AdminLoginResponse {
  token: string;
}

export interface AdminMeResponse {
  user: string;
  role: string;
}

export interface AdminUser {
  username: string;
  is_active: boolean;
  created_at?: string;
}

export interface AdminTrend {
  id: string;
  slug: string;
  name: string;
  description?: string | null;
  image_url?: string | null;
  order_index: number;
}

export interface AdminTag {
  id: string;
  slug: string;
  title: string;
  category?: string | null;
  description?: string | null;
}

export interface AdminOrganization extends Organization {}

export interface AdminMetric {
  id: string;
  name: string;
  type: MetricType;
  field_key?: string | null;
  orderable: boolean;
  description?: string | null;
}

export interface AdminSDG {
  id?: string;
  code: string;
  title: string;
  icon?: string | null;
  description?: string | null;
}

export interface AdminTechnology {
  id: string;
  slug: string;
  index: number;
  name: string;
  trend_slug: string;
  trend_name?: string;
  trl: number;
  source_link?: string | null;
  image_url?: string | null;
  description_short?: string | null;
  description_full?: string | null;
  tag_slugs?: string[];
  sdg_codes?: string[];
  organization_slugs?: string[];
  custom_metric_1?: number | null;
  custom_metric_2?: number | null;
  custom_metric_3?: number | null;
  custom_metric_4?: number | null;
  deleted_at?: string | null;
}

export interface AdminTechnologyListResponse {
  items: AdminTechnology[];
  total: number;
  page: number;
  limit: number;
}

// ================================ COMMON ====================================

export interface ApiErrorResponse {
  error: string;
  code?: string;
  details?: unknown;
}

export interface TechnologyListParams {
  page?: number;
  limit?: number;
  cursor?: string;
  search?: string;
  trend_id?: string;
  sdg_id?: string;
  tag_id?: string;
  organization_id?: string;
  trl_min?: number;
  trl_max?: number;
  sort_by?: string;
  order?: 'asc' | 'desc';
  highlight?: string[];
  locale?: string;
}

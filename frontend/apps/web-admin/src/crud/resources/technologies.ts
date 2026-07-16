import { z } from 'zod';
import type { AdminTechnology } from '@radartcell/api';
import type { CrudResource } from '../types';
import { joinCSV, slugSchema, splitCSV } from './shared';

const schema = z.object({
  slug: slugSchema,
  index: z.coerce.number().int().min(0).max(999999),
  name: z.string().trim().min(1).max(300),
  trend_slug: slugSchema,
  trl: z.coerce.number().int().min(1).max(9),
  source_link: z.string().trim().max(1000).optional().nullable().or(z.literal('')),
  image_url: z.string().trim().max(1000).optional().nullable().or(z.literal('')),
  description_short: z.string().optional().nullable(),
  description_full: z.string().optional().nullable(),
  tag_slugs: z.string().optional().nullable(),
  sdg_codes: z.string().optional().nullable(),
  organization_slugs: z.string().optional().nullable(),
  custom_metric_1: z.union([z.coerce.number(), z.literal('')]).optional(),
  custom_metric_2: z.union([z.coerce.number(), z.literal('')]).optional(),
  custom_metric_3: z.union([z.coerce.number(), z.literal('')]).optional(),
  custom_metric_4: z.union([z.coerce.number(), z.literal('')]).optional(),
});

export const technologiesResource: CrudResource<AdminTechnology & Record<string, unknown>> = {
  name: 'technologies',
  title: 'Технологии',
  keyField: 'slug',
  listPath: '/api/admin/technologies?include_deleted=true&limit=200',
  createPath: '/api/admin/technologies',
  updatePath: (k) => `/api/admin/technologies/${encodeURIComponent(k)}`,
  removePath: (k) => `/api/admin/technologies/${encodeURIComponent(k)}`,
  restorePath: (k) => `/api/admin/technologies/${encodeURIComponent(k)}/restore`,
  softDelete: true,
  columns: [
    { key: 'slug', header: 'Slug' },
    { key: 'name', header: 'Название' },
    { key: 'trend_slug', header: 'Тренд' },
    { key: 'trl', header: 'TRL' },
    { key: 'deleted_at', header: 'Удалена', render: (r) => (r.deleted_at ? 'да' : '—') },
  ],
  fields: [
    { name: 'slug', label: 'Slug', type: 'text', requiredOnCreate: true },
    { name: 'index', label: 'Индекс', type: 'number', min: 0, required: true },
    { name: 'name', label: 'Название', type: 'text', required: true },
    { name: 'trend_slug', label: 'Slug тренда', type: 'text', required: true },
    { name: 'trl', label: 'TRL', type: 'number', min: 1, max: 9, required: true },
    { name: 'source_link', label: 'Источник', type: 'url' },
    { name: 'image_url', label: 'Обложка (URL)', type: 'url' },
    { name: 'description_short', label: 'Краткое описание', type: 'textarea' },
    { name: 'description_full', label: 'Полное описание', type: 'textarea' },
    { name: 'tag_slugs', label: 'Теги', type: 'multi', hint: 'slug через запятую' },
    { name: 'sdg_codes', label: 'ЦУР', type: 'multi', hint: 'коды через запятую (SDG 09, ...)' },
    { name: 'organization_slugs', label: 'Организации', type: 'multi', hint: 'slug через запятую' },
    { name: 'custom_metric_1', label: 'Метрика 1 (0..1)', type: 'number', step: 0.01 },
    { name: 'custom_metric_2', label: 'Метрика 2 (0..1)', type: 'number', step: 0.01 },
    { name: 'custom_metric_3', label: 'Метрика 3 (0..1)', type: 'number', step: 0.01 },
    { name: 'custom_metric_4', label: 'Метрика 4 (0..1)', type: 'number', step: 0.01 },
  ],
  schema,
  parseList: (data) => {
    if (Array.isArray(data)) return data as (AdminTechnology & Record<string, unknown>)[];
    if (data && typeof data === 'object' && Array.isArray((data as { items?: unknown }).items)) {
      return (data as { items: (AdminTechnology & Record<string, unknown>)[] }).items;
    }
    return [];
  },
  serialize: (values) => ({
    ...values,
    tag_slugs: splitCSV(values.tag_slugs),
    sdg_codes: splitCSV(values.sdg_codes),
    organization_slugs: splitCSV(values.organization_slugs),
    // Ensure empty numbers become null (bypass empty-string coercion)
    custom_metric_1: values.custom_metric_1 === '' ? null : values.custom_metric_1,
    custom_metric_2: values.custom_metric_2 === '' ? null : values.custom_metric_2,
    custom_metric_3: values.custom_metric_3 === '' ? null : values.custom_metric_3,
    custom_metric_4: values.custom_metric_4 === '' ? null : values.custom_metric_4,
  }),
  deserialize: (row) => ({
    ...row,
    tag_slugs: joinCSV(row.tag_slugs),
    sdg_codes: joinCSV(row.sdg_codes),
    organization_slugs: joinCSV(row.organization_slugs),
  }),
};

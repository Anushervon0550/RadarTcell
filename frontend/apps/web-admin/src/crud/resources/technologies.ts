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
  // Entered as percent (0..100) in the form; converted to 0..1 on submit.
  custom_metric_1: z.union([z.coerce.number().min(0).max(100), z.literal('')]).optional(),
  custom_metric_2: z.union([z.coerce.number().min(0).max(100), z.literal('')]).optional(),
  custom_metric_3: z.union([z.coerce.number().min(0).max(100), z.literal('')]).optional(),
  custom_metric_4: z.union([z.coerce.number().min(0).max(100), z.literal('')]).optional(),
});

/** Percent (0..100) shown in the form → fraction (0..1) stored in the API. */
function pctToFraction(v: unknown): number | null {
  if (v === '' || v == null) return null;
  const n = Number(v);
  if (!Number.isFinite(n)) return null;
  return Math.round((n / 100) * 10000) / 10000;
}

/** Fraction (0..1) from the API → percent (0..100) for the form. */
function fractionToPct(v: unknown): number | '' {
  if (v === '' || v == null) return '';
  const n = Number(v);
  if (!Number.isFinite(n)) return '';
  return Math.round(n * 100);
}

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
    {
      name: 'slug',
      label: 'Slug (адрес в URL)',
      type: 'text',
      requiredOnCreate: true,
      placeholder: 'edge-llm',
      hint: 'латиница, цифры и дефис. Пример: edge-llm',
    },
    {
      name: 'index',
      label: 'Порядковый номер',
      type: 'number',
      min: 0,
      required: true,
      hint: 'позиция в списке, например 31',
    },
    { name: 'name', label: 'Название', type: 'text', required: true, placeholder: 'Edge LLM' },
    {
      name: 'trend_slug',
      label: 'Тренд',
      type: 'select',
      required: true,
      optionsSource: { path: '/api/admin/trends', valueField: 'slug', labelField: 'name' },
    },
    {
      name: 'trl',
      label: 'Уровень готовности (TRL 1–9)',
      type: 'number',
      min: 1,
      max: 9,
      required: true,
      hint: '1–3 идея · 4–5 эксперимент · 6–7 внедрение · 8–9 продукт',
    },
    { name: 'source_link', label: 'Источник', type: 'url' },
    { name: 'image_url', label: 'Обложка (URL)', type: 'url' },
    { name: 'description_short', label: 'Краткое описание', type: 'textarea' },
    { name: 'description_full', label: 'Полное описание', type: 'textarea' },
    {
      name: 'tag_slugs',
      label: 'Теги',
      type: 'multiselect',
      hint: 'отметьте нужные теги',
      optionsSource: { path: '/api/admin/tags', valueField: 'slug', labelField: 'title' },
    },
    {
      name: 'sdg_codes',
      label: 'ЦУР (цели устойчивого развития)',
      type: 'multiselect',
      hint: 'отметьте подходящие цели',
      optionsSource: { path: '/api/admin/sdgs', valueField: 'code', labelField: 'title' },
    },
    {
      name: 'organization_slugs',
      label: 'Организации',
      type: 'multiselect',
      hint: 'отметьте связанные организации',
      optionsSource: { path: '/api/admin/organizations', valueField: 'slug', labelField: 'name' },
    },
    { name: 'custom_metric_1', label: 'Зрелость, %', type: 'number', min: 0, max: 100, step: 1, hint: '0–100 %' },
    { name: 'custom_metric_2', label: 'Влияние, %', type: 'number', min: 0, max: 100, step: 1, hint: '0–100 %' },
    { name: 'custom_metric_3', label: 'Покрытие, %', type: 'number', min: 0, max: 100, step: 1, hint: '0–100 %' },
    { name: 'custom_metric_4', label: 'Стоимость, %', type: 'number', min: 0, max: 100, step: 1, hint: '0–100 %' },
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
    // Form holds percent (0..100); API expects fraction (0..1).
    custom_metric_1: pctToFraction(values.custom_metric_1),
    custom_metric_2: pctToFraction(values.custom_metric_2),
    custom_metric_3: pctToFraction(values.custom_metric_3),
    custom_metric_4: pctToFraction(values.custom_metric_4),
  }),
  deserialize: (row) => ({
    ...row,
    tag_slugs: joinCSV(row.tag_slugs),
    sdg_codes: joinCSV(row.sdg_codes),
    organization_slugs: joinCSV(row.organization_slugs),
    // API stores fraction (0..1); show as percent (0..100) in the form.
    custom_metric_1: fractionToPct(row.custom_metric_1),
    custom_metric_2: fractionToPct(row.custom_metric_2),
    custom_metric_3: fractionToPct(row.custom_metric_3),
    custom_metric_4: fractionToPct(row.custom_metric_4),
  }),
};

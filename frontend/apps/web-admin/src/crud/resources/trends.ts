import { z } from 'zod';
import type { AdminTrend } from '@radartcell/api';
import type { CrudResource } from '../types';
import { slugSchema } from './shared';

const schema = z.object({
  slug: slugSchema,
  name: z.string().trim().min(1, 'обязательно').max(200),
  order_index: z.coerce.number().int().min(0).max(9999),
  description: z.string().optional().nullable(),
  image_url: z
    .string()
    .trim()
    .max(1000)
    .optional()
    .refine((v) => !v || /^https?:\/\//i.test(v), 'должен быть http(s) URL')
    .nullable()
    .or(z.literal('')),
});

export const trendsResource: CrudResource<AdminTrend & Record<string, unknown>> = {
  name: 'trends',
  title: 'Тренды',
  keyField: 'slug',
  listPath: '/api/admin/trends',
  createPath: '/api/admin/trends',
  updatePath: (k) => `/api/admin/trends/${encodeURIComponent(k)}`,
  removePath: (k) => `/api/admin/trends/${encodeURIComponent(k)}`,
  columns: [
    { key: 'slug', header: 'Slug' },
    { key: 'name', header: 'Название' },
    { key: 'order_index', header: 'Порядок' },
  ],
  fields: [
    { name: 'slug', label: 'Slug', type: 'text', requiredOnCreate: true, hint: 'a-z, 0-9, дефис' },
    { name: 'name', label: 'Название', type: 'text', required: true },
    { name: 'order_index', label: 'Порядок', type: 'number', min: 0, required: true },
    { name: 'description', label: 'Описание', type: 'textarea' },
    { name: 'image_url', label: 'URL изображения', type: 'url' },
  ],
  schema,
};

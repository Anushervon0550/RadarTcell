import { z } from 'zod';
import type { AdminTag } from '@radartcell/api';
import type { CrudResource } from '../types';
import { slugSchema } from './shared';

const schema = z.object({
  slug: slugSchema,
  title: z.string().trim().min(1).max(200),
  category: z.string().trim().min(1).max(120),
  description: z.string().optional().nullable(),
});

export const tagsResource: CrudResource<AdminTag & Record<string, unknown>> = {
  name: 'tags',
  title: 'Теги',
  keyField: 'slug',
  listPath: '/api/admin/tags',
  createPath: '/api/admin/tags',
  updatePath: (k) => `/api/admin/tags/${encodeURIComponent(k)}`,
  removePath: (k) => `/api/admin/tags/${encodeURIComponent(k)}`,
  columns: [
    { key: 'slug', header: 'Slug' },
    { key: 'title', header: 'Название' },
    { key: 'category', header: 'Категория' },
  ],
  fields: [
    { name: 'slug', label: 'Slug', type: 'text', requiredOnCreate: true },
    { name: 'title', label: 'Название', type: 'text', required: true },
    { name: 'category', label: 'Категория', type: 'text', required: true },
    { name: 'description', label: 'Описание', type: 'textarea' },
  ],
  schema,
};

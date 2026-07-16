import { z } from 'zod';
import type { AdminOrganization } from '@radartcell/api';
import type { CrudResource } from '../types';
import { slugSchema } from './shared';

const schema = z.object({
  slug: slugSchema,
  name: z.string().trim().min(1).max(200),
  logo_url: z.string().trim().max(1000).optional().nullable().or(z.literal('')),
  website: z.string().trim().max(500).optional().nullable().or(z.literal('')),
  headquarters: z.string().trim().max(200).optional().nullable().or(z.literal('')),
  description: z.string().optional().nullable(),
});

export const organizationsResource: CrudResource<AdminOrganization & Record<string, unknown>> = {
  name: 'organizations',
  title: 'Организации',
  keyField: 'slug',
  listPath: '/api/admin/organizations',
  createPath: '/api/admin/organizations',
  updatePath: (k) => `/api/admin/organizations/${encodeURIComponent(k)}`,
  removePath: (k) => `/api/admin/organizations/${encodeURIComponent(k)}`,
  columns: [
    { key: 'slug', header: 'Slug' },
    { key: 'name', header: 'Название' },
    { key: 'website', header: 'Сайт' },
    { key: 'headquarters', header: 'HQ' },
  ],
  fields: [
    { name: 'slug', label: 'Slug', type: 'text', requiredOnCreate: true },
    { name: 'name', label: 'Название', type: 'text', required: true },
    { name: 'logo_url', label: 'Логотип (URL)', type: 'url' },
    { name: 'website', label: 'Сайт', type: 'url' },
    { name: 'headquarters', label: 'Штаб-квартира', type: 'text' },
    { name: 'description', label: 'Описание', type: 'textarea' },
  ],
  schema,
};

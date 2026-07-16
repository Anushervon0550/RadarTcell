import { z } from 'zod';
import type { AdminSDG } from '@radartcell/api';
import type { CrudResource } from '../types';

const schema = z.object({
  code: z.string().trim().min(1).max(30),
  title: z.string().trim().min(1).max(300),
  icon: z.string().trim().max(1000).optional().nullable().or(z.literal('')),
  description: z.string().optional().nullable(),
});

export const sdgsResource: CrudResource<AdminSDG & Record<string, unknown>> = {
  name: 'sdgs',
  title: 'ЦУР',
  keyField: 'code',
  listPath: '/api/admin/sdgs',
  createPath: '/api/admin/sdgs',
  updatePath: (k) => `/api/admin/sdgs/${encodeURIComponent(k)}`,
  removePath: (k) => `/api/admin/sdgs/${encodeURIComponent(k)}`,
  columns: [
    { key: 'code', header: 'Код' },
    { key: 'title', header: 'Название' },
  ],
  fields: [
    { name: 'code', label: 'Код', type: 'text', requiredOnCreate: true, hint: 'например SDG 09' },
    { name: 'title', label: 'Название', type: 'text', required: true },
    { name: 'icon', label: 'Иконка (URL)', type: 'url' },
    { name: 'description', label: 'Описание', type: 'textarea' },
  ],
  schema,
};

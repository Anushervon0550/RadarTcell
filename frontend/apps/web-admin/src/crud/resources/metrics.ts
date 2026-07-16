import { z } from 'zod';
import type { AdminMetric } from '@radartcell/api';
import type { CrudResource } from '../types';

const schema = z.object({
  name: z.string().trim().min(1).max(200),
  type: z.enum(['bubble', 'bar', 'distance']),
  field_key: z.string().trim().max(60).optional().nullable().or(z.literal('')),
  orderable: z.coerce.boolean(),
  description: z.string().optional().nullable(),
});

export const metricsResource: CrudResource<AdminMetric & Record<string, unknown>> = {
  name: 'metrics',
  title: 'Метрики',
  keyField: 'id',
  listPath: '/api/admin/metrics',
  createPath: '/api/admin/metrics',
  updatePath: (k) => `/api/admin/metrics/${encodeURIComponent(k)}`,
  removePath: (k) => `/api/admin/metrics/${encodeURIComponent(k)}`,
  columns: [
    { key: 'id', header: 'ID', width: '260px' },
    { key: 'name', header: 'Название' },
    { key: 'type', header: 'Тип' },
    { key: 'field_key', header: 'field_key' },
    { key: 'orderable', header: 'Orderable' },
  ],
  fields: [
    { name: 'name', label: 'Название', type: 'text', required: true },
    {
      name: 'type',
      label: 'Тип',
      type: 'select',
      options: [
        { value: 'bubble', label: 'bubble' },
        { value: 'bar', label: 'bar' },
        { value: 'distance', label: 'distance' },
      ],
      required: true,
    },
    { name: 'field_key', label: 'Field key (для legacy)', type: 'text' },
    { name: 'orderable', label: 'Использовать в сортировке', type: 'checkbox' },
    { name: 'description', label: 'Описание', type: 'textarea' },
  ],
  schema,
};

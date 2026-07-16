import { z } from 'zod';
import type { AdminSDG } from '@radartcell/api';
import type { CrudResource } from '../types';

/**
 * Normalize an SDG code: accepts "3", "03", "sdg 3", "SDG 03" → "SDG 03".
 * Falls back to the trimmed input if no number is found.
 */
function normalizeSdgCode(input: unknown): string {
  const s = String(input ?? '').trim();
  const m = s.match(/\d+/);
  if (!m) return s;
  const n = parseInt(m[0], 10);
  if (Number.isNaN(n)) return s;
  return `SDG ${String(n).padStart(2, '0')}`;
}

const schema = z.object({
  code: z
    .string()
    .trim()
    .min(1, 'обязательно')
    .max(30)
    .transform(normalizeSdgCode)
    .refine((v) => /^SDG \d{2}$/.test(v), 'введите номер цели, например 9'),
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
    {
      name: 'code',
      label: 'Номер цели',
      type: 'text',
      requiredOnCreate: true,
      placeholder: '9',
      hint: 'просто число: 9 → SDG 09',
    },
    { name: 'title', label: 'Название', type: 'text', required: true },
    { name: 'icon', label: 'Иконка (URL)', type: 'url' },
    { name: 'description', label: 'Описание', type: 'textarea' },
  ],
  schema,
};

import { z } from 'zod';

/**
 * Split a comma-separated string into a clean list of slugs/codes.
 * Empty items are removed.
 */
export function splitCSV(input: unknown): string[] {
  if (Array.isArray(input)) return input.map(String).map((s) => s.trim()).filter(Boolean);
  const s = String(input ?? '').trim();
  if (!s) return [];
  return s.split(',').map((p) => p.trim()).filter(Boolean);
}

export function joinCSV(input: unknown): string {
  if (!Array.isArray(input)) return '';
  return input.map(String).join(', ');
}

/** Slug validation: kebab-case, digits allowed, no leading/trailing dash. */
export const slugSchema = z
  .string()
  .trim()
  .min(1, 'обязательно')
  .max(120)
  .regex(/^[a-z0-9]+(?:-[a-z0-9]+)*$/, 'kebab-case: a-z, 0-9, дефис');

export type Stage = 'idea' | 'prototype' | 'product';

export function stageFromTrl(trl: number | string | null | undefined): Stage {
  const n = Number(trl) || 0;
  if (n <= 3) return 'idea';
  if (n <= 6) return 'prototype';
  return 'product';
}

const STAGE_LABEL_RU: Record<Stage, string> = {
  idea: 'Идея',
  prototype: 'Прототип',
  product: 'Продукт',
};

export function stageLabel(stage: Stage | string | null | undefined): string {
  if (stage && stage in STAGE_LABEL_RU) return STAGE_LABEL_RU[stage as Stage];
  return String(stage ?? '—');
}

/**
 * Whitelist-based URL sanitizer. Returns '' for anything that isn't
 * http(s), a data:image, blob:, or a relative path. Blocks javascript: etc.
 */
export function safeUrl(value: unknown): string {
  const s = String(value ?? '').trim();
  if (!s) return '';
  if (/^(https?:|\/|#|mailto:|tel:)/i.test(s)) return s;
  if (/^data:image\//i.test(s)) return s;
  if (/^blob:/i.test(s)) return s;
  return '';
}

/**
 * Safe CSS `url()` fragment for inline background-image styles.
 * Encodes characters that could break out of the style attribute.
 */
export function cssUrl(value: unknown): string {
  const u = safeUrl(value);
  if (!u) return 'none';
  const safe = u.replace(/[)'\s\\<>]/g, encodeURIComponent);
  return `url('${safe}')`;
}

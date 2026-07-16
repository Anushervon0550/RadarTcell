import type { ReactNode } from 'react';
import type { ZodTypeAny } from 'zod';

export type FieldType =
  | 'text'
  | 'textarea'
  | 'number'
  | 'checkbox'
  | 'select'
  | 'multi'
  | 'url';

export interface FieldOption {
  value: string;
  label: string;
}

export interface FieldDef {
  name: string;
  label: string;
  type: FieldType;
  placeholder?: string;
  required?: boolean;
  requiredOnCreate?: boolean;
  options?: FieldOption[];
  hint?: string;
  min?: number;
  max?: number;
  step?: number;
  disabled?: boolean;
}

export interface ColumnDef<T> {
  key: keyof T | string;
  header: string;
  render?: (row: T) => ReactNode;
  width?: string;
}

export interface CrudResource<T extends Record<string, unknown>> {
  name: string;
  title: string;
  keyField: keyof T & string;
  listPath: string;
  createPath: string;
  updatePath: (key: string) => string;
  removePath: (key: string) => string;
  restorePath?: (key: string) => string;
  columns: ColumnDef<T>[];
  fields: FieldDef[];
  schema: ZodTypeAny;
  /** Optional response unwrapper. Default: array or `{items:[]}`. */
  parseList?: (data: unknown) => T[];
  /** Convert form values before POST/PUT (e.g. split "a,b" → ["a","b"]). */
  serialize?: (values: Record<string, unknown>) => Record<string, unknown>;
  /** Convert row into form defaults (e.g. join array → "a,b"). */
  deserialize?: (row: T) => Record<string, unknown>;
  /** Optional: show restore button when row has `deleted_at`. */
  softDelete?: boolean;
}

export function unwrapList<T>(data: unknown): T[] {
  if (Array.isArray(data)) return data as T[];
  if (data && typeof data === 'object' && Array.isArray((data as { items?: unknown }).items)) {
    return (data as { items: T[] }).items;
  }
  return [];
}

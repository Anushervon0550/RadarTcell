import { useEffect } from 'react';
import { useForm, type SubmitHandler } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useQuery } from '@tanstack/react-query';
import { Button, Input, Select, Textarea } from '@radartcell/ui';
import { api } from '@/api/client';
import { unwrapList, type CrudResource, type FieldDef, type FieldOption } from './types';

interface Props<T extends Record<string, unknown>> {
  resource: CrudResource<T>;
  initial?: T | null;
  onSubmit: (values: Record<string, unknown>) => Promise<void> | void;
  onCancel?: () => void;
  submitting?: boolean;
}

export function ResourceForm<T extends Record<string, unknown>>({
  resource,
  initial,
  onSubmit,
  onCancel,
  submitting,
}: Props<T>) {
  const isCreate = !initial;
  const defaultValues = isCreate
    ? blankDefaults(resource.fields)
    : (resource.deserialize?.(initial as T) ?? (initial as unknown as Record<string, unknown>));

  const form = useForm<Record<string, unknown>>({
    defaultValues,
    resolver: zodResolver(resource.schema),
  });

  useEffect(() => {
    form.reset(defaultValues);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [initial]);

  const handleSubmit: SubmitHandler<Record<string, unknown>> = async (values) => {
    const payload = resource.serialize ? resource.serialize(values) : values;
    await onSubmit(payload);
  };

  return (
    <form onSubmit={form.handleSubmit(handleSubmit)} className="grid grid-cols-1 gap-3 sm:grid-cols-2">
      {resource.fields.map((f) => {
        const isWide = f.type === 'textarea' || f.type === 'multi' || f.type === 'multiselect';
        return (
          <div key={f.name} className={isWide ? 'sm:col-span-2' : ''}>
            <FieldRenderer field={f} form={form} isCreate={isCreate} />
          </div>
        );
      })}
      <div className="sticky bottom-0 -mx-6 mt-2 flex flex-wrap gap-2 border-t border-line bg-[#0c1326]/95 px-6 py-3 backdrop-blur sm:col-span-2">
        <Button variant="primary" type="submit" loading={submitting}>
          {isCreate ? 'Создать' : 'Сохранить'}
        </Button>
        {onCancel && (
          <Button variant="ghost" type="button" onClick={onCancel}>
            Отмена
          </Button>
        )}
        <span className="ml-auto self-center text-xs text-ink-muted-2">
          Esc — закрыть
        </span>
      </div>
    </form>
  );
}

function blankDefaults(fields: FieldDef[]): Record<string, unknown> {
  const out: Record<string, unknown> = {};
  for (const f of fields) {
    if (f.type === 'checkbox') out[f.name] = false;
    else if (f.type === 'number') out[f.name] = undefined;
    else out[f.name] = '';
  }
  return out;
}

function FieldRenderer({
  field,
  form,
  isCreate,
}: {
  field: FieldDef;
  form: ReturnType<typeof useForm<Record<string, unknown>>>;
  isCreate: boolean;
}) {
  const err = form.formState.errors[field.name]?.message as string | undefined;
  const disabled = field.disabled || (!isCreate && field.name === 'slug' && field.requiredOnCreate);

  switch (field.type) {
    case 'textarea':
      return (
        <Textarea
          label={field.label}
          error={err}
          hint={field.hint}
          placeholder={field.placeholder}
          disabled={disabled}
          {...form.register(field.name)}
        />
      );
    case 'number':
      return (
        <Input
          type="number"
          label={field.label}
          error={err}
          hint={field.hint}
          placeholder={field.placeholder}
          min={field.min}
          max={field.max}
          step={field.step ?? 'any'}
          disabled={disabled}
          {...form.register(field.name, { valueAsNumber: true })}
        />
      );
    case 'checkbox':
      return (
        <label className="flex items-center gap-2 text-sm text-ink">
          <input
            type="checkbox"
            className="h-4 w-4 rounded border-line bg-bg-soft"
            disabled={disabled}
            {...form.register(field.name)}
          />
          {field.label}
          {err && <span className="text-xs text-red-400">{err}</span>}
        </label>
      );
    case 'select':
      return (
        <SelectField field={field} form={form} err={err} disabled={disabled} />
      );
    case 'multiselect':
      return (
        <MultiSelectField field={field} form={form} err={err} disabled={disabled} />
      );
    case 'multi':
      return (
        <Input
          label={field.label}
          error={err}
          hint={field.hint ?? 'Значения через запятую'}
          placeholder={field.placeholder ?? 'a, b, c'}
          disabled={disabled}
          {...form.register(field.name)}
        />
      );
    default:
      return (
        <Input
          label={field.label}
          error={err}
          hint={field.hint}
          placeholder={field.placeholder}
          disabled={disabled}
          {...form.register(field.name)}
        />
      );
  }
}

function SelectField({
  field,
  form,
  err,
  disabled,
}: {
  field: FieldDef;
  form: ReturnType<typeof useForm<Record<string, unknown>>>;
  err?: string;
  disabled?: boolean;
}) {
  const src = field.optionsSource;
  const { data: dynamicOptions, isLoading } = useQuery({
    queryKey: ['field-options', src?.path],
    queryFn: async (): Promise<FieldOption[]> => {
      const raw = await api.get<unknown>(src!.path);
      const rows = unwrapList<Record<string, unknown>>(raw);
      const valueKey = src!.valueField ?? 'slug';
      const labelKey = src!.labelField ?? 'name';
      return rows.map((r) => ({
        value: String(r[valueKey] ?? ''),
        label: String(r[labelKey] ?? r[valueKey] ?? ''),
      }));
    },
    enabled: !!src,
    staleTime: 60_000,
  });

  const options = src ? dynamicOptions ?? [] : field.options ?? [];

  return (
    <Select
      label={field.label}
      error={err}
      hint={field.hint}
      placeholder={isLoading ? 'Загрузка…' : '— выберите —'}
      options={options}
      disabled={disabled}
      {...form.register(field.name)}
    />
  );
}

function MultiSelectField({
  field,
  form,
  err,
  disabled,
}: {
  field: FieldDef;
  form: ReturnType<typeof useForm<Record<string, unknown>>>;
  err?: string;
  disabled?: boolean;
}) {
  const src = field.optionsSource;
  const { data: options, isLoading } = useQuery({
    queryKey: ['field-options', src?.path],
    queryFn: async (): Promise<FieldOption[]> => {
      const raw = await api.get<unknown>(src!.path);
      const rows = unwrapList<Record<string, unknown>>(raw);
      const valueKey = src!.valueField ?? 'slug';
      const labelKey = src!.labelField ?? 'name';
      return rows.map((r) => ({
        value: String(r[valueKey] ?? ''),
        label: String(r[labelKey] ?? r[valueKey] ?? ''),
      }));
    },
    enabled: !!src,
    staleTime: 60_000,
  });

  // The stored value is a comma-separated string (kept for API compatibility).
  const raw = form.watch(field.name);
  const selected = String(raw ?? '')
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean);

  const toggle = (value: string) => {
    const set = new Set(selected);
    if (set.has(value)) set.delete(value);
    else set.add(value);
    form.setValue(field.name, Array.from(set).join(', '), { shouldValidate: true });
  };

  return (
    <div className="flex flex-col gap-1.5">
      <div className="flex items-center justify-between">
        <label className="text-xs text-ink-muted">{field.label}</label>
        {selected.length > 0 && (
          <span className="text-[11px] text-brand-400">выбрано: {selected.length}</span>
        )}
      </div>
      <div className="max-h-44 overflow-y-auto rounded-lg border border-line bg-bg-soft p-2">
        {isLoading ? (
          <div className="px-1 py-2 text-xs text-ink-muted">Загрузка…</div>
        ) : !options || options.length === 0 ? (
          <div className="px-1 py-2 text-xs text-ink-muted">Нет доступных вариантов</div>
        ) : (
          <div className="grid grid-cols-1 gap-1 sm:grid-cols-2">
            {options.map((o) => {
              const checked = selected.includes(o.value);
              return (
                <label
                  key={o.value}
                  className={
                    'flex cursor-pointer items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors ' +
                    (checked ? 'bg-brand-600/15 text-white' : 'hover:bg-white/5 text-ink')
                  }
                >
                  <input
                    type="checkbox"
                    className="h-4 w-4 rounded border-line bg-bg-soft"
                    checked={checked}
                    disabled={disabled}
                    onChange={() => toggle(o.value)}
                  />
                  <span className="truncate">{o.label}</span>
                </label>
              );
            })}
          </div>
        )}
      </div>
      {(err || field.hint) && (
        <span className={err ? 'text-xs text-red-400' : 'text-xs text-ink-muted-2'}>
          {err ?? field.hint}
        </span>
      )}
    </div>
  );
}

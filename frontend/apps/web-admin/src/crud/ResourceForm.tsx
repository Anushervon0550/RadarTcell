import { useEffect } from 'react';
import { useForm, type SubmitHandler } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Button, Input, Select, Textarea } from '@radartcell/ui';
import type { CrudResource, FieldDef } from './types';

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
    <form onSubmit={form.handleSubmit(handleSubmit)} className="grid grid-cols-1 gap-3">
      {resource.fields.map((f) => (
        <FieldRenderer key={f.name} field={f} form={form} isCreate={isCreate} />
      ))}
      <div className="mt-2 flex flex-wrap gap-2">
        <Button variant="primary" type="submit" loading={submitting}>
          {isCreate ? 'Создать' : 'Сохранить'}
        </Button>
        {onCancel && (
          <Button variant="ghost" type="button" onClick={onCancel}>
            Отмена
          </Button>
        )}
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
        <Select
          label={field.label}
          error={err}
          hint={field.hint}
          options={field.options ?? []}
          disabled={disabled}
          {...form.register(field.name)}
        />
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

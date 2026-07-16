import { forwardRef, type SelectHTMLAttributes, useId } from 'react';
import { cn } from '../lib/cn';

export interface SelectOption {
  value: string;
  label: string;
}

export interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
  label?: string;
  error?: string;
  hint?: string;
  options: SelectOption[];
  placeholder?: string;
}

export const Select = forwardRef<HTMLSelectElement, SelectProps>(function Select(
  { className, label, error, hint, options, placeholder, id, ...rest },
  ref,
) {
  const autoId = useId();
  const inputId = id ?? autoId;
  return (
    <div className="flex flex-col gap-1.5">
      {label && (
        <label htmlFor={inputId} className="text-xs text-ink-muted">
          {label}
        </label>
      )}
      <select
        id={inputId}
        ref={ref}
        className={cn(
          'w-full rounded-lg border border-line bg-bg-soft px-3 py-2 text-sm text-ink',
          'rt-focus-ring appearance-none',
          error && 'border-red-500 focus-visible:ring-red-500',
          className,
        )}
        {...rest}
      >
        {placeholder && <option value="">{placeholder}</option>}
        {options.map((o) => (
          <option key={o.value} value={o.value}>
            {o.label}
          </option>
        ))}
      </select>
      {(error || hint) && (
        <span className={cn('text-xs', error ? 'text-red-400' : 'text-ink-muted-2')}>
          {error ?? hint}
        </span>
      )}
    </div>
  );
});

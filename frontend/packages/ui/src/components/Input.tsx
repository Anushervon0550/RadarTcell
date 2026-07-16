import { forwardRef, type InputHTMLAttributes, useId } from 'react';
import { cn } from '../lib/cn';

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  hint?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(function Input(
  { className, label, error, hint, id, ...rest },
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
      <input
        id={inputId}
        ref={ref}
        className={cn(
          'w-full rounded-lg border border-line bg-bg-soft px-3 py-2 text-sm text-ink placeholder:text-ink-muted-2',
          'rt-focus-ring',
          error && 'border-red-500 focus-visible:ring-red-500',
          className,
        )}
        aria-invalid={error ? true : undefined}
        aria-describedby={error || hint ? `${inputId}-help` : undefined}
        {...rest}
      />
      {(error || hint) && (
        <span
          id={`${inputId}-help`}
          className={cn('text-xs', error ? 'text-red-400' : 'text-ink-muted-2')}
        >
          {error ?? hint}
        </span>
      )}
    </div>
  );
});

import { forwardRef, type TextareaHTMLAttributes, useId } from 'react';
import { cn } from '../lib/cn';

export interface TextareaProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string;
  error?: string;
  hint?: string;
}

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(function Textarea(
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
      <textarea
        id={inputId}
        ref={ref}
        className={cn(
          'min-h-[90px] w-full resize-y rounded-lg border border-line bg-bg-soft px-3 py-2 text-sm text-ink placeholder:text-ink-muted-2',
          'rt-focus-ring',
          error && 'border-red-500 focus-visible:ring-red-500',
          className,
        )}
        aria-invalid={error ? true : undefined}
        {...rest}
      />
      {(error || hint) && (
        <span className={cn('text-xs', error ? 'text-red-400' : 'text-ink-muted-2')}>
          {error ?? hint}
        </span>
      )}
    </div>
  );
});

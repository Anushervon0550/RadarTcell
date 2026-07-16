import { forwardRef, type ButtonHTMLAttributes, type ReactNode } from 'react';
import { cn } from '../lib/cn';

export type ButtonVariant = 'primary' | 'secondary' | 'ghost' | 'danger' | 'success' | 'warn';
export type ButtonSize = 'sm' | 'md' | 'lg';

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  size?: ButtonSize;
  loading?: boolean;
  icon?: ReactNode;
  iconRight?: ReactNode;
  fullWidth?: boolean;
}

const VARIANTS: Record<ButtonVariant, string> = {
  primary:
    'bg-gradient-to-br from-brand-600 to-[#4f46e5] text-white border-brand-700 hover:brightness-110 shadow-glow',
  secondary:
    'bg-bg-panel border-line text-ink hover:border-brand-600/60 hover:bg-bg-panel-2',
  ghost:
    'bg-transparent border-transparent text-ink hover:bg-white/5',
  danger:
    'bg-gradient-to-br from-red-800 to-red-500 text-white border-red-900 hover:brightness-110',
  success:
    'bg-gradient-to-br from-emerald-700 to-emerald-500 text-white border-emerald-800 hover:brightness-110',
  warn:
    'bg-gradient-to-br from-amber-700 to-amber-500 text-white border-amber-800 hover:brightness-110',
};

const SIZES: Record<ButtonSize, string> = {
  sm: 'px-2.5 py-1.5 text-xs rounded-lg',
  md: 'px-3.5 py-2 text-sm rounded-lg',
  lg: 'px-5 py-2.5 text-base rounded-xl',
};

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(function Button(
  {
    className,
    variant = 'secondary',
    size = 'md',
    loading,
    icon,
    iconRight,
    children,
    disabled,
    type = 'button',
    fullWidth,
    ...rest
  },
  ref,
) {
  return (
    <button
      ref={ref}
      type={type}
      disabled={disabled || loading}
      className={cn(
        'inline-flex items-center justify-center gap-2 border font-medium transition-all',
        'rt-focus-ring active:translate-y-px disabled:opacity-60 disabled:cursor-not-allowed',
        VARIANTS[variant],
        SIZES[size],
        fullWidth && 'w-full',
        className,
      )}
      {...rest}
    >
      {loading ? (
        <span
          aria-hidden
          className="inline-block h-4 w-4 animate-spin rounded-full border-2 border-white/40 border-t-white"
        />
      ) : (
        icon
      )}
      {children}
      {iconRight}
    </button>
  );
});

import { useState } from 'react';
import { cn, safeUrl } from '@radartcell/ui';

interface Props {
  name: string;
  logoUrl?: string | null;
  /** Extra classes for sizing / rounding (e.g. "h-14 w-14 rounded-xl"). */
  className?: string;
}

/** Derive up to two initials from an organization name. */
function initials(name: string): string {
  const parts = name.trim().split(/\s+/).filter(Boolean);
  if (parts.length === 0) return '?';
  if (parts.length === 1) return parts[0]!.slice(0, 2).toUpperCase();
  return (parts[0]![0]! + parts[parts.length - 1]![0]!).toUpperCase();
}

/**
 * Organization logo with a graceful fallback: shows the image when it loads,
 * otherwise renders a coloured initials avatar. Avoids the "white square"
 * placeholder that appears when logo_url is empty or points to a dead URL.
 */
export function OrgLogo({ name, logoUrl, className }: Props) {
  const url = safeUrl(logoUrl);
  const [failed, setFailed] = useState(false);
  const showImage = url !== '' && !failed;

  if (showImage) {
    return (
      <img
        src={url}
        alt={`Логотип ${name}`}
        loading="lazy"
        onError={() => setFailed(true)}
        className={cn(
          'shrink-0 border border-line bg-white object-cover',
          className,
        )}
      />
    );
  }

  return (
    <span
      role="img"
      aria-label={`Логотип ${name}`}
      className={cn(
        'flex shrink-0 items-center justify-center border border-line bg-brand-600/20 font-semibold uppercase text-brand-200',
        className,
      )}
    >
      {initials(name)}
    </span>
  );
}

import { Link } from 'react-router-dom';
import type { TechnologyListItem } from '@radartcell/api';
import { cssUrl, Pill, StageChip, TrlChip } from '@radartcell/ui';
import { FALLBACK_COVER } from '@/pages/radar/constants';

export function TechCard({ item }: { item: TechnologyListItem }) {
  const cover = item.image_url ?? FALLBACK_COVER;
  return (
    <Link
      to={`/technology/${encodeURIComponent(item.slug)}`}
      className="group flex flex-col overflow-hidden rounded-card border border-line bg-bg-panel/95 shadow-elevated transition-all hover:-translate-y-1 hover:border-brand-600/50 hover:shadow-[0_22px_42px_rgba(124,58,237,0.18)]"
    >
      <div
        className="aspect-[16/9] w-full bg-cover bg-center"
        style={{ backgroundImage: cssUrl(cover) }}
      />
      <div className="flex flex-1 flex-col gap-2 p-4">
        <h3 className="text-base font-semibold leading-snug">{item.name}</h3>
        <div className="flex flex-wrap items-center gap-1.5">
          <TrlChip trl={item.trl} />
          <StageChip trl={item.trl} />
          <Pill>{item.trend_name || item.trend_slug}</Pill>
        </div>
        {item.description_short && (
          <p className="line-clamp-3 text-sm leading-relaxed text-ink-muted">
            {item.description_short}
          </p>
        )}
      </div>
    </Link>
  );
}

import { cn } from '../lib/cn';
import { stageFromTrl, stageLabel, type Stage } from '../lib/utils';

const STAGE_COLORS: Record<Stage, string> = {
  idea: 'border-amber-500/45 bg-amber-500/15 text-amber-300',
  prototype: 'border-sky-500/45 bg-sky-500/15 text-sky-300',
  product: 'border-emerald-500/45 bg-emerald-500/15 text-emerald-300',
};

export function StageChip({
  trl,
  stage,
  className,
}: {
  trl?: number;
  stage?: Stage;
  className?: string;
}) {
  const resolved: Stage = stage ?? stageFromTrl(trl);
  return (
    <span
      className={cn(
        'inline-flex items-center gap-1 rounded-full border px-2.5 py-0.5 text-xs font-semibold',
        STAGE_COLORS[resolved],
        className,
      )}
    >
      {stageLabel(resolved)}
    </span>
  );
}

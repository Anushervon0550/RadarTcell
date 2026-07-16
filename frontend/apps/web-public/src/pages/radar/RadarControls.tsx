export function RadarControls({
  onZoomIn,
  onZoomOut,
  onReset,
}: {
  onZoomIn: () => void;
  onZoomOut: () => void;
  onReset: () => void;
}) {
  return (
    <div className="absolute bottom-3 right-3 z-[4] flex flex-col gap-1 rounded-xl border border-line bg-black/70 p-1 backdrop-blur">
      <button
        type="button"
        onClick={onZoomIn}
        title="Приблизить"
        aria-label="Приблизить"
        className="h-9 w-9 rounded-lg text-lg font-bold hover:bg-brand-600/25"
      >
        +
      </button>
      <button
        type="button"
        onClick={onZoomOut}
        title="Отдалить"
        aria-label="Отдалить"
        className="h-9 w-9 rounded-lg text-lg font-bold hover:bg-brand-600/25"
      >
        −
      </button>
      <button
        type="button"
        onClick={onReset}
        title="Сбросить"
        aria-label="Сбросить масштаб"
        className="h-9 w-9 rounded-lg text-sm hover:bg-brand-600/25"
      >
        ⟲
      </button>
    </div>
  );
}

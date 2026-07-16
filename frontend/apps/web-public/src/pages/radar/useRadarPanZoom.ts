import { useCallback, useEffect, useRef, useState, type RefObject } from 'react';

interface Transform {
  scale: number;
  tx: number;
  ty: number;
}

interface Options {
  min?: number;
  max?: number;
  viewSize: number;
}

/**
 * Attach pan (drag) + zoom (wheel) interactions to an inner <g> group inside <svg>.
 * All transforms are applied via the SVG `transform` attribute on the inner group.
 *
 * Translation is clamped so the centre of the radar can never be dragged/zoomed
 * completely outside the viewBox (otherwise the whole chart "disappears").
 */
export function useRadarPanZoom(
  svgRef: RefObject<SVGSVGElement | null>,
  rootRef: RefObject<SVGGElement | null>,
  { min = 0.5, max = 5, viewSize }: Options,
) {
  const [zoom, setZoom] = useState(1);
  const transformRef = useRef<Transform>({ scale: 1, tx: 0, ty: 0 });

  const clamp = (v: number, lo: number, hi: number) => Math.min(hi, Math.max(lo, v));

  // Keep the geometric centre of the (scaled) radar inside the viewBox.
  const clampTranslate = useCallback(
    (t: Transform) => {
      const c = viewSize / 2;
      const centerX = t.tx + t.scale * c;
      const centerY = t.ty + t.scale * c;
      t.tx = clamp(centerX, 0, viewSize) - t.scale * c;
      t.ty = clamp(centerY, 0, viewSize) - t.scale * c;
    },
    [viewSize],
  );

  const applyTransform = useCallback(() => {
    const t = transformRef.current;
    const root = rootRef.current;
    if (root) {
      root.setAttribute('transform', `translate(${t.tx},${t.ty}) scale(${t.scale})`);
    }
    setZoom(t.scale);
  }, [rootRef]);

  useEffect(() => {
    const svg = svgRef.current;
    const root = rootRef.current;
    if (!svg || !root) return;

    let pointerDown = false;
    let dragging = false;
    let pointerId = -1;
    let startX = 0;
    let startY = 0;
    let startTx = 0;
    let startTy = 0;
    const DRAG_THRESHOLD = 4; // px before a press turns into a drag

    const toSvgCoords = (evtX: number, evtY: number) => {
      const rect = svg.getBoundingClientRect();
      const nx = (evtX - rect.left) / rect.width;
      const ny = (evtY - rect.top) / rect.height;
      return { x: nx * viewSize, y: ny * viewSize };
    };

    const onWheel = (e: WheelEvent) => {
      e.preventDefault();
      const t = transformRef.current;
      const delta = -e.deltaY / 500;
      const nextScale = clamp(t.scale * (1 + delta), min, max);
      if (nextScale === t.scale) return;
      const { x, y } = toSvgCoords(e.clientX, e.clientY);
      // Keep the cursor point stable in SVG space during zoom.
      const factor = nextScale / t.scale;
      t.tx = x - factor * (x - t.tx);
      t.ty = y - factor * (y - t.ty);
      t.scale = nextScale;
      clampTranslate(t);
      applyTransform();
    };

    const onPointerDown = (e: PointerEvent) => {
      if (e.button !== 0) return;
      // Don't capture the pointer yet: capturing here would steal the `click`
      // event from radar dots. We only start a real drag once the pointer
      // moves past a small threshold.
      pointerDown = true;
      dragging = false;
      pointerId = e.pointerId;
      startX = e.clientX;
      startY = e.clientY;
      startTx = transformRef.current.tx;
      startTy = transformRef.current.ty;
    };

    const onPointerMove = (e: PointerEvent) => {
      if (!pointerDown) return;

      if (!dragging) {
        const dx = e.clientX - startX;
        const dy = e.clientY - startY;
        if (Math.hypot(dx, dy) < DRAG_THRESHOLD) return;
        // Promote to a drag: now capture the pointer so panning is smooth.
        dragging = true;
        svg.setPointerCapture(pointerId);
        svg.classList.add('cursor-grabbing');
      }

      const rect = svg.getBoundingClientRect();
      const px = ((e.clientX - startX) / rect.width) * viewSize;
      const py = ((e.clientY - startY) / rect.height) * viewSize;
      const t = transformRef.current;
      t.tx = startTx + px;
      t.ty = startTy + py;
      clampTranslate(t);
      applyTransform();
    };

    const onPointerUp = (e: PointerEvent) => {
      if (!pointerDown) return;
      pointerDown = false;
      if (dragging) {
        dragging = false;
        svg.releasePointerCapture(e.pointerId);
        svg.classList.remove('cursor-grabbing');
      }
    };

    svg.addEventListener('wheel', onWheel, { passive: false });
    svg.addEventListener('pointerdown', onPointerDown);
    svg.addEventListener('pointermove', onPointerMove);
    svg.addEventListener('pointerup', onPointerUp);
    svg.addEventListener('pointercancel', onPointerUp);

    return () => {
      svg.removeEventListener('wheel', onWheel);
      svg.removeEventListener('pointerdown', onPointerDown);
      svg.removeEventListener('pointermove', onPointerMove);
      svg.removeEventListener('pointerup', onPointerUp);
      svg.removeEventListener('pointercancel', onPointerUp);
    };
  }, [svgRef, rootRef, min, max, viewSize, applyTransform, clampTranslate]);

  // Zoom by a multiplicative factor around the viewport centre (used by +/- buttons).
  const zoomBy = useCallback(
    (factor: number) => {
      const t = transformRef.current;
      const nextScale = clamp(t.scale * factor, min, max);
      if (nextScale === t.scale) return;
      const pivot = viewSize / 2;
      const f = nextScale / t.scale;
      t.tx = pivot - f * (pivot - t.tx);
      t.ty = pivot - f * (pivot - t.ty);
      t.scale = nextScale;
      clampTranslate(t);
      applyTransform();
    },
    [min, max, viewSize, clampTranslate, applyTransform],
  );

  const reset = useCallback(() => {
    transformRef.current = { scale: 1, tx: 0, ty: 0 };
    applyTransform();
  }, [applyTransform]);

  return { zoom, zoomBy, reset };
}

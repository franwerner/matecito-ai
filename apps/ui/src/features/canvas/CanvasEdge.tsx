// Fixed geometry matching CanvasEdge.dc.html — the reference exposes only
// `label` as a prop, the curve itself is illustrative, not derived from points.
const CANVAS_EDGE_VIEWBOX = '0 0 240 130'
const CANVAS_EDGE_PATH = 'M10,104 C96,104 96,26 230,26'
const CANVAS_EDGE_ARROW_PATH = 'M226,22 L234,26 L226,30 Z'
const CANVAS_EDGE_ORIGIN_X = 10
const CANVAS_EDGE_ORIGIN_Y = 104
const CANVAS_EDGE_ORIGIN_RADIUS = 4
const CANVAS_EDGE_STROKE_WIDTH = 2
const CANVAS_EDGE_STROKE_OPACITY = 0.5

export interface CanvasEdgeProps {
  label: string
}

export function CanvasEdge({ label }: CanvasEdgeProps) {
  return (
    <div className="relative h-full min-h-[130px] w-full min-w-[220px] font-sans">
      <svg
        width="100%"
        height="100%"
        viewBox={CANVAS_EDGE_VIEWBOX}
        preserveAspectRatio="none"
        className="absolute inset-0 overflow-visible"
        aria-hidden="true"
      >
        <path
          d={CANVAS_EDGE_PATH}
          fill="none"
          className="stroke-primary"
          strokeWidth={CANVAS_EDGE_STROKE_WIDTH}
          strokeOpacity={CANVAS_EDGE_STROKE_OPACITY}
          strokeLinecap="round"
        />
        <circle
          cx={CANVAS_EDGE_ORIGIN_X}
          cy={CANVAS_EDGE_ORIGIN_Y}
          r={CANVAS_EDGE_ORIGIN_RADIUS}
          className="fill-primary"
        />
        <path d={CANVAS_EDGE_ARROW_PATH} className="fill-primary" />
      </svg>
      <span className="border-border bg-card absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 rounded-md border px-2.5 py-[3px] font-mono text-[11px] whitespace-nowrap text-[#5A6675]">
        {label}
      </span>
    </div>
  )
}

import { useId } from 'react'

export type LabeledEdgeState = 'inactive' | 'active'

// LabeledEdge's dark palette is invariant across themes (rail-of-the-cockpit
// criterion, see globals.css) — it stays local to this module as named
// constants, never as inline literals, per structure/code-conventions.md.
// Closed set: read by LABELED_EDGE_STROKE[state]/LABELED_EDGE_DASH[state],
// never compared against the 'active'/'inactive' string literal.
const LABELED_EDGE_STROKE: Record<LabeledEdgeState, string> = {
  inactive: '#4A3D31',
  active: '#8B7BFF',
}

const LABELED_EDGE_DASH: Record<LabeledEdgeState, string> = {
  inactive: '0',
  active: '4 4',
}

const LABELED_EDGE_VIEWBOX = '0 0 220 80'
const LABELED_EDGE_PATH = 'M6,40 C 78,40 142,40 210,40'
const LABELED_EDGE_STROKE_WIDTH = 1.5
const LABELED_EDGE_STROKE_OPACITY = 0.85
const LABELED_EDGE_MARKER_SIZE = 7
const LABELED_EDGE_MARKER_REF_X = 5.5
const LABELED_EDGE_MARKER_REF_Y = 3
const LABELED_EDGE_ARROWHEAD_PATH = 'M0,0 L6,3 L0,6 Z'

const LABELED_EDGE_CHIP_BACKGROUND = '#1A1410'
const LABELED_EDGE_CHIP_BORDER = '#3A2F26'
const LABELED_EDGE_CHIP_FOREGROUND = '#C7CAD0'

export interface LabeledEdgeProps {
  label: string
  active?: boolean
}

export function LabeledEdge({ label, active = false }: LabeledEdgeProps) {
  // useId keeps the <marker> reference unique per instance, so several
  // LabeledEdge renders (mixed active/inactive) never fight over one id.
  const markerId = `labeled-edge-arrow-${useId()}`
  const state: LabeledEdgeState = active ? 'active' : 'inactive'
  const strokeColor = LABELED_EDGE_STROKE[state]
  const dashArray = LABELED_EDGE_DASH[state]

  return (
    <div className="relative flex h-full min-h-[80px] w-full items-center justify-center">
      <svg
        viewBox={LABELED_EDGE_VIEWBOX}
        width="100%"
        height="100%"
        preserveAspectRatio="none"
        className="absolute inset-0 overflow-visible"
        aria-hidden="true"
      >
        <defs>
          <marker
            id={markerId}
            markerWidth={LABELED_EDGE_MARKER_SIZE}
            markerHeight={LABELED_EDGE_MARKER_SIZE}
            refX={LABELED_EDGE_MARKER_REF_X}
            refY={LABELED_EDGE_MARKER_REF_Y}
            orient="auto"
          >
            <path d={LABELED_EDGE_ARROWHEAD_PATH} fill={strokeColor} />
          </marker>
        </defs>
        <path
          d={LABELED_EDGE_PATH}
          fill="none"
          stroke={strokeColor}
          strokeWidth={LABELED_EDGE_STROKE_WIDTH}
          strokeDasharray={dashArray}
          strokeOpacity={LABELED_EDGE_STROKE_OPACITY}
          markerEnd={`url(#${markerId})`}
        />
      </svg>
      <span
        className="relative z-1 rounded-[5px] border px-2.5 py-0.5 font-mono text-[10.5px] font-medium tracking-[.03em] whitespace-nowrap"
        style={{
          backgroundColor: LABELED_EDGE_CHIP_BACKGROUND,
          borderColor: LABELED_EDGE_CHIP_BORDER,
          color: LABELED_EDGE_CHIP_FOREGROUND,
        }}
      >
        {label}
      </span>
    </div>
  )
}

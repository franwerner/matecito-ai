import type { KeyboardEvent } from 'react'

import { cn } from '@/shared/lib/cn'
import { Button } from '@/shared/ui'

export type TimelineNodeState = 'done' | 'active' | 'pending'

export interface TimelineNode {
  state: TimelineNodeState
}

interface TimelineNodeStateConfig {
  readonly fillClassName: string
  readonly ringClassName: string
}

// Closed set: read TIMELINE_NODE_STATES[node.state], never compare against the string literal (structure/code-conventions.md).
const TIMELINE_NODE_STATES: Record<TimelineNodeState, TimelineNodeStateConfig> = {
  done: {
    fillClassName: 'bg-timeline-done',
    ringClassName: 'shadow-[0_0_0_1px_var(--timeline-done-border)]',
  },
  active: {
    fillClassName: 'bg-timeline-active',
    ringClassName: 'shadow-[0_0_0_1px_var(--timeline-active-border)]',
  },
  pending: {
    fillClassName: 'bg-timeline-pending',
    ringClassName: 'shadow-[0_0_0_1px_var(--timeline-pending-border)]',
  },
}

const PROGRESS_MIN = 0
const PROGRESS_MAX = 100
const KEYBOARD_STEP = 1
const SINGLE_NODE_COUNT = 1
const MIDPOINT_LEFT_PERCENT = 50

export interface TimelineScrubberProps {
  nodes: TimelineNode[]
  progress: number
  elapsed: string
  total: string
  isPlaying?: boolean
  onPlayToggle?: () => void
  onProgressChange?: (progress: number) => void
}

export function TimelineScrubber({
  nodes,
  progress,
  elapsed,
  total,
  isPlaying = false,
  onPlayToggle,
  onProgressChange,
}: TimelineScrubberProps) {
  const nodeCount = nodes.length

  const handlePlayheadKeyDown = (event: KeyboardEvent<HTMLDivElement>) => {
    if (!onProgressChange) return

    if (event.key === 'ArrowLeft' || event.key === 'ArrowDown') {
      event.preventDefault()
      onProgressChange(Math.max(PROGRESS_MIN, progress - KEYBOARD_STEP))
    } else if (event.key === 'ArrowRight' || event.key === 'ArrowUp') {
      event.preventDefault()
      onProgressChange(Math.min(PROGRESS_MAX, progress + KEYBOARD_STEP))
    } else if (event.key === 'Home') {
      event.preventDefault()
      onProgressChange(PROGRESS_MIN)
    } else if (event.key === 'End') {
      event.preventDefault()
      onProgressChange(PROGRESS_MAX)
    }
  }

  return (
    <div className="flex w-full items-center gap-[18px] border-t border-border bg-card px-5 py-[13px] font-sans">
      <div className="flex shrink-0 items-center gap-[10px]">
        <Button
          type="button"
          variant="secondary"
          size="icon"
          onClick={onPlayToggle}
          aria-label={isPlaying ? 'Pausar corrida' : 'Reproducir corrida'}
          className="border-secondary-border size-[34px] rounded-[10px] border p-0 text-[13px]"
        >
          <span aria-hidden="true">{isPlaying ? '❚❚' : '▶'}</span>
        </Button>
        <span className="text-chrome-body font-display text-[12.5px] font-semibold whitespace-nowrap">
          Reproducir corrida
        </span>
      </div>
      <div className="relative flex h-[34px] flex-1 items-center">
        <div className="absolute inset-x-0 h-[3px] rounded-full bg-border" />
        <div
          className="absolute left-0 h-[3px] rounded-full bg-gradient-to-r from-primary to-accent-foreground"
          style={{ width: `${progress}%` }}
        />
        {nodes.map((node, index) => {
          const state = TIMELINE_NODE_STATES[node.state]
          const left =
            nodeCount > SINGLE_NODE_COUNT
              ? (index / (nodeCount - SINGLE_NODE_COUNT)) * PROGRESS_MAX
              : MIDPOINT_LEFT_PERCENT
          return (
            <span
              key={index}
              aria-hidden="true"
              className={cn(
                'absolute size-[11px] -translate-x-1/2 -translate-y-1/2 rounded-full border-2 border-card',
                state.fillClassName,
                state.ringClassName,
              )}
              style={{ left: `${left}%`, top: '50%' }}
            />
          )
        })}
        <div
          role="slider"
          tabIndex={0}
          aria-label="Posición de la corrida"
          aria-valuemin={PROGRESS_MIN}
          aria-valuemax={PROGRESS_MAX}
          aria-valuenow={progress}
          onKeyDown={handlePlayheadKeyDown}
          className="border-primary absolute size-[15px] -translate-x-1/2 -translate-y-1/2 rounded-full border-[3px] bg-card shadow-[0_2px_8px_rgba(75,84,212,.4)] outline-offset-2"
          style={{ left: `${progress}%`, top: '50%' }}
        />
      </div>
      <span className="text-chrome-meta shrink-0 font-mono text-[11.5px] whitespace-nowrap">
        {elapsed} / {total}
      </span>
    </div>
  )
}

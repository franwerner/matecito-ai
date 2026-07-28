import { cn } from '@/shared/lib/cn'
import { Badge, Separator } from '@/shared/ui'

export interface CockpitHeaderProps {
  runId?: string
  live: boolean
}

type LiveState = 'live' | 'idle'

interface LiveStateConfig {
  readonly label: string
  readonly dotClassName: string
  readonly badgeClassName: string
}

// Closed set derived from the `live` boolean — read LIVE_STATES[state], never compare `live` against a string literal (structure/code-conventions.md).
const LIVE_STATES: Record<LiveState, LiveStateConfig> = {
  live: {
    label: 'En vivo',
    dotClassName: 'bg-agent-running',
    badgeClassName: 'border-primary/35 bg-primary/10 text-accent-foreground',
  },
  idle: {
    label: 'Detenido',
    dotClassName: 'bg-muted-foreground',
    badgeClassName: 'border-border bg-muted text-chrome-body',
  },
}

const RUN_ID_PLACEHOLDER = 'sin corrida activa'

export function CockpitHeader({ runId, live }: CockpitHeaderProps) {
  const state: LiveState = live ? 'live' : 'idle'
  const config = LIVE_STATES[state]

  return (
    <header className="flex h-14 shrink-0 items-center justify-between gap-4 border-b border-border bg-card px-[18px]">
      <div className="flex items-center gap-3">
        <span
          aria-hidden="true"
          className="font-display inline-flex size-[30px] items-center justify-center rounded-[9px] bg-gradient-to-br from-secondary to-secondary-border text-base font-bold text-secondary-foreground"
        >
          m
        </span>
        <span className="font-display text-lg font-semibold">Matecito</span>
        <Separator orientation="vertical" className="h-[22px]" />
        <span className="text-chrome-body font-mono text-xs">
          corrida <span className="text-foreground">{runId ?? RUN_ID_PLACEHOLDER}</span>
        </span>
        <Badge
          variant="outline"
          className={cn('gap-1.5 px-2.5 py-[3px] text-[11px]', config.badgeClassName)}
        >
          <span aria-hidden="true" className={cn('size-1.5 rounded-full', config.dotClassName)} />
          {config.label}
        </Badge>
      </div>
      {/* Command-palette hint: visual-only, no palette exists yet — a <span> stays out of the tab order and has no handler (frontend/routing.md, out of scope). */}
      <span className="text-chrome-meta inline-flex items-center gap-2 rounded-lg border border-border bg-muted px-2.5 py-[5px] font-mono text-xs">
        <span className="text-chrome-strong">⌘K</span> comandos
      </span>
    </header>
  )
}

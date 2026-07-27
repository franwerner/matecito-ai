import { cn } from '@/shared/lib/cn'
import { Badge } from '@/shared/ui'

export type AgentStatus = 'running' | 'done' | 'waiting'

interface AgentStatusConfig {
  readonly label: string
  readonly dotClassName: string
  readonly borderClassName: string
  readonly glowClassName: string
  readonly showActivityIndicator: boolean
}

// Closed set: read AGENT_STATUSES[status], never compare against the string literal (structure/code-conventions.md).
const AGENT_STATUSES: Record<AgentStatus, AgentStatusConfig> = {
  running: {
    label: 'En curso',
    dotClassName: 'bg-agent-running',
    borderClassName: 'border-agent-running-border',
    glowClassName: 'shadow-[0_0_0_1px_rgba(75,84,212,.35),0_10px_30px_-14px_rgba(75,84,212,.5)]',
    showActivityIndicator: true,
  },
  done: {
    label: 'Completado',
    dotClassName: 'bg-agent-done',
    borderClassName: 'border-agent-done-border',
    glowClassName: 'shadow-[0_6px_20px_-14px_rgba(30,45,70,.16)]',
    showActivityIndicator: false,
  },
  waiting: {
    label: 'En espera',
    dotClassName: 'bg-agent-waiting',
    borderClassName: 'border-agent-waiting-border',
    glowClassName: 'shadow-[0_6px_20px_-14px_rgba(30,45,70,.16)]',
    showActivityIndicator: false,
  },
}

export interface AgentNodeProps {
  type: string
  title: string
  phase: string
  status: AgentStatus
  mandate?: string
  tools?: string[]
}

export function AgentNode({ type, title, phase, status, mandate, tools }: AgentNodeProps) {
  const config = AGENT_STATUSES[status]
  const hasTools = Boolean(tools?.length)

  return (
    <div
      className={cn(
        'w-full rounded-[14px] border-[1.5px] bg-card p-[13px_14px] font-sans',
        config.borderClassName,
        config.glowClassName,
      )}
    >
      <div className="mb-[10px] flex items-center justify-between gap-2">
        <Badge
          variant="outline"
          className="gap-1.5 border-primary/40 bg-primary/[0.11] px-[9px] py-[3px] text-[11px] tracking-[.02em] text-accent-foreground"
        >
          <span aria-hidden="true" className="size-1.5 rounded-full bg-primary" />
          {type}
        </Badge>
        <span className="text-chrome-body inline-flex items-center gap-1.5 text-[11px] font-bold">
          <span aria-hidden="true" className={cn('size-2 rounded-full', config.dotClassName)} />
          {config.label}
        </span>
      </div>
      <div className="font-display text-[17px] leading-[1.15] font-semibold text-card-foreground">
        {title}
      </div>
      <div className="text-chrome-meta mt-1 font-mono text-[10.5px] tracking-[.06em] uppercase">
        {phase}
      </div>
      {mandate ? (
        <div className="text-chrome-body mt-[9px] text-[12.5px] leading-[1.45]">{mandate}</div>
      ) : null}
      {hasTools ? (
        <div className="mt-[11px] flex flex-wrap gap-[5px] border-t border-border pt-[11px]">
          {tools?.map((tool) => (
            <span
              key={tool}
              className="text-chrome-strong rounded-[5px] border border-border bg-muted px-[7px] py-0.5 font-mono text-[10.5px]"
            >
              {tool}
            </span>
          ))}
        </div>
      ) : null}
      {config.showActivityIndicator ? (
        <div className="mt-[11px] flex items-center gap-[7px]">
          <span
            aria-hidden="true"
            className="animate-mc-pulse motion-reduce:animate-none size-[7px] rounded-full bg-agent-running"
          />
          <span className="font-mono text-[10.5px] text-agent-running">escribiendo diff…</span>
        </div>
      ) : null}
    </div>
  )
}

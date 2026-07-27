import { cn } from '@/shared/lib/cn'
import { Badge } from '@/shared/ui'

export type DecisionStatus = 'borrador' | 'ratificada'

interface DecisionStatusConfig {
  readonly label: string
  readonly pillClassName: string
}

// Closed set: read DECISION_STATUSES[status], never compare against the string literal (structure/code-conventions.md).
const DECISION_STATUSES: Record<DecisionStatus, DecisionStatusConfig> = {
  borrador: {
    label: 'Borrador',
    pillClassName:
      'border-decision-borrador-border bg-decision-borrador text-decision-borrador-foreground',
  },
  ratificada: {
    label: 'Ratificada',
    pillClassName:
      'border-decision-ratificada-border bg-decision-ratificada text-decision-ratificada-foreground',
  },
}

export interface DecisionCardProps {
  id: string
  title: string
  status: DecisionStatus
  summary?: string
  relations: number
}

export function DecisionCard({ id, title, status, summary, relations }: DecisionCardProps) {
  const config = DECISION_STATUSES[status]

  return (
    <div className="w-full rounded-[10px] border border-border bg-card p-[14px_15px] font-sans">
      <div className="mb-2 flex items-center justify-between gap-[10px]">
        <span className="font-mono text-[10.5px] tracking-[.03em] text-[#95A0AE]">{id}</span>
        <Badge
          variant="outline"
          className={cn('gap-1.5 px-[9px] py-[3px] text-[11px]', config.pillClassName)}
        >
          <span aria-hidden="true" className="size-1.5 rounded-full bg-current" />
          {config.label}
        </Badge>
      </div>
      <div className="font-display text-[16px] leading-[1.2] font-semibold text-card-foreground">
        {title}
      </div>
      {summary ? (
        <div className="mt-[7px] text-[12.5px] leading-[1.5] text-[#7A8694]">{summary}</div>
      ) : null}
      <div className="mt-3 flex items-center gap-2 border-t border-border pt-[11px]">
        <span className="inline-flex items-center gap-1.5 font-mono text-[11px] text-[#5A6675]">
          <span aria-hidden="true" className="text-primary">
            ◇
          </span>
          {relations} piezas de código cumplen
        </span>
      </div>
    </div>
  )
}

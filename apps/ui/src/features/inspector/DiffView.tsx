import { cn } from '@/shared/lib/cn'
import { Badge } from '@/shared/ui'

import { type SeverityLevel, SeverityTag } from './SeverityTag'

export type DiffLineType = 'add' | 'del' | 'ctx'

interface DiffLineStyleConfig {
  readonly backgroundClassName: string
  readonly barClassName: string
  readonly signClassName: string
  readonly sign: string
}

// Closed set: read DIFF_LINE_STYLES[type], never compare against the string literal (structure/code-conventions.md).
const DIFF_LINE_STYLES: Record<DiffLineType, DiffLineStyleConfig> = {
  add: {
    backgroundClassName: 'bg-diff-add',
    barClassName: 'border-diff-add-border',
    signClassName: 'text-diff-add-foreground',
    sign: '+',
  },
  del: {
    backgroundClassName: 'bg-diff-del',
    barClassName: 'border-diff-del-border',
    signClassName: 'text-diff-del-foreground',
    sign: '-',
  },
  ctx: {
    backgroundClassName: 'bg-diff-ctx',
    barClassName: 'border-diff-ctx-border',
    signClassName: 'text-diff-ctx-foreground',
    sign: ' ',
  },
}

export interface DiffViewLine {
  type: DiffLineType
  lineNumber: string
  text: string
}

export interface DiffViewTag {
  level: SeverityLevel
  text?: string
}

export interface DiffViewProps {
  filename: string
  lines: DiffViewLine[]
  tags?: DiffViewTag[]
  justification: string
  decisionRef?: string
}

export function DiffView({ filename, lines, tags, justification, decisionRef }: DiffViewProps) {
  const hasTags = Boolean(tags?.length)

  return (
    <div className="w-full overflow-hidden rounded-[10px] border border-border bg-card font-sans">
      <div className="flex flex-wrap items-center justify-between gap-[10px] border-b border-border bg-muted px-3 py-[9px]">
        <span className="font-mono text-xs font-medium text-card-foreground">{filename}</span>
        {hasTags ? (
          <div className="flex flex-wrap gap-1.5">
            {tags?.map((tag, index) => (
              <SeverityTag key={`${tag.level}-${index}`} level={tag.level} text={tag.text} />
            ))}
          </div>
        ) : null}
      </div>
      <div className="overflow-x-auto bg-surface-sunken py-1.5">
        <div className="min-w-max">
          {lines.map((line, index) => {
            const style = DIFF_LINE_STYLES[line.type]
            return (
              <div
                key={index}
                className={cn(
                  'flex border-l-2 font-mono text-xs leading-[1.75]',
                  style.backgroundClassName,
                  style.barClassName,
                )}
              >
                <span className="text-diff-ctx-foreground w-9 shrink-0 pr-[9px] text-right select-none">
                  {line.lineNumber}
                </span>
                <span className={cn('w-[15px] shrink-0 select-none', style.signClassName)}>
                  {style.sign}
                </span>
                <span className="text-chrome-value pr-2.5 whitespace-pre">{line.text}</span>
              </div>
            )
          })}
        </div>
      </div>
      <div className="border-t border-border bg-card px-3 py-2.5">
        <div className="text-chrome-meta mb-1 text-[10.5px] font-extrabold tracking-[.06em] uppercase">
          Justificación
        </div>
        <div className="text-chrome-body text-[12.5px] leading-[1.5]">{justification}</div>
        {decisionRef ? (
          <div className="mt-2.5 flex">
            <Badge
              variant="outline"
              className="cursor-pointer gap-1.5 border-primary/35 bg-primary/10 px-[9px] py-0.5 font-mono text-[11px] text-accent-foreground"
            >
              <span aria-hidden="true">↩</span> remonta a {decisionRef}
            </Badge>
          </div>
        ) : null}
      </div>
    </div>
  )
}

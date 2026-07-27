import type { ComponentProps } from 'react'

import { cn } from '@/shared/lib/cn'

export type SeverityLevel = 'critical' | 'warning' | 'suggestion' | 'success'

interface SeverityLevelConfig {
  readonly label: string
  readonly glyph: string
  readonly chipClassName: string
  readonly iconClassName: string
}

// Closed set: read SEVERITY_LEVELS[level], never compare against the string literal (structure/code-conventions.md).
const SEVERITY_LEVELS: Record<SeverityLevel, SeverityLevelConfig> = {
  critical: {
    label: 'Crítico',
    glyph: '!',
    chipClassName:
      'border-severity-critical-border bg-severity-critical text-severity-critical-foreground',
    iconClassName: 'bg-severity-critical-border',
  },
  warning: {
    label: 'Advertencia',
    glyph: '△',
    chipClassName:
      'border-severity-warning-border bg-severity-warning text-severity-warning-foreground',
    iconClassName: 'bg-severity-warning-border',
  },
  suggestion: {
    label: 'Sugerencia',
    glyph: 'i',
    chipClassName:
      'border-severity-suggestion-border bg-severity-suggestion text-severity-suggestion-foreground',
    iconClassName: 'bg-severity-suggestion-border',
  },
  success: {
    label: 'Cumple',
    glyph: '✓',
    chipClassName:
      'border-severity-success-border bg-severity-success text-severity-success-foreground',
    iconClassName: 'bg-severity-success-border',
  },
}

export interface SeverityTagProps extends Omit<ComponentProps<'span'>, 'children'> {
  level: SeverityLevel
  text?: string
}

export function SeverityTag({ level, text, className, ...props }: SeverityTagProps) {
  const config = SEVERITY_LEVELS[level]

  return (
    <span
      data-slot="severity-tag"
      className={cn(
        'inline-flex items-center gap-[7px] rounded-[6px] border py-1 pr-2.5 pl-[5px] font-sans text-xs leading-[1.2] font-bold',
        config.chipClassName,
        className,
      )}
      {...props}
    >
      <span
        aria-hidden="true"
        className={cn(
          'text-background inline-flex size-4 items-center justify-center rounded-[4px] font-mono text-[10px] font-extrabold',
          config.iconClassName,
        )}
      >
        {config.glyph}
      </span>
      <span>{config.label}</span>
      {text ? <span className="text-[11.5px] font-semibold opacity-[.82]">· {text}</span> : null}
    </span>
  )
}

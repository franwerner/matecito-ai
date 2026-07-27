import { Fragment } from 'react'

import { Badge } from '@/shared/ui'

export interface ArtifactCardField {
  label: string
  value: string
}

export interface ArtifactCardProps {
  glyph: string
  type: string
  title: string
  fields: ArtifactCardField[]
  refs?: string[]
}

export function ArtifactCard({ glyph, type, title, fields, refs }: ArtifactCardProps) {
  const hasRefs = Boolean(refs?.length)

  return (
    <div className="w-full rounded-[10px] border border-border bg-card p-[14px_15px] font-sans">
      <div className="mb-[9px] flex items-center gap-2">
        <span
          aria-hidden="true"
          className="text-chrome-strong inline-flex size-6 items-center justify-center rounded-[6px] border border-border bg-muted font-mono text-xs"
        >
          {glyph}
        </span>
        <span className="text-chrome-meta font-mono text-[10.5px] font-bold tracking-[.06em] uppercase">
          {type}
        </span>
      </div>
      <div className="font-display text-base leading-[1.2] font-semibold text-card-foreground">
        {title}
      </div>
      <div className="mt-3 grid grid-cols-[auto_1fr] gap-x-[14px] gap-y-[6px]">
        {fields.map((field) => (
          <Fragment key={field.label}>
            <span className="text-chrome-meta font-mono text-[11px] whitespace-nowrap">
              {field.label}
            </span>
            <span className="text-chrome-value text-right font-mono text-[11.5px]">
              {field.value}
            </span>
          </Fragment>
        ))}
      </div>
      {hasRefs ? (
        <div className="mt-3 flex flex-wrap gap-[6px] border-t border-border pt-[11px]">
          {refs?.map((ref) => (
            <Badge
              key={ref}
              variant="outline"
              className="gap-[5px] border-primary/30 bg-primary/10 px-2 py-0.5 font-mono text-[10.5px] text-accent-foreground"
            >
              <span aria-hidden="true">↩</span> {ref}
            </Badge>
          ))}
        </div>
      ) : null}
    </div>
  )
}

import type { CSSProperties } from 'react'
import { Toaster as SonnerToaster, type ToasterProps } from 'sonner'

const TOASTER_TOKEN_STYLE = {
  '--normal-bg': 'var(--popover)',
  '--normal-text': 'var(--popover-foreground)',
  '--normal-border': 'var(--border)',
} as CSSProperties

export function Toaster(props: ToasterProps) {
  return (
    <SonnerToaster
      data-slot="toaster"
      theme="light"
      className="toaster group"
      style={TOASTER_TOKEN_STYLE}
      {...props}
    />
  )
}

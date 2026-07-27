import type { ComponentProps } from 'react'

import { cn } from '@/shared/lib/cn'

export function Input({ className, type, ...props }: ComponentProps<'input'>) {
  return (
    <input
      type={type}
      data-slot="input"
      className={cn(
        'flex h-9 w-full min-w-0 rounded-md border border-input bg-card px-3 py-1 text-sm text-foreground shadow-xs transition-colors outline-none placeholder:text-muted-foreground focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50',
        'aria-invalid:border-destructive aria-invalid:ring-destructive/40',
        className,
      )}
      {...props}
    />
  )
}

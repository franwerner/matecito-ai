import type { ComponentProps } from 'react'

import { cn } from '@/shared/lib/cn'

export function Card({ className, ...props }: ComponentProps<'div'>) {
  return (
    <div
      data-slot="card"
      className={cn(
        'flex flex-col gap-4 rounded-xl border border-border bg-card p-4 text-card-foreground shadow-sm',
        className,
      )}
      {...props}
    />
  )
}

export function CardHeader({ className, ...props }: ComponentProps<'div'>) {
  return (
    <div data-slot="card-header" className={cn('flex flex-col gap-1.5', className)} {...props} />
  )
}

export function CardTitle({ className, ...props }: ComponentProps<'div'>) {
  return (
    <div
      data-slot="card-title"
      className={cn('font-display text-base leading-none font-semibold', className)}
      {...props}
    />
  )
}

export function CardDescription({ className, ...props }: ComponentProps<'div'>) {
  return (
    <div
      data-slot="card-description"
      className={cn('text-sm text-muted-foreground', className)}
      {...props}
    />
  )
}

export function CardContent({ className, ...props }: ComponentProps<'div'>) {
  return <div data-slot="card-content" className={className} {...props} />
}

export function CardFooter({ className, ...props }: ComponentProps<'div'>) {
  return (
    <div
      data-slot="card-footer"
      className={cn('flex items-center gap-2 border-t border-border pt-4', className)}
      {...props}
    />
  )
}

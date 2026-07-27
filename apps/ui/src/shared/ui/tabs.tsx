import * as TabsPrimitive from '@radix-ui/react-tabs'
import type { ComponentProps } from 'react'

import { cn } from '@/shared/lib/cn'

export function Tabs({ className, ...props }: ComponentProps<typeof TabsPrimitive.Root>) {
  return (
    <TabsPrimitive.Root
      data-slot="tabs"
      className={cn('flex flex-col gap-2', className)}
      {...props}
    />
  )
}

export function TabsList({ className, ...props }: ComponentProps<typeof TabsPrimitive.List>) {
  return (
    <TabsPrimitive.List
      data-slot="tabs-list"
      className={cn('inline-flex items-center gap-1 border-b border-border', className)}
      {...props}
    />
  )
}

export function TabsTrigger({ className, ...props }: ComponentProps<typeof TabsPrimitive.Trigger>) {
  return (
    <TabsPrimitive.Trigger
      data-slot="tabs-trigger"
      className={cn(
        'inline-flex items-center justify-center border-b-2 border-transparent px-3 py-2 text-sm font-bold text-muted-foreground transition-colors outline-none focus-visible:ring-2 focus-visible:ring-ring data-[state=active]:border-primary data-[state=active]:text-foreground disabled:pointer-events-none disabled:opacity-50',
        className,
      )}
      {...props}
    />
  )
}

export function TabsContent({ className, ...props }: ComponentProps<typeof TabsPrimitive.Content>) {
  return (
    <TabsPrimitive.Content
      data-slot="tabs-content"
      className={cn('outline-none', className)}
      {...props}
    />
  )
}

import * as TooltipPrimitive from '@radix-ui/react-tooltip'
import type { ComponentProps } from 'react'

import { cn } from '@/shared/lib/cn'

export function TooltipProvider({
  delayDuration = 200,
  ...props
}: ComponentProps<typeof TooltipPrimitive.Provider>) {
  return (
    <TooltipPrimitive.Provider
      data-slot="tooltip-provider"
      delayDuration={delayDuration}
      {...props}
    />
  )
}

export function Tooltip(props: ComponentProps<typeof TooltipPrimitive.Root>) {
  return (
    <TooltipProvider>
      <TooltipPrimitive.Root data-slot="tooltip" {...props} />
    </TooltipProvider>
  )
}

export function TooltipTrigger(props: ComponentProps<typeof TooltipPrimitive.Trigger>) {
  return <TooltipPrimitive.Trigger data-slot="tooltip-trigger" {...props} />
}

export function TooltipContent({
  className,
  sideOffset = 6,
  children,
  ...props
}: ComponentProps<typeof TooltipPrimitive.Content>) {
  return (
    <TooltipPrimitive.Portal>
      <TooltipPrimitive.Content
        data-slot="tooltip-content"
        sideOffset={sideOffset}
        className={cn(
          'z-50 origin-[var(--radix-tooltip-content-transform-origin)] rounded-md border border-border bg-popover px-3 py-1.5 font-mono text-xs text-popover-foreground shadow-md transition-[opacity,transform] duration-150 data-[state=closed]:scale-95 data-[state=closed]:opacity-0 data-[state=delayed-open]:scale-100 data-[state=delayed-open]:opacity-100',
          className,
        )}
        {...props}
      >
        {children}
        <TooltipPrimitive.Arrow className="fill-popover" />
      </TooltipPrimitive.Content>
    </TooltipPrimitive.Portal>
  )
}

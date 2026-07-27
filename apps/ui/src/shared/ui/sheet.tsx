import * as DialogPrimitive from '@radix-ui/react-dialog'
import { cva, type VariantProps } from 'class-variance-authority'
import { X } from 'lucide-react'
import type { ComponentProps } from 'react'

import { cn } from '@/shared/lib/cn'

export function Sheet(props: ComponentProps<typeof DialogPrimitive.Root>) {
  return <DialogPrimitive.Root data-slot="sheet" {...props} />
}

export function SheetTrigger(props: ComponentProps<typeof DialogPrimitive.Trigger>) {
  return <DialogPrimitive.Trigger data-slot="sheet-trigger" {...props} />
}

export function SheetClose(props: ComponentProps<typeof DialogPrimitive.Close>) {
  return <DialogPrimitive.Close data-slot="sheet-close" {...props} />
}

function SheetOverlay({ className, ...props }: ComponentProps<typeof DialogPrimitive.Overlay>) {
  return (
    <DialogPrimitive.Overlay
      data-slot="sheet-overlay"
      className={cn(
        'fixed inset-0 z-50 bg-foreground/40 transition-opacity duration-200 data-[state=closed]:opacity-0 data-[state=open]:opacity-100',
        className,
      )}
      {...props}
    />
  )
}

const sheetVariants = cva(
  'fixed z-50 flex flex-col gap-4 border-border bg-card p-6 text-card-foreground shadow-lg transition-transform duration-200 ease-out',
  {
    variants: {
      side: {
        top: 'inset-x-0 top-0 border-b data-[state=closed]:-translate-y-full data-[state=open]:translate-y-0',
        bottom:
          'inset-x-0 bottom-0 border-t data-[state=closed]:translate-y-full data-[state=open]:translate-y-0',
        left: 'inset-y-0 left-0 h-full w-3/4 border-r data-[state=closed]:-translate-x-full data-[state=open]:translate-x-0 sm:max-w-sm',
        right:
          'inset-y-0 right-0 h-full w-3/4 border-l data-[state=closed]:translate-x-full data-[state=open]:translate-x-0 sm:max-w-sm',
      },
    },
    defaultVariants: {
      side: 'right',
    },
  },
)

export interface SheetContentProps
  extends ComponentProps<typeof DialogPrimitive.Content>, VariantProps<typeof sheetVariants> {}

export function SheetContent({ className, children, side, ...props }: SheetContentProps) {
  return (
    <DialogPrimitive.Portal>
      <SheetOverlay />
      <DialogPrimitive.Content
        data-slot="sheet-content"
        className={cn(sheetVariants({ side }), className)}
        {...props}
      >
        {children}
        <DialogPrimitive.Close
          data-slot="sheet-close"
          className="absolute top-4 right-4 rounded-md text-muted-foreground transition-colors outline-none hover:text-foreground focus-visible:ring-2 focus-visible:ring-ring"
        >
          <X className="size-4" />
          <span className="sr-only">Cerrar</span>
        </DialogPrimitive.Close>
      </DialogPrimitive.Content>
    </DialogPrimitive.Portal>
  )
}

export function SheetHeader({ className, ...props }: ComponentProps<'div'>) {
  return (
    <div data-slot="sheet-header" className={cn('flex flex-col gap-1.5', className)} {...props} />
  )
}

export function SheetFooter({ className, ...props }: ComponentProps<'div'>) {
  return (
    <div
      data-slot="sheet-footer"
      className={cn('mt-auto flex flex-col gap-2', className)}
      {...props}
    />
  )
}

export function SheetTitle({ className, ...props }: ComponentProps<typeof DialogPrimitive.Title>) {
  return (
    <DialogPrimitive.Title
      data-slot="sheet-title"
      className={cn('font-display text-lg font-semibold', className)}
      {...props}
    />
  )
}

export function SheetDescription({
  className,
  ...props
}: ComponentProps<typeof DialogPrimitive.Description>) {
  return (
    <DialogPrimitive.Description
      data-slot="sheet-description"
      className={cn('text-sm text-muted-foreground', className)}
      {...props}
    />
  )
}

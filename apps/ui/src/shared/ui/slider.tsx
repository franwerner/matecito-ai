import * as SliderPrimitive from '@radix-ui/react-slider'
import type { ComponentProps } from 'react'
import { useMemo } from 'react'

import { cn } from '@/shared/lib/cn'

export function Slider({
  className,
  defaultValue,
  value,
  min = 0,
  max = 100,
  ...props
}: ComponentProps<typeof SliderPrimitive.Root>) {
  const values = useMemo(
    () => (Array.isArray(value) ? value : Array.isArray(defaultValue) ? defaultValue : [min, max]),
    [value, defaultValue, min, max],
  )

  return (
    <SliderPrimitive.Root
      data-slot="slider"
      defaultValue={defaultValue}
      value={value}
      min={min}
      max={max}
      className={cn(
        'relative flex w-full touch-none items-center select-none data-[disabled]:opacity-50 data-[orientation=vertical]:h-full data-[orientation=vertical]:min-h-44 data-[orientation=vertical]:w-auto data-[orientation=vertical]:flex-col',
        className,
      )}
      {...props}
    >
      <SliderPrimitive.Track
        data-slot="slider-track"
        className="relative grow overflow-hidden rounded-full bg-muted data-[orientation=horizontal]:h-1.5 data-[orientation=horizontal]:w-full data-[orientation=vertical]:h-full data-[orientation=vertical]:w-1.5"
      >
        <SliderPrimitive.Range
          data-slot="slider-range"
          className="absolute bg-primary data-[orientation=horizontal]:h-full data-[orientation=vertical]:w-full"
        />
      </SliderPrimitive.Track>
      {values.map((_, index) => (
        <SliderPrimitive.Thumb
          data-slot="slider-thumb"
          key={`slider-thumb-${index}`}
          className="block size-4 shrink-0 rounded-full border border-primary bg-card shadow-sm transition-colors outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50"
        />
      ))}
    </SliderPrimitive.Root>
  )
}

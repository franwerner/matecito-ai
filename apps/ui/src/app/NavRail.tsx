import {
  Diamond,
  GitCompare,
  LayoutGrid,
  ListChecks,
  type LucideIcon,
  Settings,
} from 'lucide-react'

import { cn } from '@/shared/lib/cn'

export type RailItemId = 'canvas' | 'diffs' | 'decisions' | 'queue' | 'settings'

interface RailItem {
  readonly id: RailItemId
  readonly label: string
  readonly icon: LucideIcon
  readonly anchorToFoot?: boolean
}

// Closed set: fixed order and content, rendered as a list — never rebuilt or reordered inline (structure/code-conventions.md).
const RAIL_ITEMS: readonly RailItem[] = [
  { id: 'canvas', label: 'Canvas', icon: LayoutGrid },
  { id: 'diffs', label: 'Diffs', icon: GitCompare },
  { id: 'decisions', label: 'Decisiones', icon: Diamond },
  { id: 'queue', label: 'Cola', icon: ListChecks },
  { id: 'settings', label: 'Ajustes', icon: Settings, anchorToFoot: true },
]

const DEFAULT_ICON_STROKE_WIDTH = 2
const ACTIVE_ICON_STROKE_WIDTH = 2.5

export interface NavRailProps {
  active: RailItemId
  onSelect: (id: RailItemId) => void
}

export function NavRail({ active, onSelect }: NavRailProps) {
  return (
    <nav
      aria-label="Navegación principal"
      className="flex h-full flex-col items-center gap-1.5 bg-[#16181C] py-3.5"
    >
      {RAIL_ITEMS.map((item) => {
        const isActive = item.id === active
        const Icon = item.icon

        return (
          <button
            key={item.id}
            type="button"
            onClick={() => onSelect(item.id)}
            aria-label={item.label}
            aria-current={isActive ? 'page' : undefined}
            className={cn(
              'relative inline-flex size-[38px] items-center justify-center rounded-[11px] transition-colors outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-[#16181C]',
              isActive
                ? 'bg-secondary text-secondary-foreground'
                : 'text-white/55 hover:bg-white/10',
              item.anchorToFoot && 'mt-auto',
            )}
          >
            {/* Shape-based active marker so the current item is never signaled by color alone (frontend/accessibility.md). */}
            {isActive ? (
              <span
                aria-hidden="true"
                className="absolute -left-[7px] h-4 w-[3px] rounded-full bg-secondary"
              />
            ) : null}
            <Icon
              aria-hidden
              className="size-[18px]"
              strokeWidth={isActive ? ACTIVE_ICON_STROKE_WIDTH : DEFAULT_ICON_STROKE_WIDTH}
            />
          </button>
        )
      })}
    </nav>
  )
}

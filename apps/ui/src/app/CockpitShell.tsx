import { type ReactNode, useState } from 'react'

import { CockpitHeader } from './CockpitHeader'
import { NavRail, type RailItemId } from './NavRail'

export interface CockpitShellProps {
  children: ReactNode
}

const INITIAL_RAIL_ITEM: RailItemId = 'canvas'

// Run identity (runId/live) is not wired yet — routing and the broker connection are out of scope
// for this change (see intake `sdd/ui-cockpit-shell/intake`), so the shell renders the header's
// no-run state until a real run is threaded through.
export function CockpitShell({ children }: CockpitShellProps) {
  const [active, setActive] = useState<RailItemId>(INITIAL_RAIL_ITEM)

  return (
    <div className="grid h-dvh grid-cols-[60px_1fr] grid-rows-[56px_1fr] overflow-hidden bg-background">
      <div className="col-span-2">
        <CockpitHeader live={false} />
      </div>
      <NavRail active={active} onSelect={setActive} />
      <main className="min-h-0 overflow-y-auto">{children}</main>
    </div>
  )
}

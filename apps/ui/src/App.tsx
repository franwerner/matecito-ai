import { useState } from 'react'

import { CockpitShell } from '@/app'
import { AgentNode, type AgentStatus, CanvasEdge, LabeledEdge } from '@/features/canvas'
import {
  ArtifactCard,
  DecisionCard,
  DiffView,
  type DiffViewLine,
  type SeverityLevel,
} from '@/features/inspector'
import { type TimelineNode, TimelineScrubber } from '@/features/timeline'

interface DemoAgent {
  readonly type: string
  readonly title: string
  readonly phase: string
  readonly status: AgentStatus
  readonly mandate: string
  readonly tools?: string[]
}

// Static pipeline mock: no broker, no WebSocket — see Requirement "Composición demo del contenido".
const DEMO_AGENTS: readonly DemoAgent[] = [
  {
    type: 'Orquestador',
    title: 'Orquestador',
    phase: 'Raíz',
    status: 'done',
    mandate: 'Reparte el trabajo entre agentes.',
  },
  {
    type: 'Agente',
    title: 'Nodo 1',
    phase: 'Fase 01 · Análisis',
    status: 'done',
    mandate: 'Explora el repositorio y arma el contexto.',
  },
  {
    type: 'Agente',
    title: 'Nodo 2',
    phase: 'Fase 02 · Plan',
    status: 'done',
    mandate: 'Propone el plan y las decisiones a ratificar.',
  },
  {
    type: 'Agente',
    title: 'Nodo 3',
    phase: 'Fase 04 · Código',
    status: 'running',
    mandate: 'Implementa la validación de sesión.',
    tools: ['read', 'edit', 'bash'],
  },
  {
    type: 'Agente',
    title: 'Nodo 4',
    phase: 'Fase 06 · Tests',
    status: 'waiting',
    mandate: 'Cubre los cambios con pruebas.',
  },
  {
    type: 'Agente',
    title: 'Nodo 5',
    phase: 'Fase 07 · Revisión',
    status: 'waiting',
    mandate: 'Audita severidad y cumplimiento.',
  },
]

type DemoEdgeKind = 'canvas' | 'labeled'

interface DemoEdge {
  readonly kind: DemoEdgeKind
  readonly label: string
  readonly active?: boolean
}

// One entry per gap between consecutive DEMO_AGENTS (length - 1) — indexed by
// position instead of derived from it, so no modulo/magic-number arithmetic is needed.
const DEMO_EDGES: readonly DemoEdge[] = [
  { kind: 'canvas', label: 'contexto' },
  { kind: 'labeled', label: 'plan' },
  { kind: 'canvas', label: 'diffs' },
  { kind: 'labeled', label: 'plan', active: true },
  { kind: 'canvas', label: 'revisión' },
]

const DEMO_DIFF_LINES: DiffViewLine[] = [
  { type: 'ctx', lineNumber: '18', text: 'export async function loadSession(token: string) {' },
  { type: 'del', lineNumber: '19', text: '  return decodeToken(token)' },
  { type: 'add', lineNumber: '19', text: '  const claims = decodeToken(token)' },
  { type: 'add', lineNumber: '20', text: '  if (!claims) throw new SessionError()' },
  { type: 'add', lineNumber: '21', text: '  return claims' },
  { type: 'ctx', lineNumber: '22', text: '}' },
]

const DEMO_DIFF_TAGS: { level: SeverityLevel; text?: string }[] = [
  { level: 'critical', text: 'token sin validar' },
]

const DEMO_TIMELINE_NODES: TimelineNode[] = [
  { state: 'done' },
  { state: 'done' },
  { state: 'done' },
  { state: 'active' },
  { state: 'pending' },
  { state: 'pending' },
  { state: 'pending' },
]

const INITIAL_TIMELINE_PROGRESS = 46

export function App() {
  const [progress, setProgress] = useState(INITIAL_TIMELINE_PROGRESS)
  const [isPlaying, setIsPlaying] = useState(false)

  return (
    <CockpitShell>
      <div className="grid h-full grid-rows-[1fr_auto]">
        <div className="grid min-h-0 grid-cols-[1fr_400px]">
          <div className="flex min-h-0 items-center gap-3 overflow-auto bg-background p-8">
            {DEMO_AGENTS.map((agent, index) => {
              const edge = DEMO_EDGES[index]
              return (
                <div key={agent.title} className="flex shrink-0 items-center gap-3">
                  <div className="w-[236px]">
                    <AgentNode
                      type={agent.type}
                      title={agent.title}
                      phase={agent.phase}
                      status={agent.status}
                      mandate={agent.mandate}
                      tools={agent.tools}
                    />
                  </div>
                  {edge ? (
                    <div className="w-[150px]">
                      {edge.kind === 'canvas' ? (
                        <CanvasEdge label={edge.label} />
                      ) : (
                        <LabeledEdge label={edge.label} active={edge.active} />
                      )}
                    </div>
                  ) : null}
                </div>
              )
            })}
          </div>

          <aside className="flex min-h-0 flex-col gap-[18px] overflow-y-auto border-l border-border bg-card p-[18px]">
            <section>
              <div className="text-chrome-meta mb-[9px] font-mono text-[10.5px] font-bold tracking-[.06em] uppercase">
                Diff escrito
              </div>
              <DiffView
                filename="src/auth/session.ts"
                lines={DEMO_DIFF_LINES}
                tags={DEMO_DIFF_TAGS}
                justification="La sesión debe invalidarse si el token no trae claims válidas."
                decisionRef="DEC-014"
              />
            </section>
            <section>
              <div className="text-chrome-meta mb-[9px] font-mono text-[10.5px] font-bold tracking-[.06em] uppercase">
                Decisión aplicada
              </div>
              <DecisionCard
                id="DEC-014"
                title="Autenticación por token de sesión"
                status="ratificada"
                summary="Toda ruta protegida valida claims antes de resolver el handler."
                relations={3}
              />
            </section>
            <section>
              <div className="text-chrome-meta mb-[9px] font-mono text-[10.5px] font-bold tracking-[.06em] uppercase">
                Artefacto producido
              </div>
              <ArtifactCard
                glyph="◆"
                type="Módulo"
                title="session.ts"
                fields={[
                  { label: 'Líneas', value: '+3 / −1' },
                  { label: 'Cobertura', value: '92%' },
                ]}
                refs={['DEC-014']}
              />
            </section>
          </aside>
        </div>

        <div role="region" aria-label="Línea de tiempo de la corrida">
          <TimelineScrubber
            nodes={DEMO_TIMELINE_NODES}
            progress={progress}
            elapsed="02:14"
            total="04:30"
            isPlaying={isPlaying}
            onPlayToggle={() => setIsPlaying((playing) => !playing)}
            onProgressChange={setProgress}
          />
        </div>
      </div>
    </CockpitShell>
  )
}

"use client";

import {
  DndContext,
  type DragEndEvent,
  DragOverlay,
  type DragStartEvent,
  KeyboardSensor,
  PointerSensor,
  TouchSensor,
  useDraggable,
  useDroppable,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import { CSS } from "@dnd-kit/utilities";
import { GripVerticalIcon, Loader2Icon } from "lucide-react";
import Link from "next/link";
import { useCallback, useEffect, useMemo, useState } from "react";

import { DashboardShell } from "@/components/layout/dashboard-shell";
import { showToast } from "@/components/shared/toast";
import { Button } from "@/components/ui/button";
import { DemandStatusBadge } from "@/features/demands/components/demand-status-badge";
import {
  DemandTransitionDialog,
  type TimelinePayload,
} from "@/features/demands/components/demand-transition-dialog";
import {
  apiDemandAction,
  apiGetManagedOffice,
  apiListOfficeDemands,
  type Demand,
  type DemandStatus,
} from "@/lib/api";
import { apiErrorMessage } from "@/lib/forms";
import { cn } from "@/lib/utils";

type BoardAction = "claim" | "start" | "request-confirmation";
type PendingTransition = { demand: Demand; action: BoardAction };

const columns: DemandStatus[] = [
  "registered",
  "under_review",
  "in_progress",
  "awaiting_confirmation",
  "completed",
];
const boardColumnLabel: Record<DemandStatus, string> = {
  registered: "Registradas",
  under_review: "Em análise",
  in_progress: "Em execução",
  awaiting_confirmation: "Aguardando validação",
  completed: "Concluídas",
};

const transitionCopy: Record<
  BoardAction,
  { title: string; description: string; submitLabel: string }
> = {
  claim: {
    title: "Assumir atendimento",
    description: "Registre como o gabinete iniciará a análise desta demanda.",
    submitLabel: "Assumir demanda",
  },
  start: {
    title: "Iniciar execução",
    description: "Informe qual ação prática passou a ser executada pelo gabinete.",
    submitLabel: "Iniciar execução",
  },
  "request-confirmation": {
    title: "Solicitar validação",
    description: "Explique a solução entregue para que a pessoa autora possa validar o resultado.",
    submitLabel: "Solicitar validação",
  },
};

function actionForTransition(from: DemandStatus, to: DemandStatus): BoardAction | undefined {
  if (from === "registered" && to === "under_review") return "claim";
  if (from === "under_review" && to === "in_progress") return "start";
  if (from === "in_progress" && to === "awaiting_confirmation") return "request-confirmation";
  return undefined;
}

export default function OfficeDemandBoardPage() {
  const [demands, setDemands] = useState<Demand[]>([]);
  const [officeID, setOfficeID] = useState<number>();
  const [officeSlug, setOfficeSlug] = useState<string>();
  const [loading, setLoading] = useState(true);
  const [activeDemand, setActiveDemand] = useState<Demand>();
  const [pending, setPending] = useState<PendingTransition>();

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 8 } }),
    useSensor(TouchSensor, { activationConstraint: { delay: 180, tolerance: 8 } }),
    useSensor(KeyboardSensor),
  );

  const load = useCallback(async () => {
    setLoading(true);
    const [officeResult, demandResult] = await Promise.all([
      apiGetManagedOffice(),
      apiListOfficeDemands(),
    ]);
    if (officeResult.ok) {
      setOfficeSlug(officeResult.data?.slug);
      setOfficeID(officeResult.data?.office_id);
    }
    if (!demandResult.ok) {
      showToast(
        apiErrorMessage(demandResult.error, "Não foi possível carregar o acompanhamento."),
        "error",
      );
    } else {
      setDemands(demandResult.data ?? []);
    }
    setLoading(false);
  }, []);

  useEffect(() => {
    const timeout = window.setTimeout(() => void load(), 0);
    return () => window.clearTimeout(timeout);
  }, [load]);

  const demandsByStatus = useMemo(
    () =>
      Object.fromEntries(
        columns.map((status) => [
          status,
          demands
            .filter((demand) => demand.status === status)
            .sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime()),
        ]),
      ) as Record<DemandStatus, Demand[]>,
    [demands],
  );

  function onDragStart(event: DragStartEvent) {
    const id = Number(String(event.active.id).replace("demand-", ""));
    setActiveDemand(demands.find((demand) => demand.id === id));
  }

  function onDragEnd(event: DragEndEvent) {
    setActiveDemand(undefined);
    if (!event.over) return;
    const id = Number(String(event.active.id).replace("demand-", ""));
    const demand = demands.find((item) => item.id === id);
    const destination = event.over.id as DemandStatus;
    if (!demand || demand.status === destination) return;
    const action = actionForTransition(demand.status, destination);
    if (!action) {
      showToast("Esta demanda só pode avançar para a próxima situação autorizada.", "error");
      return;
    }
    setPending({ demand, action });
  }

  async function submitTransition(payload: TimelinePayload) {
    if (!pending) return false;
    const result = await apiDemandAction(pending.demand.id, pending.action, payload);
    if (!result.ok || !result.data) {
      showToast(apiErrorMessage(result.error, "Não foi possível mover a demanda."), "error");
      return false;
    }
    setDemands((current) =>
      current.map((item) => (item.id === result.data?.id ? result.data : item)),
    );
    showToast("Situação atualizada e registrada na linha do tempo.");
    return true;
  }

  return (
    <DashboardShell
      title="Acompanhamento"
      subtitle="Organize as demandas do gabinete por situação, com cada mudança registrada."
      officeSlug={officeSlug}
    >
      <div className="mb-6 flex flex-wrap items-end justify-between gap-3">
        <div>
          <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">
            Fluxo operacional
          </p>
          <h2 className="font-display mt-1 text-2xl font-semibold">Board de demandas</h2>
          <p className="text-ink-soft mt-2 max-w-2xl text-sm leading-6">
            Arraste uma demanda para a próxima situação. Antes de confirmar, uma justificativa é
            obrigatória.
          </p>
        </div>
        <Button variant="outline" onClick={() => void load()} disabled={loading}>
          {loading ? <Loader2Icon className="animate-spin" /> : null}
          Atualizar
        </Button>
      </div>
      {loading ? (
        <div className="text-ink-soft grid min-h-[45vh] place-items-center text-sm">
          Carregando board…
        </div>
      ) : (
        <DndContext
          sensors={sensors}
          onDragStart={onDragStart}
          onDragEnd={onDragEnd}
          onDragCancel={() => setActiveDemand(undefined)}
        >
          <div className="overflow-x-auto pb-2">
            <div className="grid min-w-[72rem] grid-cols-5 gap-4" aria-label="Board de demandas">
              {columns.map((status) => (
                <BoardColumn
                  key={status}
                  status={status}
                  demands={demandsByStatus[status]}
                  officeID={officeID}
                />
              ))}
            </div>
          </div>
          <DragOverlay dropAnimation={null}>
            {activeDemand ? <DemandBoardCardPreview demand={activeDemand} /> : null}
          </DragOverlay>
        </DndContext>
      )}
      <p className="text-ink-faint mt-3 text-xs lg:hidden">
        Deslize horizontalmente para ver todas as situações do board.
      </p>
      {pending ? (
        <DemandTransitionDialog
          open
          title={transitionCopy[pending.action].title}
          description={transitionCopy[pending.action].description}
          submitLabel={transitionCopy[pending.action].submitLabel}
          onClose={() => setPending(undefined)}
          onSubmit={submitTransition}
        />
      ) : null}
    </DashboardShell>
  );
}

function BoardColumn({
  status,
  demands,
  officeID,
}: {
  status: DemandStatus;
  demands: Demand[];
  officeID?: number;
}) {
  const { isOver, setNodeRef } = useDroppable({ id: status });
  return (
    <section
      ref={setNodeRef}
      className={cn(
        "border-line bg-paper-2/60 min-h-[31rem] rounded-xl border p-3 transition",
        isOver && "ring-lime bg-lime-pale/45 ring-2",
      )}
      aria-label={boardColumnLabel[status]}
    >
      <header className="mb-3 flex items-center justify-between gap-2 px-1">
        <span className="text-sm font-semibold">{boardColumnLabel[status]}</span>
        <span className="bg-card text-ink-soft grid min-w-6 place-items-center rounded-full px-2 py-1 text-xs font-bold">
          {demands.length}
        </span>
      </header>
      <div className="space-y-3">
        {demands.map((demand) => (
          <DemandBoardCard
            key={demand.id}
            demand={demand}
            canMove={
              ["registered", "under_review", "in_progress"].includes(demand.status) &&
              (!demand.responsible_office_id || demand.responsible_office_id === officeID)
            }
          />
        ))}
      </div>
      {!demands.length ? (
        <p className="text-ink-faint px-2 py-8 text-center text-xs">Sem demandas</p>
      ) : null}
    </section>
  );
}

function DemandBoardCard({ demand, canMove }: { demand: Demand; canMove: boolean }) {
  const { attributes, listeners, setNodeRef, transform, isDragging } = useDraggable({
    id: `demand-${demand.id}`,
    disabled: !canMove,
  });
  const style = transform ? { transform: CSS.Translate.toString(transform) } : undefined;
  return (
    <article
      ref={setNodeRef}
      style={style}
      className={cn(
        "border-line bg-card rounded-lg border p-3 shadow-sm transition",
        isDragging && "opacity-35",
      )}
    >
      <div className="flex items-start gap-1.5">
        <button
          type="button"
          className={cn(
            "text-ink-faint hover:text-ink mt-0.5 -ml-1 grid size-6 shrink-0 place-items-center rounded",
            canMove ? "cursor-grab active:cursor-grabbing" : "cursor-not-allowed opacity-40",
          )}
          aria-label={
            canMove ? `Mover ${demand.title}` : "Demanda sob responsabilidade de outro gabinete"
          }
          disabled={!canMove}
          {...attributes}
          {...listeners}
        >
          <GripVerticalIcon className="size-4" />
        </button>
        <Link
          href={`/demandas/${demand.id}`}
          className="min-w-0 flex-1 rounded-sm focus-visible:outline-2"
        >
          <p className="text-ink-faint font-mono text-[10px] font-semibold">{demand.protocol}</p>
          <p className="mt-1 line-clamp-3 text-sm leading-5 font-semibold">{demand.title}</p>
          <p className="text-ink-soft mt-2 line-clamp-2 text-xs leading-5">{demand.neighborhood}</p>
          <DemandStatusBadge status={demand.status} className="mt-3" />
        </Link>
      </div>
    </article>
  );
}

function DemandBoardCardPreview({ demand }: { demand: Demand }) {
  return (
    <article className="border-line bg-card w-[13.5rem] rotate-1 rounded-lg border p-3 shadow-xl">
      <p className="text-ink-faint font-mono text-[10px] font-semibold">{demand.protocol}</p>
      <p className="mt-1 line-clamp-3 text-sm leading-5 font-semibold">{demand.title}</p>
      <DemandStatusBadge status={demand.status} className="mt-3" />
    </article>
  );
}

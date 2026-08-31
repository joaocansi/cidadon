import type { DemandStatus } from "@/lib/api";
import { cn } from "@/lib/utils";

const STATUS_LABELS: Record<DemandStatus, string> = {
  registered: "Registrada",
  under_review: "Em análise",
  in_progress: "Em andamento",
  awaiting_confirmation: "Aguardando validação",
  completed: "Concluída",
};

const STATUS_STYLES: Record<DemandStatus, string> = {
  registered: "bg-paper-2 text-ink-soft ring-line",
  under_review: "bg-lime-pale text-lime-deep ring-lime/30",
  in_progress: "bg-amber/15 text-amber-deep ring-amber/30",
  awaiting_confirmation: "bg-sky-100 text-sky-800 ring-sky-200",
  completed: "bg-lime text-pine ring-lime-deep/20",
};

export function DemandStatusBadge({
  status,
  className,
}: {
  status: DemandStatus;
  className?: string;
}) {
  return (
    <span
      className={cn(
        "inline-flex min-h-7 items-center rounded-full px-3 text-xs font-semibold ring-1",
        STATUS_STYLES[status],
        className,
      )}
    >
      {STATUS_LABELS[status]}
    </span>
  );
}

export function getDemandStatusLabel(status: DemandStatus) {
  return STATUS_LABELS[status];
}

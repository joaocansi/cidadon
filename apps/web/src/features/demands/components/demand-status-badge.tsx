import type { DemandStatus } from "@/lib/api";
import { cn } from "@/lib/utils";

export const DEMAND_STATUS_META: Record<
  DemandStatus,
  { label: string; badgeClassName: string; markerColor: string }
> = {
  registered: {
    label: "Registrada",
    badgeClassName: "bg-paper-2 text-ink-soft ring-line",
    markerColor: "#7b8169",
  },
  under_review: {
    label: "Em análise",
    badgeClassName: "bg-lime-pale text-lime-deep ring-lime/30",
    markerColor: "#6f9524",
  },
  in_progress: {
    label: "Em execução",
    badgeClassName: "bg-amber/15 text-amber-deep ring-amber/30",
    markerColor: "#e3a23a",
  },
  awaiting_confirmation: {
    label: "Aguardando validação",
    badgeClassName: "bg-sky-100 text-sky-800 ring-sky-200",
    markerColor: "#3b82f6",
  },
  completed: {
    label: "Concluída",
    badgeClassName: "bg-lime text-pine ring-lime-deep/20",
    markerColor: "#1f3d2a",
  },
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
        DEMAND_STATUS_META[status].badgeClassName,
        className,
      )}
    >
      {DEMAND_STATUS_META[status].label}
    </span>
  );
}

export function getDemandStatusLabel(status: DemandStatus) {
  return DEMAND_STATUS_META[status].label;
}

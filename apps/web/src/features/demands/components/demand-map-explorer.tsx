"use client";

import {
  CalendarClockIcon,
  ChevronRightIcon,
  CircleXIcon,
  MapPinIcon,
  MessageSquareIcon,
  ThumbsUpIcon,
} from "lucide-react";
import dynamic from "next/dynamic";
import Link from "next/link";
import { useEffect, useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import { DemandStatusBadge } from "@/features/demands/components/demand-status-badge";
import type { Demand } from "@/lib/api";
import { cn } from "@/lib/utils";

const DemandMap = dynamic(
  () => import("@/features/demands/maps/demand-map").then((module) => module.DemandMap),
  { ssr: false, loading: () => <div className="bg-paper-2 h-full animate-pulse" /> },
);

const dateFormatter = new Intl.DateTimeFormat("pt-BR", {
  dateStyle: "medium",
  timeStyle: "short",
});

export function DemandMapExplorer({
  demands,
  selected,
  onSelect,
  className,
  mapClassName,
  emptyMessage = "Nenhuma demanda com localização foi encontrada.",
}: {
  demands: Demand[];
  selected?: Demand;
  onSelect: (demand: Demand) => void;
  className?: string;
  mapClassName?: string;
  emptyMessage?: string;
}) {
  const [previewOpen, setPreviewOpen] = useState(false);
  const previousSelectedID = useRef<number | undefined>(undefined);

  useEffect(() => {
    if (selected?.id && previousSelectedID.current !== selected.id) {
      setPreviewOpen(true);
      previousSelectedID.current = selected.id;
    }
  }, [selected?.id]);

  function selectDemand(demand: Demand) {
    onSelect(demand);
    setPreviewOpen(true);
  }

  const hasMappedDemand = demands.some(
    (demand) => Number.isFinite(demand.latitude) && Number.isFinite(demand.longitude),
  );

  return (
    <section className={cn("relative", className)} aria-label="Mapa de demandas">
      <DemandMap
        className={mapClassName ?? "h-[min(70dvh,680px)] min-h-[500px] w-full rounded-xl"}
        demands={demands}
        selectedId={selected?.id}
        onSelect={selectDemand}
        fitToDemands
        showLegend
      />

      {!hasMappedDemand ? (
        <div className="bg-card/96 absolute inset-x-5 top-1/2 z-20 mx-auto max-w-sm -translate-y-1/2 rounded-xl border p-5 text-center shadow-[var(--shadow-card)] backdrop-blur">
          <MapPinIcon className="text-lime-deep mx-auto size-7" />
          <p className="mt-3 text-sm font-semibold">Nenhuma localização disponível</p>
          <p className="text-ink-soft mt-1 text-sm">{emptyMessage}</p>
        </div>
      ) : null}

      {selected && previewOpen ? (
        <>
          <div className="absolute inset-x-2 bottom-2 z-20 sm:hidden">
            <DemandMapPreview demand={selected} onClose={() => setPreviewOpen(false)} mobile />
          </div>
          <div className="absolute top-3 right-3 bottom-3 z-20 hidden w-[min(23rem,calc(100%-1.5rem))] sm:block">
            <DemandMapPreview demand={selected} onClose={() => setPreviewOpen(false)} />
          </div>
        </>
      ) : null}
    </section>
  );
}

function DemandMapPreview({
  demand,
  onClose,
  mobile = false,
}: {
  demand: Demand;
  onClose: () => void;
  mobile?: boolean;
}) {
  const address = [demand.street, demand.number, demand.neighborhood].filter(Boolean).join(", ");
  return (
    <article
      className={cn(
        "border-line bg-card/98 relative flex h-full flex-col overflow-y-auto rounded-xl border p-4 shadow-[0_20px_50px_-24px_rgba(30,36,23,.58)] backdrop-blur-xl",
        mobile && "max-h-[min(62dvh,30rem)] rounded-t-2xl",
      )}
      aria-label={`Prévia da demanda ${demand.protocol}`}
    >
      {mobile ? (
        <span className="bg-line mx-auto mb-3 h-1 w-10 rounded-full" aria-hidden="true" />
      ) : null}
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <p className="text-ink-faint font-mono text-xs font-semibold">{demand.protocol}</p>
          <div className="mt-2">
            <DemandStatusBadge status={demand.status} />
          </div>
        </div>
        <Button
          type="button"
          variant="ghost"
          size="icon"
          onClick={onClose}
          aria-label="Fechar prévia da demanda"
          className="text-ink-soft shrink-0"
        >
          <CircleXIcon className="size-4" />
        </Button>
      </div>

      <div className="mt-4">
        <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">
          {demand.category}
        </p>
        <h3 className="font-display mt-1 text-xl leading-tight font-semibold">{demand.title}</h3>
        <p className="text-ink-soft mt-2 line-clamp-4 text-sm leading-6">{demand.description}</p>
      </div>

      <div className="border-line-soft text-ink-soft mt-4 space-y-2 border-y py-3 text-xs leading-5">
        <p className="flex gap-2">
          <MapPinIcon className="text-lime-deep mt-0.5 size-4 shrink-0" />
          <span>
            {address} · {demand.city}/{demand.state}
          </span>
        </p>
        <p className="flex gap-2">
          <CalendarClockIcon className="text-lime-deep mt-0.5 size-4 shrink-0" />
          Atualizada em {dateFormatter.format(new Date(demand.updated_at))}
        </p>
      </div>

      <div className="text-ink-soft mt-3 flex items-center gap-4 text-xs font-medium">
        <span className="inline-flex items-center gap-1.5">
          <ThumbsUpIcon className="text-lime-deep size-3.5" />
          {demand.support_count} apoios
        </span>
        <span className="inline-flex items-center gap-1.5">
          <MessageSquareIcon className="text-lime-deep size-3.5" />
          {demand.comment_count} comentários
        </span>
      </div>

      <Link
        href={`/demandas/${demand.id}`}
        target="_blank"
        rel="noopener noreferrer"
        className="bg-pine text-paper hover:bg-pine/90 mt-5 inline-flex h-10 items-center justify-center gap-2 rounded-lg px-4 text-sm font-semibold transition"
      >
        Mais detalhes
        <ChevronRightIcon className="size-4" />
      </Link>
    </article>
  );
}

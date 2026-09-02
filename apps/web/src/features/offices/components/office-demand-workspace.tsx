"use client";

import { DemandMapExplorer } from "@/features/demands/components/demand-map-explorer";
import type { Demand } from "@/lib/api";

export function OfficeDemandWorkspace({
  demands,
  selected,
  onSelect,
}: {
  demands: Demand[];
  selected?: Demand;
  onSelect: (demand: Demand) => void;
}) {
  return (
    <section id="demandas" className="scroll-mt-24 space-y-3">
      <div>
        <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">Território</p>
        <h2 className="font-display mt-1 text-xl font-semibold">Demandas atribuídas ao gabinete</h2>
        <p className="text-ink-soft mt-1 text-sm">
          Selecione um ponto para visualizar a prévia sem sair do mapa.
        </p>
      </div>
      <DemandMapExplorer
        demands={demands}
        selected={selected}
        onSelect={onSelect}
        mapClassName="h-[min(70dvh,680px)] min-h-[500px] w-full rounded-xl"
        emptyMessage="As demandas direcionadas a este gabinete aparecerão aqui quando forem registradas."
      />
    </section>
  );
}

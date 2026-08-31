"use client";

import { CalendarClockIcon, ImageIcon, MapPinIcon } from "lucide-react";
import dynamic from "next/dynamic";
import Link from "next/link";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DemandStatusBadge } from "@/features/demands/components/demand-status-badge";
import type { Demand } from "@/lib/api";

const DemandMap = dynamic(
  () => import("@/features/demands/maps/demand-map").then((module) => module.DemandMap),
  { ssr: false },
);

const dateFormatter = new Intl.DateTimeFormat("pt-BR", { dateStyle: "medium", timeStyle: "short" });

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
      </div>
      <div className="grid min-w-0 gap-4 xl:grid-cols-[minmax(0,1fr)_380px]">
        <Card className="border-line min-w-0 py-0">
          <CardContent className="p-0">
            <DemandMap
              className="h-[min(64vh,640px)] min-h-[480px] w-full"
              demands={demands}
              selectedId={selected?.id}
              onSelect={onSelect}
            />
          </CardContent>
        </Card>

        <Card className="border-line xl:max-h-[640px]">
          <CardHeader className="border-line-soft border-b">
            <CardTitle className="font-display text-lg">Detalhes da demanda</CardTitle>
          </CardHeader>
          <CardContent className="overflow-y-auto">
            {selected ? (
              <div className="space-y-5">
                <div className="flex flex-wrap items-center justify-between gap-2">
                  <span className="text-ink-faint font-mono text-xs font-semibold">
                    {selected.protocol}
                  </span>
                  <DemandStatusBadge status={selected.status} />
                </div>

                <div>
                  <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">
                    {selected.category}
                  </p>
                  <h3 className="font-display mt-1 text-xl font-semibold">{selected.title}</h3>
                  <p className="text-ink-soft mt-2 text-sm leading-6">{selected.description}</p>
                </div>

                <div className="bg-paper-2/60 text-ink-soft space-y-2 rounded-xl p-3 text-sm">
                  <p className="flex gap-2">
                    <MapPinIcon className="text-lime-deep mt-0.5 size-4 shrink-0" />
                    <span>
                      {[selected.street, selected.number, selected.neighborhood]
                        .filter(Boolean)
                        .join(", ")}{" "}
                      · {selected.city}/{selected.state}
                    </span>
                  </p>
                  <p className="flex items-center gap-2">
                    <CalendarClockIcon className="text-lime-deep size-4 shrink-0" />
                    Registrada em {dateFormatter.format(new Date(selected.created_at))}
                  </p>
                </div>

                {selected.images.length ? (
                  <div>
                    <p className="mb-2 flex items-center gap-2 text-sm font-semibold">
                      <ImageIcon className="text-lime-deep size-4" />
                      Imagens anexadas
                    </p>
                    <div className="grid grid-cols-3 gap-2">
                      {selected.images.map((image, index) => (
                        // eslint-disable-next-line @next/next/no-img-element
                        <img
                          key={`${selected.id}-${index}`}
                          src={image}
                          alt={`Anexo ${index + 1} da demanda`}
                          className="aspect-square rounded-lg object-cover"
                        />
                      ))}
                    </div>
                  </div>
                ) : null}

                <Link
                  href={`/demands/${selected.id}`}
                  className="bg-pine text-paper inline-flex h-10 items-center rounded-lg px-4 text-sm font-semibold"
                >
                  Abrir histórico e ações
                </Link>
              </div>
            ) : (
              <div className="grid min-h-64 place-items-center text-center">
                <div>
                  <MapPinIcon className="text-lime-deep mx-auto size-7" />
                  <p className="mt-3 font-semibold">Nenhuma demanda selecionada</p>
                  <p className="text-ink-soft mt-1 text-sm">
                    Clique em um ponto do mapa para ver todas as informações.
                  </p>
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </section>
  );
}

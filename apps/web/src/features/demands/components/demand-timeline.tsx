"use client";

import { CalendarClockIcon, FlagIcon } from "lucide-react";

import { UserAvatar } from "@/components/shared/user-avatar";
import { DemandImageGallery } from "@/features/demands/components/demand-image-gallery";
import type { DemandEvent } from "@/lib/api";

const eventLabel: Record<string, string> = {
  created: "Demanda registrada",
  claimed: "Gabinete assumiu o atendimento",
  execution_started: "Execução iniciada",
  confirmation_requested: "Validação solicitada ao cidadão",
  confirmed: "Conclusão confirmada pelo cidadão",
  reopened: "Demanda reaberta para nova análise",
  automatically_completed: "Concluída após o prazo de validação",
  milestone: "Marco publicado pelo gabinete",
};

const eventDescription: Record<string, string> = {
  created: "A solicitação foi registrada e entrou no fluxo de atendimento.",
  claimed: "O gabinete assumiu a responsabilidade por este atendimento.",
  execution_started: "O gabinete iniciou as ações para tratar a demanda.",
  confirmation_requested: "O gabinete informou uma entrega e aguarda a validação da pessoa autora.",
  confirmed: "A pessoa autora confirmou que a solução foi concluída.",
  reopened: "A pessoa autora sinalizou que ainda há pendências; a demanda voltou para análise.",
  automatically_completed:
    "O prazo de 120 horas para validação terminou sem manifestação da pessoa autora.",
  milestone: "Atualização relevante do andamento publicada pelo gabinete.",
};

export function DemandTimeline({ events }: { events: DemandEvent[] }) {
  if (!events.length) {
    return <p className="text-ink-soft py-6 text-sm">Ainda não há atualizações neste protocolo.</p>;
  }
  return (
    <ol className="before:bg-line relative space-y-1 before:absolute before:top-4 before:bottom-4 before:left-[7px] before:w-px">
      {events.map((event) => {
        const fromOffice =
          event.actor_role === "councillor" || event.actor_role === "office_member";
        const images = event.images ?? [];
        return (
          <li key={event.id} className="relative py-3 pl-7">
            <span className="border-card bg-lime absolute top-5 left-0 size-[15px] rounded-full border-4" />
            <div
              className={event.type === "milestone" ? "bg-lime-pale/55 rounded-xl p-3" : undefined}
            >
              <div className="flex items-start gap-2">
                {event.type === "milestone" ? (
                  <FlagIcon className="text-lime-deep mt-0.5 size-4 shrink-0" />
                ) : (
                  <CalendarClockIcon className="text-lime-deep mt-0.5 size-4 shrink-0" />
                )}
                <div className="min-w-0 flex-1">
                  <p className="text-sm font-semibold">
                    {eventLabel[event.type] ?? "Atualização da demanda"}
                  </p>
                  <p className="text-ink-soft mt-1 text-xs leading-5">
                    {eventDescription[event.type] ??
                      "Uma atualização foi registrada no atendimento."}
                  </p>
                  {event.actor_name ? (
                    <div className="mt-2 flex items-center gap-2">
                      <UserAvatar
                        name={event.actor_name}
                        imageUrl={event.actor_image_url}
                        className="size-6 text-[9px]"
                      />
                      <span className="text-ink-soft text-xs">
                        {event.actor_name}
                        {fromOffice ? " · Gabinete" : ""}
                      </span>
                    </div>
                  ) : null}
                  {event.message ? (
                    <p className="border-lime/35 bg-card/70 text-ink-soft mt-3 rounded-lg border-l-2 px-3 py-2 text-sm leading-6 whitespace-pre-wrap">
                      {event.message}
                    </p>
                  ) : null}
                  {images.length ? (
                    <DemandImageGallery
                      images={images}
                      altPrefix="Anexo da atualização"
                      className="mt-3 grid-cols-2"
                    />
                  ) : null}
                  <p className="text-ink-faint mt-2 text-xs">
                    {new Date(event.created_at).toLocaleString("pt-BR")}
                  </p>
                </div>
              </div>
            </div>
          </li>
        );
      })}
    </ol>
  );
}

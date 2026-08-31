import { CalendarClockIcon, MapPinIcon, MessageSquareIcon, ThumbsUpIcon } from "lucide-react";
import Link from "next/link";

import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import type { Demand } from "@/lib/api";

import { DemandStatusBadge } from "./demand-status-badge";

const dateFormatter = new Intl.DateTimeFormat("pt-BR", {
  day: "2-digit",
  month: "short",
  year: "numeric",
});

export function DemandCard({ demand }: { demand: Demand }) {
  const address = [demand.street, demand.number, demand.neighborhood].filter(Boolean).join(", ");

  return (
    <Card className="border-line bg-card rounded-[8px] border shadow-[var(--shadow-card)] transition hover:-translate-y-0.5 hover:shadow-md">
      <CardHeader className="gap-3">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <span className="text-ink-faint font-mono text-xs font-semibold tracking-[0.06em] uppercase">
            {demand.protocol}
          </span>
          <DemandStatusBadge status={demand.status} />
        </div>
        <div>
          <CardTitle className="font-display text-ink text-xl font-semibold">
            {demand.title}
          </CardTitle>
          <p className="text-ink-soft mt-2 line-clamp-3 text-[14.5px] leading-6">
            {demand.description}
          </p>
        </div>
      </CardHeader>
      <CardContent className="text-ink-soft flex flex-col gap-3 text-sm">
        <span className="inline-flex items-center gap-2">
          <MapPinIcon className="text-lime-deep size-4" />
          {address} · {demand.city}/{demand.state}
        </span>
        <span className="bg-lime-pale text-lime-deep inline-flex w-fit rounded-full px-3 py-1 text-xs font-semibold">
          {demand.category}
        </span>
      </CardContent>
      <CardFooter className="border-line-soft text-ink-faint flex flex-wrap justify-between gap-3 border-t pt-4 text-xs font-medium">
        <span className="inline-flex items-center gap-1.5">
          <CalendarClockIcon className="size-3.5" />
          {dateFormatter.format(new Date(demand.created_at))}
        </span>
        <span className="inline-flex items-center gap-4">
          <span className="inline-flex items-center gap-1.5">
            <ThumbsUpIcon className="text-lime-deep size-3.5" />
            {demand.support_count} apoios
          </span>
          <span className="inline-flex items-center gap-1.5">
            <MessageSquareIcon className="text-lime-deep size-3.5" />
            {demand.comment_count} comentários
          </span>
        </span>
        <Link
          href={`/demands/${demand.id}`}
          className="text-pine text-sm font-semibold hover:underline"
        >
          Ver detalhes
        </Link>
      </CardFooter>
    </Card>
  );
}

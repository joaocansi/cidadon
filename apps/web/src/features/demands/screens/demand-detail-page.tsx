"use client";

import {
  ArrowLeftIcon,
  CalendarClockIcon,
  ImageIcon,
  Loader2Icon,
  MapPinIcon,
  SendIcon,
} from "lucide-react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { type FormEvent, useCallback, useEffect, useMemo, useState } from "react";

import { CitizenShell } from "@/components/layout/citizen-shell";
import { DashboardShell } from "@/components/layout/dashboard-shell";
import { FormField } from "@/components/shared/forms/form-field";
import { showToast } from "@/components/shared/toast";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { RoleGate } from "@/features/auth/components/role-gate";
import { useSession } from "@/features/auth/components/session-provider";
import {
  DemandStatusBadge,
  getDemandStatusLabel,
} from "@/features/demands/components/demand-status-badge";
import {
  apiCommentDemand,
  apiDemandAction,
  apiDemandActivity,
  apiGetDemand,
  apiGetManagedOffice,
  apiListOfficeDemands,
  type Demand,
  type DemandActivity,
} from "@/lib/api";
import { apiErrorMessage } from "@/lib/forms";

const eventLabel: Record<string, string> = {
  created: "Demanda registrada",
  claimed: "Gabinete assumiu o atendimento",
  execution_started: "Execução iniciada",
  confirmation_requested: "Validação solicitada ao cidadão",
  confirmed: "Conclusão confirmada pelo cidadão",
  reopened: "Demanda reaberta para nova análise",
  comment_hidden: "Conteúdo moderado",
  automatically_completed: "Concluída após o prazo de validação",
  migrated: "Histórico migrado",
};
type Action = "claim" | "start" | "request-confirmation" | "confirm" | "reopen";

export default function DemandDetailPage() {
  const { id: rawID } = useParams<{ id: string }>();
  const id = Number(rawID);
  const { user } = useSession();
  const staff = user?.role === "councillor" || user?.role === "office_member";
  const [demand, setDemand] = useState<Demand>();
  const [activity, setActivity] = useState<DemandActivity>();
  const [officeID, setOfficeID] = useState<number>();
  const [isAssigned, setIsAssigned] = useState(false);
  const [body, setBody] = useState("");
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const load = useCallback(async () => {
    setLoading(true);
    const [demandResult, activityResult] = await Promise.all([
      apiGetDemand(id),
      apiDemandActivity(id),
    ]);
    if (demandResult.ok) setDemand(demandResult.data);
    else
      showToast(
        apiErrorMessage(demandResult.error, "Não foi possível carregar a demanda."),
        "error",
      );
    if (activityResult.ok) setActivity(activityResult.data);
    setLoading(false);
  }, [id]);
  useEffect(() => {
    if (!id) return;
    const timeout = window.setTimeout(() => void load(), 0);
    return () => window.clearTimeout(timeout);
  }, [id, load]);
  useEffect(() => {
    if (!staff) return;
    void Promise.all([apiGetManagedOffice(), apiListOfficeDemands()]).then(
      ([officeResult, demandsResult]) => {
        if (officeResult.ok) setOfficeID(officeResult.data?.office_id);
        if (demandsResult.ok)
          setIsAssigned((demandsResult.data ?? []).some((item) => item.id === id));
      },
    );
  }, [id, staff]);
  const canAct =
    staff &&
    isAssigned &&
    (!demand?.responsible_office_id || demand.responsible_office_id === officeID);
  const timeline = useMemo(
    () =>
      [
        ...(activity?.events ?? []).map((item) => ({
          kind: "event" as const,
          id: item.id,
          createdAt: item.created_at,
          item,
        })),
        ...(activity?.comments ?? []).map((item) => ({
          kind: "comment" as const,
          id: item.id,
          createdAt: item.created_at,
          item,
        })),
      ].sort((a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()),
    [activity],
  );
  async function runAction(action: Action) {
    const updateRequired = action === "request-confirmation" || action === "reopen";
    if (updateRequired && !body.trim()) {
      showToast("Escreva uma atualização antes de continuar.", "error");
      return;
    }
    setSaving(true);
    const result = await apiDemandAction(id, action, updateRequired ? { body } : undefined);
    setSaving(false);
    if (!result.ok) {
      showToast(apiErrorMessage(result.error, "Não foi possível atualizar a demanda."), "error");
      return;
    }
    setBody("");
    showToast("Demanda atualizada com sucesso.");
    await load();
  }
  async function comment(event: FormEvent) {
    event.preventDefault();
    if (!body.trim()) return;
    setSaving(true);
    const result = await apiCommentDemand(id, { body });
    setSaving(false);
    if (!result.ok) {
      showToast(apiErrorMessage(result.error, "Não foi possível publicar o comentário."), "error");
      return;
    }
    setBody("");
    showToast("Comentário publicado.");
    await load();
  }
  const content =
    loading || !demand ? (
      <div className="grid min-h-[56vh] place-items-center">
        <span className="text-ink-soft flex items-center gap-2 text-sm">
          <Loader2Icon className="size-4 animate-spin" />
          Carregando demanda…
        </span>
      </div>
    ) : (
      <div className="grid gap-5 xl:grid-cols-[minmax(0,1fr)_360px]">
        <div className="min-w-0 space-y-5">
          <Card className="border-line">
            <CardHeader className="border-line-soft border-b">
              <div className="flex flex-wrap items-start justify-between gap-4">
                <div>
                  <p className="text-ink-faint font-mono text-xs font-semibold">
                    {demand.protocol}
                  </p>
                  <CardTitle className="font-display mt-2 max-w-3xl text-2xl leading-tight sm:text-3xl">
                    {demand.title}
                  </CardTitle>
                </div>
                <DemandStatusBadge status={demand.status} />
              </div>
            </CardHeader>
            <CardContent className="space-y-5">
              <p className="text-ink-soft max-w-3xl text-[15px] leading-7">{demand.description}</p>
              <div className="border-line-soft text-ink-soft grid gap-3 border-y py-4 text-sm sm:grid-cols-2">
                <p className="flex gap-2">
                  <MapPinIcon className="text-lime-deep mt-0.5 size-4 shrink-0" />
                  {[demand.street, demand.number, demand.neighborhood]
                    .filter(Boolean)
                    .join(", ")} · {demand.city}/{demand.state}
                </p>
                <p className="flex gap-2">
                  <CalendarClockIcon className="text-lime-deep mt-0.5 size-4 shrink-0" />
                  Registrada em {new Date(demand.created_at).toLocaleDateString("pt-BR")}
                </p>
              </div>
              {demand.images.length ? (
                <div>
                  <p className="mb-3 flex items-center gap-2 text-sm font-semibold">
                    <ImageIcon className="text-lime-deep size-4" />
                    Anexos da solicitação
                  </p>
                  <div className="grid grid-cols-2 gap-3 sm:grid-cols-3">
                    {demand.images.map((image, index) => (
                      // Local data URLs are intentionally rendered without Next image optimization.
                      // eslint-disable-next-line @next/next/no-img-element
                      <img
                        key={image}
                        src={image}
                        alt={`Anexo ${index + 1}`}
                        className="border-line aspect-video rounded-xl border object-cover"
                      />
                    ))}
                  </div>
                </div>
              ) : null}
            </CardContent>
          </Card>
          <Card className="border-line">
            <CardHeader className="border-line-soft border-b">
              <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">
                Registro auditável
              </p>
              <CardTitle className="font-display text-xl">Linha do tempo</CardTitle>
            </CardHeader>
            <CardContent className="pt-1">
              {timeline.length ? (
                <ol className="before:bg-line relative space-y-1 before:absolute before:top-4 before:bottom-4 before:left-[7px] before:w-px">
                  {timeline.map((entry) => (
                    <li key={`${entry.kind}-${entry.id}`} className="relative py-3 pl-7">
                      <span className="border-card bg-lime absolute top-5 left-0 size-[15px] rounded-full border-4" />
                      {entry.kind === "event" ? (
                        <>
                          <p className="text-sm font-semibold">
                            {eventLabel[entry.item.type] ?? "Atualização da demanda"}
                          </p>
                          <p className="text-ink-faint mt-1 text-xs">
                            {new Date(entry.createdAt).toLocaleString("pt-BR")}
                          </p>
                        </>
                      ) : (
                        <>
                          <div className="bg-paper-2 rounded-xl px-4 py-3">
                            <p className="text-sm font-semibold">{entry.item.author_name}</p>
                            <p className="text-ink-soft mt-1 text-sm leading-6 whitespace-pre-wrap">
                              {entry.item.hidden ? "Conteúdo moderado." : entry.item.body}
                            </p>
                          </div>
                          <p className="text-ink-faint mt-1 text-xs">
                            {new Date(entry.createdAt).toLocaleString("pt-BR")}
                          </p>
                        </>
                      )}
                    </li>
                  ))}
                </ol>
              ) : (
                <p className="text-ink-soft py-6 text-sm">
                  Ainda não há atualizações neste protocolo.
                </p>
              )}
            </CardContent>
          </Card>
        </div>
        <aside className="space-y-5 xl:sticky xl:top-22 xl:self-start">
          <Card className="border-line">
            <CardHeader className="border-line-soft border-b">
              <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">
                Situação atual
              </p>
              <CardTitle>{getDemandStatusLabel(demand.status)}</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <ActionPanel
                demand={demand}
                canAct={canAct}
                isAuthor={user?.role === "citizen"}
                saving={saving}
                onAction={runAction}
              />
              <p className="bg-paper-2 text-ink-soft rounded-lg px-3 py-2 text-xs leading-5">
                As ações ficam disponíveis apenas para quem participa deste atendimento.
              </p>
            </CardContent>
          </Card>
          <Card className="border-line">
            <CardHeader className="border-line-soft border-b">
              <CardTitle>Adicionar atualização</CardTitle>
            </CardHeader>
            <CardContent>
              <form onSubmit={comment} className="space-y-3">
                <FormField
                  id="comment"
                  label="Comentário público"
                  hint="Pessoas autenticadas podem acompanhar."
                >
                  <textarea
                    id="comment"
                    value={body}
                    onChange={(event) => setBody(event.target.value)}
                    placeholder="Escreva uma atualização útil"
                    className="field-textarea min-h-32"
                  />
                </FormField>
                <Button type="submit" disabled={saving || !body.trim()} className="w-full">
                  <SendIcon />
                  Publicar comentário
                </Button>
              </form>
            </CardContent>
          </Card>
        </aside>
      </div>
    );
  return (
    <RoleGate allowed={["citizen", "councillor", "office_member"]}>
      {staff ? (
        <DashboardShell
          title="Demanda"
          subtitle="Histórico e ações do protocolo."
          officeId={officeID}
        >
          <div className="mb-5">
            <Link
              href="/office"
              className="text-ink-soft hover:text-ink inline-flex items-center gap-1.5 text-sm font-semibold"
            >
              <ArrowLeftIcon className="size-4" />
              Voltar às demandas
            </Link>
          </div>
          {content}
        </DashboardShell>
      ) : (
        <CitizenShell
          title="Acompanhar demanda"
          subtitle="Consulte o histórico e participe do atendimento."
        >
          {content}
        </CitizenShell>
      )}
    </RoleGate>
  );
}

function ActionPanel({
  demand,
  canAct,
  isAuthor,
  saving,
  onAction,
}: {
  demand: Demand;
  canAct: boolean;
  isAuthor: boolean;
  saving: boolean;
  onAction: (action: Action) => void;
}) {
  if (demand.status === "registered" && !demand.responsible_office_id)
    return canAct ? (
      <Button className="w-full" disabled={saving} onClick={() => onAction("claim")}>
        Assumir atendimento
      </Button>
    ) : (
      <p className="text-ink-soft text-sm">Aguardando um gabinete assumir o atendimento.</p>
    );
  if (demand.status === "under_review")
    return canAct ? (
      <Button className="w-full" disabled={saving} onClick={() => onAction("start")}>
        Iniciar execução
      </Button>
    ) : (
      <p className="text-ink-soft text-sm">O gabinete responsável está analisando esta demanda.</p>
    );
  if (demand.status === "in_progress")
    return canAct ? (
      <Button className="w-full" disabled={saving} onClick={() => onAction("request-confirmation")}>
        Solicitar validação
      </Button>
    ) : (
      <p className="text-ink-soft text-sm">
        O gabinete responsável está trabalhando nesta demanda.
      </p>
    );
  if (demand.status === "awaiting_confirmation")
    return isAuthor ? (
      <div className="grid gap-2">
        <Button disabled={saving} onClick={() => onAction("confirm")}>
          Confirmar conclusão
        </Button>
        <Button variant="outline" disabled={saving} onClick={() => onAction("reopen")}>
          Reabrir demanda
        </Button>
      </div>
    ) : (
      <p className="text-ink-soft text-sm">Aguardando a validação de quem abriu a demanda.</p>
    );
  return <p className="text-ink-soft text-sm">Esta demanda foi concluída.</p>;
}

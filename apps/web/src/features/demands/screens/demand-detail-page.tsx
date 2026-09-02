"use client";

import {
  ArrowLeftIcon,
  CalendarClockIcon,
  FlagIcon,
  HeartIcon,
  ImageIcon,
  Loader2Icon,
  MapPinIcon,
} from "lucide-react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

import { CitizenShell } from "@/components/layout/citizen-shell";
import { DashboardShell } from "@/components/layout/dashboard-shell";
import { showToast } from "@/components/shared/toast";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { RoleGate } from "@/features/auth/components/role-gate";
import { useSession } from "@/features/auth/components/session-provider";
import { DemandComments } from "@/features/demands/components/demand-comments";
import { DemandImageGallery } from "@/features/demands/components/demand-image-gallery";
import {
  DemandStatusBadge,
  getDemandStatusLabel,
} from "@/features/demands/components/demand-status-badge";
import { DemandTimeline } from "@/features/demands/components/demand-timeline";
import {
  DemandTransitionDialog,
  type TimelinePayload,
} from "@/features/demands/components/demand-transition-dialog";
import {
  apiAddDemandSupport,
  apiCreateDemandMilestone,
  apiDemandAction,
  apiDemandActivity,
  apiDemandSupport,
  apiGetDemand,
  apiGetManagedOffice,
  apiListOfficeDemands,
  apiRemoveDemandSupport,
  type Demand,
  type DemandActivity,
  type DemandSupport,
} from "@/lib/api";
import { apiErrorMessage } from "@/lib/forms";

type Action = "claim" | "start" | "request-confirmation" | "confirm" | "reopen";
type TimelineAction = Exclude<Action, "confirm"> | "milestone";

const actionCopy: Record<
  TimelineAction,
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
  reopen: {
    title: "Reabrir demanda",
    description:
      "Descreva o que ainda não foi resolvido para que o gabinete possa revisar o atendimento.",
    submitLabel: "Reabrir demanda",
  },
  milestone: {
    title: "Novo marco",
    description: "Publique uma atualização relevante sobre o andamento desta demanda.",
    submitLabel: "Publicar marco",
  },
};

export default function DemandDetailPage() {
  const { id: rawID } = useParams<{ id: string }>();
  const id = Number(rawID);
  const { user } = useSession();
  const userID = user?.id;
  const staff = user?.role === "councillor" || user?.role === "office_member";
  const [demand, setDemand] = useState<Demand>();
  const [activity, setActivity] = useState<DemandActivity>();
  const [officeID, setOfficeID] = useState<number>();
  const [officeSlug, setOfficeSlug] = useState<string>();
  const [officeImageUrl, setOfficeImageUrl] = useState<string>();
  const [isAssigned, setIsAssigned] = useState(false);
  const [support, setSupport] = useState<DemandSupport>();
  const [pendingAction, setPendingAction] = useState<TimelineAction>();
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
    if (!id || !userID) return;
    const timeout = window.setTimeout(() => void load(), 0);
    return () => window.clearTimeout(timeout);
  }, [id, load, userID]);
  useEffect(() => {
    if (!staff) return;
    void Promise.all([apiGetManagedOffice(), apiListOfficeDemands()]).then(
      ([officeResult, demandsResult]) => {
        if (officeResult.ok && officeResult.data) {
          setOfficeID(officeResult.data.office_id);
          setOfficeSlug(officeResult.data.slug);
          setOfficeImageUrl(officeResult.data.image_url);
        }
        if (demandsResult.ok)
          setIsAssigned((demandsResult.data ?? []).some((item) => item.id === id));
      },
    );
  }, [id, staff]);
  useEffect(() => {
    if (user?.role !== "citizen" || !id) return;
    void apiDemandSupport(id).then((result) => {
      if (result.ok) setSupport(result.data);
    });
  }, [id, user?.role]);
  const canAct =
    staff &&
    isAssigned &&
    (!demand?.responsible_office_id || demand.responsible_office_id === officeID);
  const isAuthor =
    user?.role === "citizen" &&
    activity?.events.find((event) => event.type === "created")?.actor_user_id === user.id;
  async function runAction(action: Exclude<TimelineAction, "milestone">, payload: TimelinePayload) {
    setSaving(true);
    const result = await apiDemandAction(id, action, {
      message: payload.message,
      images: payload.images,
    });
    setSaving(false);
    if (!result.ok) {
      showToast(apiErrorMessage(result.error, "Não foi possível atualizar a demanda."), "error");
      return false;
    }
    showToast("Demanda atualizada com sucesso.");
    await load();
    return true;
  }
  async function publishMilestone(payload: TimelinePayload) {
    setSaving(true);
    const result = await apiCreateDemandMilestone(id, payload);
    setSaving(false);
    if (!result.ok) {
      showToast(apiErrorMessage(result.error, "Não foi possível publicar o marco."), "error");
      return false;
    }
    showToast("Marco publicado com sucesso.");
    await load();
    return true;
  }
  async function confirmDemand() {
    setSaving(true);
    const result = await apiDemandAction(id, "confirm");
    setSaving(false);
    if (!result.ok) {
      showToast(apiErrorMessage(result.error, "Não foi possível confirmar a conclusão."), "error");
      return;
    }
    showToast("Conclusão confirmada. Obrigado pela validação.");
    await load();
  }
  async function toggleSupport() {
    if (!support?.can_support) return;
    setSaving(true);
    const result = support.supported
      ? await apiRemoveDemandSupport(id)
      : await apiAddDemandSupport(id);
    setSaving(false);
    if (!result.ok) {
      showToast(apiErrorMessage(result.error, "Não foi possível atualizar seu apoio."), "error");
      return;
    }
    setSupport(result.data);
    showToast(result.data?.supported ? "Você apoiou esta demanda." : "Seu apoio foi removido.");
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
      <div className="grid min-w-0 gap-5 xl:grid-cols-[minmax(0,1fr)_360px]">
        <div className="min-w-0 space-y-5">
          <Card className="border-line">
            <CardHeader className="border-line-soft border-b">
              <div className="flex flex-wrap items-start justify-between gap-4">
                <div className="min-w-0">
                  <p className="text-ink-faint font-mono text-xs font-semibold">
                    {demand.protocol}
                  </p>
                  <CardTitle className="font-display mt-2 max-w-3xl text-2xl leading-tight break-words sm:text-3xl">
                    {demand.title}
                  </CardTitle>
                </div>
                <DemandStatusBadge status={demand.status} />
              </div>
            </CardHeader>
            <CardContent className="space-y-5">
              <p className="text-ink-soft max-w-3xl text-[15px] leading-7 break-words">
                {demand.description}
              </p>
              <div className="border-line-soft text-ink-soft grid gap-3 border-y py-4 text-sm sm:grid-cols-2">
                <p className="flex min-w-0 gap-2">
                  <MapPinIcon className="text-lime-deep mt-0.5 size-4 shrink-0" />
                  <span className="min-w-0 break-words">
                    {[demand.street, demand.number, demand.neighborhood].filter(Boolean).join(", ")}{" "}
                    · {demand.city}/{demand.state}
                  </span>
                </p>
                <p className="flex min-w-0 gap-2">
                  <CalendarClockIcon className="text-lime-deep mt-0.5 size-4 shrink-0" />
                  <span className="break-words">
                    Registrada em {new Date(demand.created_at).toLocaleDateString("pt-BR")}
                  </span>
                </p>
              </div>
              <div className="text-ink-soft flex items-center gap-2 text-sm">
                <HeartIcon className="text-lime-deep size-4" />
                <span>
                  <strong className="text-ink">
                    {support?.support_count ?? demand.support_count}
                  </strong>{" "}
                  apoio{(support?.support_count ?? demand.support_count) === 1 ? "" : "s"}
                </span>
              </div>
              {(demand.images ?? []).length ? (
                <div>
                  <p className="mb-3 flex items-center gap-2 text-sm font-semibold">
                    <ImageIcon className="text-lime-deep size-4" />
                    Anexos da solicitação
                  </p>
                  <DemandImageGallery
                    images={demand.images ?? []}
                    altPrefix="Anexo da solicitação"
                  />
                </div>
              ) : null}
            </CardContent>
          </Card>
          <Card className="border-line">
            <CardHeader className="border-line-soft border-b">
              <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">
                Conversa pública
              </p>
              <CardTitle className="font-display text-xl">Comentários</CardTitle>
            </CardHeader>
            <CardContent className="pt-5">
              <DemandComments
                demandID={id}
                comments={activity?.comments ?? []}
                currentUserID={user?.id}
                userRole={user?.role}
                officeImageUrl={officeImageUrl}
                onPublished={load}
              />
            </CardContent>
          </Card>
        </div>
        <aside className="min-w-0 space-y-5 xl:sticky xl:top-22 xl:self-start">
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
                isAuthor={Boolean(isAuthor)}
                saving={saving}
                onAction={(action) => {
                  if (action === "confirm") void confirmDemand();
                  else setPendingAction(action);
                }}
              />
              {canAct ? (
                <Button
                  variant="outline"
                  className="w-full"
                  disabled={saving}
                  onClick={() => setPendingAction("milestone")}
                >
                  <FlagIcon />
                  Novo marco
                </Button>
              ) : null}
              {user?.role === "citizen" && support?.can_support ? (
                <Button
                  variant={support.supported ? "outline" : "secondary"}
                  className="w-full"
                  disabled={saving}
                  onClick={() => void toggleSupport()}
                >
                  <HeartIcon className={support.supported ? "fill-current" : undefined} />
                  {support.supported ? "Apoiado" : "Apoiar demanda"}
                </Button>
              ) : null}
              <p className="bg-paper-2 text-ink-soft rounded-lg px-3 py-2 text-xs leading-5">
                As ações ficam disponíveis apenas para quem participa deste atendimento.
              </p>
            </CardContent>
          </Card>
          <Card className="border-line max-h-[calc(100vh-22rem)] overflow-hidden">
            <CardHeader className="border-line-soft border-b">
              <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">
                Registro auditável
              </p>
              <CardTitle className="font-display text-xl">Linha do tempo</CardTitle>
            </CardHeader>
            <CardContent className="max-h-[calc(100vh-30rem)] scrollbar-thin overflow-y-auto pt-1">
              <DemandTimeline events={activity?.events ?? []} />
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
          officeSlug={officeSlug}
        >
          <div className="mb-5">
            <Link
              href="/gabinete"
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
      {pendingAction ? (
        <DemandTransitionDialog
          open
          title={actionCopy[pendingAction].title}
          description={actionCopy[pendingAction].description}
          submitLabel={actionCopy[pendingAction].submitLabel}
          onClose={() => setPendingAction(undefined)}
          onSubmit={(payload) => {
            if (pendingAction === "milestone") return publishMilestone(payload);
            return runAction(pendingAction, payload);
          }}
        />
      ) : null}
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

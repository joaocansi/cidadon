"use client";

import { CheckCircle2Icon, Clock3Icon, InboxIcon } from "lucide-react";
import { useEffect, useState } from "react";

import { DashboardShell } from "@/components/layout/dashboard-shell";
import { showToast } from "@/components/shared/toast";
import { Card, CardContent } from "@/components/ui/card";
import { RoleGate } from "@/features/auth/components/role-gate";
import { useSession } from "@/features/auth/components/session-provider";
import {
  DemandStatusBadge,
  getDemandStatusLabel,
} from "@/features/demands/components/demand-status-badge";
import { OfficeDemandWorkspace } from "@/features/offices/components/office-demand-workspace";
import {
  apiCreateOffice,
  apiGetManagedOffice,
  apiListOfficeDemands,
  type Demand,
  type ManagedOffice,
} from "@/lib/api";
import { apiErrorMessage } from "@/lib/forms";

const dateFormatter = new Intl.DateTimeFormat("pt-BR", { dateStyle: "medium", timeStyle: "short" });

export default function OfficePage() {
  const { user } = useSession();
  const [office, setOffice] = useState<ManagedOffice | null>(null);
  const [demands, setDemands] = useState<Demand[]>([]);
  const [selected, setSelected] = useState<Demand>();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!user) return;
    let active = true;
    void apiGetManagedOffice().then(async (initialResult) => {
      let officeResult = initialResult;
      if (
        !officeResult.ok &&
        officeResult.error?.code === "NOT_FOUND" &&
        user.role === "councillor"
      ) {
        const created = await apiCreateOffice();
        if (created.ok) officeResult = await apiGetManagedOffice();
      }
      if (!active) return;
      if (!officeResult.ok || !officeResult.data) {
        setLoading(false);
        showToast(
          apiErrorMessage(officeResult.error, "Não foi possível carregar o gabinete."),
          "error",
        );
        return;
      }
      setOffice(officeResult.data);
      const demandResult = await apiListOfficeDemands();
      if (!active) return;
      setLoading(false);
      if (!demandResult.ok) {
        showToast(
          apiErrorMessage(demandResult.error, "Não foi possível carregar as demandas."),
          "error",
        );
        return;
      }
      const items = demandResult.data ?? [];
      setDemands(items);
      setSelected(items[0]);
    });
    return () => {
      active = false;
    };
  }, [user]);

  const recent = [...demands]
    .sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())
    .slice(0, 5);
  return (
    <RoleGate allowed={["councillor", "office_member"]}>
      <DashboardShell
        title={office ? `Gabinete ${office.councillor_name}` : "Painel do gabinete"}
        subtitle="Acompanhe a região e avance as demandas com clareza."
        officeId={office?.office_id}
      >
        {loading ? (
          <div className="text-ink-soft grid min-h-[60vh] place-items-center text-sm">
            Preparando o painel do gabinete…
          </div>
        ) : (
          <div className="space-y-8">
            <section className="grid gap-3 sm:grid-cols-3">
              <Metric label="Demandas recebidas" value={demands.length} icon={InboxIcon} />
              <Metric
                label="Em andamento"
                value={demands.filter((item) => item.status === "in_progress").length}
                icon={Clock3Icon}
              />
              <Metric
                label="Concluídas"
                value={demands.filter((item) => item.status === "completed").length}
                icon={CheckCircle2Icon}
              />
            </section>
            <OfficeDemandWorkspace demands={demands} selected={selected} onSelect={setSelected} />
            <section className="space-y-3">
              <div>
                <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">
                  Acompanhamento
                </p>
                <h2 className="font-display mt-1 text-xl font-semibold">Últimas atualizações</h2>
              </div>
              <Card className="border-line">
                <CardContent className="p-0">
                  {recent.length ? (
                    <div className="divide-line-soft divide-y">
                      {recent.map((demand) => (
                        <button
                          key={demand.id}
                          type="button"
                          onClick={() => setSelected(demand)}
                          className="hover:bg-lime-pale/30 flex w-full items-center gap-3 p-4 text-left transition sm:px-5"
                        >
                          <DemandStatusBadge status={demand.status} />
                          <div className="min-w-0 flex-1">
                            <p className="truncate text-sm font-semibold">{demand.title}</p>
                            <p className="text-ink-soft mt-0.5 text-xs">
                              {getDemandStatusLabel(demand.status)} · atualizado em{" "}
                              {dateFormatter.format(new Date(demand.updated_at))}
                            </p>
                          </div>
                          <span className="text-ink-faint font-mono text-xs">
                            {demand.protocol}
                          </span>
                        </button>
                      ))}
                    </div>
                  ) : (
                    <div className="text-ink-soft p-8 text-center text-sm">
                      As próximas atualizações das demandas aparecerão aqui.
                    </div>
                  )}
                </CardContent>
              </Card>
            </section>
          </div>
        )}
      </DashboardShell>
    </RoleGate>
  );
}

function Metric({
  label,
  value,
  icon: Icon,
}: {
  label: string;
  value: number;
  icon: typeof InboxIcon;
}) {
  return (
    <Card className="border-line">
      <CardContent className="flex items-center justify-between p-4">
        <div>
          <p className="text-ink-soft text-sm font-semibold">{label}</p>
          <p className="font-display mt-1 text-3xl font-semibold">{value}</p>
        </div>
        <span className="bg-lime-pale text-lime-deep grid size-10 place-items-center rounded-xl">
          <Icon className="size-5" />
        </span>
      </CardContent>
    </Card>
  );
}

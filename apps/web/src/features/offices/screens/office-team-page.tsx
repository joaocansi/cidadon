"use client";

import { type FormEvent, useEffect, useState } from "react";

import { DashboardShell } from "@/components/layout/dashboard-shell";
import { showToast } from "@/components/shared/toast";
import { RoleGate } from "@/features/auth/components/role-gate";
import { OfficeTeamPanel } from "@/features/offices/components/office-team-panel";
import {
  apiCancelOfficeMemberRequest,
  apiCreateOffice,
  apiGetManagedOffice,
  apiInviteOfficeMember,
  apiListOfficeMemberRequests,
  apiRemoveOfficeMember,
  type ManagedOffice,
  type OfficeMemberRequest,
} from "@/lib/api";
import { apiErrorMessage } from "@/lib/forms";

export default function OfficeTeamPage() {
  const [office, setOffice] = useState<ManagedOffice | null>(null);
  const [requests, setRequests] = useState<OfficeMemberRequest[]>([]);
  const [email, setEmail] = useState("");
  const [emailError, setEmailError] = useState<string>();
  const [loading, setLoading] = useState(true);
  const [inviting, setInviting] = useState(false);
  const [removingId, setRemovingId] = useState<number>();
  const [cancellingId, setCancellingId] = useState<number>();

  async function load() {
    setLoading(true);
    let officeResult = await apiGetManagedOffice();
    if (!officeResult.ok && officeResult.error?.code === "NOT_FOUND") {
      const created = await apiCreateOffice();
      if (created.ok) officeResult = await apiGetManagedOffice();
    }
    const requestsResult = await apiListOfficeMemberRequests();
    setLoading(false);
    if (!officeResult.ok || !officeResult.data) {
      showToast(
        apiErrorMessage(officeResult.error, "Não foi possível carregar a equipe."),
        "error",
      );
      return;
    }
    if (!requestsResult.ok) {
      showToast(
        apiErrorMessage(requestsResult.error, "Não foi possível carregar os convites."),
        "error",
      );
      return;
    }
    setOffice(officeResult.data);
    setRequests(requestsResult.data ?? []);
  }

  useEffect(() => {
    let active = true;
    void (async () => {
      let officeResult = await apiGetManagedOffice();
      if (!officeResult.ok && officeResult.error?.code === "NOT_FOUND") {
        const created = await apiCreateOffice();
        if (created.ok) officeResult = await apiGetManagedOffice();
      }
      const requestsResult = await apiListOfficeMemberRequests();
      if (!active) return;
      setLoading(false);
      if (!officeResult.ok || !officeResult.data) {
        showToast(
          apiErrorMessage(officeResult.error, "Não foi possível carregar a equipe."),
          "error",
        );
        return;
      }
      if (!requestsResult.ok) {
        showToast(
          apiErrorMessage(requestsResult.error, "Não foi possível carregar os convites."),
          "error",
        );
        return;
      }
      setOffice(officeResult.data);
      setRequests(requestsResult.data ?? []);
    })();
    return () => {
      active = false;
    };
  }, []);

  async function invite(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!/^\S+@\S+\.\S+$/.test(email)) {
      setEmailError("Informe um e-mail válido.");
      return;
    }
    setEmailError(undefined);
    setInviting(true);
    const result = await apiInviteOfficeMember(email);
    setInviting(false);
    if (!result.ok) {
      showToast(apiErrorMessage(result.error, "Não foi possível enviar o convite."), "error");
      return;
    }
    setEmail("");
    await load();
    showToast("Convite enviado com sucesso.");
  }

  async function removeMember(member: ManagedOffice["members"][number]) {
    if (
      !window.confirm(
        `Excluir ${member.name} da equipe? O cadastro de membro e o acesso ao gabinete serão removidos.`,
      )
    )
      return;
    setRemovingId(member.user_id);
    const result = await apiRemoveOfficeMember(member.user_id);
    setRemovingId(undefined);
    if (!result.ok) {
      showToast(apiErrorMessage(result.error, "Não foi possível remover o membro."), "error");
      return;
    }
    setOffice((current) =>
      current
        ? { ...current, members: current.members.filter((item) => item.user_id !== member.user_id) }
        : current,
    );
    showToast("Membro removido e acesso revogado.");
  }

  async function cancelRequest(request: OfficeMemberRequest) {
    if (!window.confirm(`Cancelar o convite enviado para ${request.email}?`)) return;
    setCancellingId(request.id);
    const result = await apiCancelOfficeMemberRequest(request.id);
    setCancellingId(undefined);
    if (!result.ok) {
      showToast(apiErrorMessage(result.error, "Não foi possível cancelar o convite."), "error");
      return;
    }
    setRequests((items) => items.filter((item) => item.id !== request.id));
    showToast("Convite cancelado.");
  }

  return (
    <RoleGate allowed={["councillor"]}>
      <DashboardShell
        title="Equipe do gabinete"
        subtitle="Gerencie pessoas que podem atuar nas demandas e os convites em aberto."
        officeSlug={office?.slug}
      >
        {loading ? (
          <div className="text-ink-soft grid min-h-[60vh] place-items-center text-sm">
            Carregando equipe…
          </div>
        ) : (
          <OfficeTeamPanel
            members={office?.members ?? []}
            requests={requests}
            email={email}
            error={emailError}
            inviting={inviting}
            removingId={removingId}
            cancellingId={cancellingId}
            onEmailChange={(value) => {
              setEmail(value);
              setEmailError(undefined);
            }}
            onInvite={invite}
            onRemove={(member) => void removeMember(member)}
            onCancel={(request) => void cancelRequest(request)}
          />
        )}
      </DashboardShell>
    </RoleGate>
  );
}

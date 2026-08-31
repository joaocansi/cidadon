"use client";

import { type FormEvent, useEffect, useState } from "react";

import { DashboardShell } from "@/components/layout/dashboard-shell";
import { showToast } from "@/components/shared/toast";
import { RoleGate } from "@/features/auth/components/role-gate";
import { OfficeProfileEditor } from "@/features/offices/components/office-profile-editor";
import {
  apiCreateOffice,
  apiGetManagedOffice,
  apiUpdateOffice,
  type ManagedOffice,
  type OfficeContact,
} from "@/lib/api";
import { apiErrorMessage } from "@/lib/forms";

type ProfileErrors = { city?: string; state?: string };

export default function OfficeProfileSettingsPage() {
  const [office, setOffice] = useState<ManagedOffice | null>(null);
  const [description, setDescription] = useState("");
  const [city, setCity] = useState("");
  const [state, setState] = useState("");
  const [contacts, setContacts] = useState<OfficeContact[]>([]);
  const [socials, setSocials] = useState<OfficeContact[]>([]);
  const [errors, setErrors] = useState<ProfileErrors>({});
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  function applyOffice(data: ManagedOffice) {
    setOffice(data);
    setDescription(data.description ?? "");
    setCity(data.city ?? "");
    setState(data.state ?? "");
    setContacts(data.contacts ?? []);
    setSocials(data.social_networks ?? []);
  }

  useEffect(() => {
    let active = true;
    void apiGetManagedOffice().then(async (initialResult) => {
      let result = initialResult;
      if (!result.ok && result.error?.code === "NOT_FOUND") {
        const created = await apiCreateOffice();
        if (created.ok) result = await apiGetManagedOffice();
      }
      if (!active) return;
      setLoading(false);
      if (!result.ok || !result.data) {
        showToast(
          apiErrorMessage(result.error, "Não foi possível carregar as configurações do gabinete."),
          "error",
        );
        return;
      }
      applyOffice(result.data);
    });
    return () => {
      active = false;
    };
  }, []);

  async function save(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const nextErrors: ProfileErrors = {};
    if (city.trim().length < 2) nextErrors.city = "Informe a cidade onde o gabinete atua.";
    if (!/^[A-Za-z]{2}$/.test(state.trim())) nextErrors.state = "Informe a UF com duas letras.";
    setErrors(nextErrors);
    if (Object.keys(nextErrors).length) return;

    const invalidChannel = [...contacts, ...socials].some(
      (channel) => !channel.type.trim() || !channel.value.trim(),
    );
    if (invalidChannel) {
      showToast("Preencha ou remova os canais de contato e redes sociais em branco.", "error");
      return;
    }

    setSaving(true);
    const result = await apiUpdateOffice({
      description: description.trim(),
      city: city.trim(),
      state: state.trim().toUpperCase(),
      contacts,
      social_networks: socials,
    });
    setSaving(false);
    if (!result.ok || !result.data) {
      showToast(
        apiErrorMessage(result.error, "Não foi possível salvar as configurações do gabinete."),
        "error",
      );
      return;
    }
    setDescription(result.data.description);
    setCity(result.data.city);
    setState(result.data.state);
    setContacts(result.data.contacts);
    setSocials(result.data.social_networks);
    setOffice((current) => (current ? { ...current, ...result.data } : current));
    showToast("Perfil público do gabinete salvo com sucesso.");
  }

  return (
    <RoleGate allowed={["councillor"]}>
      <DashboardShell
        title="Configurações do gabinete"
        subtitle="Atualize as informações que cidadãos veem ao procurar sua representação."
        officeId={office?.office_id}
      >
        {loading ? (
          <div className="text-ink-soft grid min-h-[60vh] place-items-center text-sm">
            Carregando configurações…
          </div>
        ) : (
          <OfficeProfileEditor
            description={description}
            city={city}
            state={state}
            contacts={contacts}
            socials={socials}
            saving={saving}
            errors={errors}
            onDescriptionChange={setDescription}
            onCityChange={(value) => {
              setCity(value);
              setErrors((current) => ({ ...current, city: undefined }));
            }}
            onStateChange={(value) => {
              setState(value);
              setErrors((current) => ({ ...current, state: undefined }));
            }}
            onContactsChange={setContacts}
            onSocialsChange={setSocials}
            onSave={save}
          />
        )}
      </DashboardShell>
    </RoleGate>
  );
}

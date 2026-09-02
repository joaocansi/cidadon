"use client";

import { Loader2Icon } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { type FormEvent, useState } from "react";

import { AvatarCropInput } from "@/components/shared/avatar-crop-input";
import { FormField } from "@/components/shared/forms/form-field";
import { PartySelect } from "@/components/shared/party-select";
import { showToast } from "@/components/shared/toast";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { AuthShell } from "@/features/auth/components/auth-shell";
import { AuthenticatedRedirect } from "@/features/auth/components/authenticated-redirect";
import { apiRegisterCouncillor } from "@/lib/api";
import {
  apiErrorMessage,
  apiFieldErrors,
  type FieldErrors,
  invalidEmail,
  requiredFields,
} from "@/lib/forms";

const empty = {
  name: "",
  email: "",
  party: "",
  photo: null as File | null,
  city: "",
  state: "",
  password: "",
  confirmPassword: "",
};

export default function RegisterCouncillorPage() {
  const router = useRouter();
  const [form, setForm] = useState(empty);
  const [errors, setErrors] = useState<FieldErrors>({});
  const [loading, setLoading] = useState(false);

  function update(field: keyof typeof form, value: string | File | null) {
    setForm((current) => ({
      ...current,
      [field]: field === "state" && typeof value === "string" ? value.toUpperCase() : value,
    }));
    setErrors((current) => ({ ...current, [field]: "" }));
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const next = requiredFields(
      {
        name: form.name,
        email: form.email,
        party: form.party,
        photo: form.photo ? "foto" : "",
        city: form.city,
        state: form.state,
        password: form.password,
        confirmPassword: form.confirmPassword,
      },
      {
        name: "Nome completo",
        email: "E-mail",
        party: "Partido",
        photo: "Foto",
        city: "Cidade",
        state: "UF",
        password: "Senha",
        confirmPassword: "Confirmação de senha",
      },
    );
    if (form.name && form.name.trim().length < 3) next.name = "Informe seu nome completo.";
    if (form.email && invalidEmail(form.email)) next.email = "Informe um e-mail válido.";
    if (form.password && (form.password.length < 6 || form.password.length > 72))
      next.password = "A senha deve ter entre 6 e 72 caracteres.";
    if (form.password !== form.confirmPassword) next.confirmPassword = "As senhas não coincidem.";
    if (form.state && form.state.length !== 2) next.state = "Informe a UF com 2 letras.";
    setErrors(next);
    if (Object.keys(next).length) return;

    setLoading(true);
    const result = await apiRegisterCouncillor({
      name: form.name,
      email: form.email,
      party: form.party,
      photo: form.photo!,
      city: form.city,
      state: form.state,
      password: form.password,
    });
    setLoading(false);
    if (!result.ok) {
      setErrors(apiFieldErrors(result.error));
      showToast(apiErrorMessage(result.error, "Não foi possível criar a conta."), "error");
      return;
    }
    showToast("Conta de vereador criada. Entre para configurar seu gabinete.");
    router.replace("/entrar");
  }

  return (
    <AuthenticatedRedirect>
      <AuthShell
        wide
        eyebrow="Cadastro de vereador"
        title={
          <>
            Crie seu <em className="text-lime-deep not-italic">gabinete</em> digital
          </>
        }
        subtitle="Organize sua equipe e acompanhe as demandas da população."
      >
        <form onSubmit={submit} className="space-y-5">
          <div className="grid gap-4 sm:grid-cols-2">
            <FormField id="name" label="Nome completo" error={errors.name}>
              <Input
                id="name"
                autoComplete="name"
                placeholder="Ana Martins"
                className="bg-card h-11"
                value={form.name}
                aria-invalid={Boolean(errors.name)}
                onChange={(event) => update("name", event.target.value)}
              />
            </FormField>
            <FormField id="email" label="E-mail" error={errors.email}>
              <Input
                id="email"
                type="email"
                autoComplete="email"
                placeholder="ana@camara.gov.br"
                className="bg-card h-11"
                value={form.email}
                aria-invalid={Boolean(errors.email)}
                onChange={(event) => update("email", event.target.value)}
              />
            </FormField>
            <FormField id="party" label="Partido" error={errors.party} className="sm:col-span-2">
              <PartySelect
                id="party"
                value={form.party}
                invalid={Boolean(errors.party)}
                onChange={(value) => update("party", value)}
              />
            </FormField>
          </div>
          <FormField id="photo" label="Foto de perfil" error={errors.photo}>
            <AvatarCropInput
              id="photo"
              value={form.photo}
              invalid={Boolean(errors.photo)}
              onChange={(value) => update("photo", value)}
              onInvalid={(message) => setErrors((current) => ({ ...current, photo: message }))}
            />
          </FormField>
          <div className="grid gap-4 sm:grid-cols-[minmax(0,1fr)_104px]">
            <FormField id="city" label="Cidade" error={errors.city}>
              <Input
                id="city"
                autoComplete="address-level2"
                placeholder="São Paulo"
                className="bg-card h-11"
                value={form.city}
                aria-invalid={Boolean(errors.city)}
                onChange={(event) => update("city", event.target.value)}
              />
            </FormField>
            <FormField id="state" label="UF" error={errors.state}>
              <Input
                id="state"
                autoComplete="address-level1"
                placeholder="SP"
                maxLength={2}
                className="bg-card h-11"
                value={form.state}
                aria-invalid={Boolean(errors.state)}
                onChange={(event) => update("state", event.target.value)}
              />
            </FormField>
          </div>
          <div className="grid gap-4 sm:grid-cols-2">
            <FormField id="password" label="Senha" error={errors.password}>
              <Input
                id="password"
                type="password"
                autoComplete="new-password"
                placeholder="Entre 6 e 72 caracteres"
                className="bg-card h-11"
                value={form.password}
                aria-invalid={Boolean(errors.password)}
                onChange={(event) => update("password", event.target.value)}
              />
            </FormField>
            <FormField id="confirmPassword" label="Confirmar senha" error={errors.confirmPassword}>
              <Input
                id="confirmPassword"
                type="password"
                autoComplete="new-password"
                placeholder="Repita a senha"
                className="bg-card h-11"
                value={form.confirmPassword}
                aria-invalid={Boolean(errors.confirmPassword)}
                onChange={(event) => update("confirmPassword", event.target.value)}
              />
            </FormField>
          </div>
          <Button type="submit" className="h-12 w-full" disabled={loading}>
            {loading ? <Loader2Icon className="animate-spin" /> : null}
            {loading ? "Criando conta…" : "Criar conta de vereador"}
          </Button>
          <p className="text-ink-soft text-center text-sm">
            Já tem uma conta?{" "}
            <Link href="/entrar" className="text-lime-deep font-semibold">
              Entrar
            </Link>
          </p>
        </form>
      </AuthShell>
    </AuthenticatedRedirect>
  );
}

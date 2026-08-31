"use client";

import { Building2Icon, CheckCircle2Icon, Loader2Icon, MapPinIcon } from "lucide-react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { type FormEvent, Suspense, useEffect, useState } from "react";

import { FormField } from "@/components/shared/forms/form-field";
import { showToast } from "@/components/shared/toast";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { AuthShell } from "@/features/auth/components/auth-shell";
import { AuthenticatedRedirect } from "@/features/auth/components/authenticated-redirect";
import {
  apiPreviewOfficeMemberInvitation,
  apiRegisterOfficeMember,
  type OfficeMemberInvitation,
} from "@/lib/api";
import { apiErrorMessage, apiFieldErrors, type FieldErrors } from "@/lib/forms";

export default function RegisterMemberPage() {
  return (
    <Suspense
      fallback={
        <div className="bg-paper text-ink-soft grid min-h-dvh place-items-center text-sm">
          Carregando convite…
        </div>
      }
    >
      <MemberRegistration />
    </Suspense>
  );
}

function MemberRegistration() {
  const token = useSearchParams().get("token") ?? "";
  const [form, setForm] = useState({ name: "", image_url: "", password: "", confirmPassword: "" });
  const [errors, setErrors] = useState<FieldErrors>({});
  const [invitation, setInvitation] = useState<OfficeMemberInvitation | null>(null);
  const [invitationIssue, setInvitationIssue] = useState<"invalid" | "unavailable" | null>(null);
  const [loadingInvitation, setLoadingInvitation] = useState(Boolean(token));
  const [done, setDone] = useState(false);
  const [loading, setLoading] = useState(false);

  async function loadInvitation() {
    if (!token) return;
    setLoadingInvitation(true);
    setInvitationIssue(null);
    const result = await apiPreviewOfficeMemberInvitation(token);
    setLoadingInvitation(false);
    if (!result.ok || !result.data) {
      setInvitation(null);
      setInvitationIssue(
        result.error?.code === "NETWORK" ||
          result.error?.code === "UNAVAILABLE" ||
          result.error?.code === "INTERNAL"
          ? "unavailable"
          : "invalid",
      );
      return;
    }
    setInvitation(result.data);
  }

  useEffect(() => {
    if (!token) return;
    let active = true;
    void apiPreviewOfficeMemberInvitation(token).then((result) => {
      if (!active) return;
      setLoadingInvitation(false);
      if (!result.ok || !result.data) {
        setInvitation(null);
        setInvitationIssue(
          result.error?.code === "NETWORK" ||
            result.error?.code === "UNAVAILABLE" ||
            result.error?.code === "INTERNAL"
            ? "unavailable"
            : "invalid",
        );
        return;
      }
      setInvitation(result.data);
    });
    return () => {
      active = false;
    };
  }, [token]);

  function update(field: keyof typeof form, value: string) {
    setForm((current) => ({ ...current, [field]: value }));
    setErrors((current) => ({ ...current, [field]: "" }));
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const next: FieldErrors = {};
    if (!form.name.trim() || form.name.trim().length < 3) next.name = "Informe seu nome completo.";
    if (!form.image_url.trim()) next.image_url = "Informe a URL da sua foto.";
    try {
      const photo = new URL(form.image_url);
      if (!(["http:", "https:"] as string[]).includes(photo.protocol))
        next.image_url = "Use uma URL de foto http ou https.";
    } catch {
      if (form.image_url) next.image_url = "Informe uma URL de foto válida.";
    }
    if (form.password.length < 6 || form.password.length > 72)
      next.password = "A senha deve ter entre 6 e 72 caracteres.";
    if (form.password !== form.confirmPassword) next.confirmPassword = "As senhas não coincidem.";
    setErrors(next);
    if (!token || Object.keys(next).length) return;

    setLoading(true);
    const result = await apiRegisterOfficeMember({
      token,
      name: form.name,
      password: form.password,
      image_url: form.image_url,
    });
    setLoading(false);
    if (!result.ok) {
      setErrors(apiFieldErrors(result.error));
      showToast(apiErrorMessage(result.error, "Não foi possível concluir seu cadastro."), "error");
      return;
    }
    showToast("Cadastro concluído. Você já pode acessar o gabinete.");
    setDone(true);
  }

  return (
    <AuthenticatedRedirect>
      <AuthShell
        wide
        eyebrow="Convite para equipe"
        title={
          invitation ? (
            <>
              Atue com o <em className="text-lime-deep not-italic">gabinete</em>
            </>
          ) : (
            <>
              Finalize seu <em className="text-lime-deep not-italic">cadastro</em>
            </>
          )
        }
        subtitle={
          invitation
            ? `Você foi convidado(a) para ajudar o gabinete de ${invitation.councillor_name} a atuar nas demandas da região.`
            : "Informe seus dados para começar a atuar nas demandas da região."
        }
      >
        {!token ? (
          <Alert variant="destructive">
            <AlertDescription>
              Este convite está incompleto, expirado ou não é válido. Solicite um novo link ao
              vereador.
            </AlertDescription>
          </Alert>
        ) : loadingInvitation ? (
          <div className="border-line bg-card text-ink-soft grid min-h-48 place-items-center rounded-2xl border text-sm">
            <span className="inline-flex items-center gap-2">
              <Loader2Icon className="size-4 animate-spin" />
              Validando convite…
            </span>
          </div>
        ) : invitationIssue ? (
          <Alert variant="destructive">
            <AlertDescription>
              {invitationIssue === "invalid"
                ? "Este convite expirou, já foi utilizado ou não é válido. Solicite um novo link ao gabinete."
                : "Não foi possível validar este convite agora. Tente novamente em instantes."}
              {invitationIssue === "unavailable" ? (
                <Button
                  type="button"
                  variant="link"
                  className="text-destructive ml-1 h-auto px-0"
                  onClick={() => void loadInvitation()}
                >
                  Tentar novamente
                </Button>
              ) : null}
            </AlertDescription>
          </Alert>
        ) : done ? (
          <div className="border-line bg-card space-y-4 rounded-2xl border p-7 text-center">
            <CheckCircle2Icon className="text-lime-deep mx-auto size-10" />
            <div>
              <p className="font-semibold">Cadastro concluído</p>
              <p className="text-ink-soft mt-1 text-sm">Sua conta já está vinculada ao gabinete.</p>
            </div>
            <Link
              className="bg-lime text-pine inline-flex h-10 items-center justify-center rounded-xl px-5 text-sm font-semibold"
              href="/login"
            >
              Entrar na plataforma
            </Link>
          </div>
        ) : (
          <div className="space-y-5">
            <InvitationSummary invitation={invitation!} />
            <form className="grid gap-4 sm:grid-cols-2" onSubmit={submit}>
              <FormField
                id="name"
                label="Nome completo"
                error={errors.name}
                className="sm:col-span-2"
              >
                <Input
                  id="name"
                  placeholder="Mariana Souza"
                  autoComplete="name"
                  className="bg-card h-11"
                  value={form.name}
                  aria-invalid={Boolean(errors.name)}
                  onChange={(event) => update("name", event.target.value)}
                />
              </FormField>
              <FormField
                id="photo"
                label="URL da foto"
                error={errors.image_url}
                hint="Use uma foto de perfil com link público."
                className="sm:col-span-2"
              >
                <Input
                  id="photo"
                  type="url"
                  className="bg-card h-11"
                  value={form.image_url}
                  aria-invalid={Boolean(errors.image_url)}
                  onChange={(event) => update("image_url", event.target.value)}
                  placeholder="https://exemplo.com/foto.jpg"
                />
              </FormField>
              <FormField id="password" label="Crie sua senha" error={errors.password}>
                <Input
                  id="password"
                  placeholder="Entre 6 e 72 caracteres"
                  className="bg-card h-11"
                  type="password"
                  autoComplete="new-password"
                  value={form.password}
                  aria-invalid={Boolean(errors.password)}
                  onChange={(event) => update("password", event.target.value)}
                />
              </FormField>
              <FormField
                id="confirmPassword"
                label="Confirmar senha"
                error={errors.confirmPassword}
              >
                <Input
                  id="confirmPassword"
                  placeholder="Repita a senha"
                  className="bg-card h-11"
                  type="password"
                  autoComplete="new-password"
                  value={form.confirmPassword}
                  aria-invalid={Boolean(errors.confirmPassword)}
                  onChange={(event) => update("confirmPassword", event.target.value)}
                />
              </FormField>
              <Button type="submit" className="h-12 sm:col-span-2" disabled={loading}>
                {loading ? <Loader2Icon className="animate-spin" /> : null}
                {loading ? "Concluindo…" : "Concluir cadastro"}
              </Button>
            </form>
          </div>
        )}
      </AuthShell>
    </AuthenticatedRedirect>
  );
}

function InvitationSummary({ invitation }: { invitation: OfficeMemberInvitation }) {
  const expiresAt = new Intl.DateTimeFormat("pt-BR", {
    dateStyle: "long",
    timeStyle: "short",
  }).format(new Date(invitation.expires_at));
  return (
    <section className="border-lime/35 bg-lime-pale/55 rounded-2xl border p-4 sm:p-5">
      <div className="flex gap-3">
        <span className="bg-lime text-pine grid size-10 shrink-0 place-items-center rounded-xl">
          <Building2Icon className="size-5" />
        </span>
        <div className="min-w-0">
          <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">
            Você foi convidado(a) para
          </p>
          <h2 className="font-display mt-1 text-lg font-semibold">
            Gabinete de {invitation.councillor_name}
          </h2>
          <p className="text-ink-soft mt-1 text-sm">
            {[invitation.party, `${invitation.city}/${invitation.state}`]
              .filter(Boolean)
              .join(" · ")}
          </p>
        </div>
      </div>
      <p className="border-lime/25 text-ink-soft mt-4 flex items-center gap-2 border-t pt-3 text-xs">
        <MapPinIcon className="text-lime-deep size-3.5" />
        Convite válido até {expiresAt}.
      </p>
    </section>
  );
}

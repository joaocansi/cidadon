"use client";

import { CheckCircle2Icon, EyeIcon, EyeOffIcon, Loader2Icon } from "lucide-react";
import Link from "next/link";
import { type FormEvent, useState } from "react";

import { showToast } from "@/components/shared/toast";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { AuthShell } from "@/features/auth/components/auth-shell";
import { AuthenticatedRedirect } from "@/features/auth/components/authenticated-redirect";
import { apiRegisterCitizen } from "@/lib/api";
import {
  apiErrorMessage,
  apiFieldErrors,
  type FieldErrors,
  invalidEmail,
  requiredFields,
} from "@/lib/forms";

const PASSWORD_MIN_LENGTH = 6;
const PASSWORD_MAX_LENGTH = 72;

export default function RegisterPage() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [city, setCity] = useState("");
  const [state, setState] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [errors, setErrors] = useState<FieldErrors>({});
  const [loading, setLoading] = useState(false);
  const [done, setDone] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const nextErrors = requiredFields(
      { name, email, password, city, state },
      {
        name: "Nome completo",
        email: "E-mail",
        password: "Senha",
        city: "Cidade",
        state: "Estado",
      },
    );

    if (email && invalidEmail(email)) nextErrors.email = "Informe um e-mail válido.";
    if (password.length < PASSWORD_MIN_LENGTH || password.length > PASSWORD_MAX_LENGTH) {
      nextErrors.password = "A senha precisa ter entre 6 e 72 caracteres.";
    }
    if (password !== confirmPassword) {
      nextErrors.confirmPassword = "As senhas não coincidem.";
    }
    setErrors(nextErrors);
    if (Object.keys(nextErrors).length) return;

    setLoading(true);
    const result = await apiRegisterCitizen({
      name,
      email,
      password,
      city,
      state,
    });
    setLoading(false);

    if (!result.ok) {
      setErrors(apiFieldErrors(result.error));
      showToast(apiErrorMessage(result.error, "Não foi possível criar a conta."), "error");
      return;
    }

    showToast("Conta criada com sucesso. Agora você já pode acompanhar demandas.");
    setDone(true);
  }

  return (
    <AuthenticatedRedirect>
      <AuthShell
        wide
        eyebrow="Cadastro de cidadão"
        title={
          <>
            Crie sua conta e <em className="text-lime-deep not-italic">registre</em> sua primeira
            demanda
          </>
        }
        subtitle="É gratuito. Leva menos de um minuto e dá acesso ao acompanhamento de cada protocolo."
      >
        {done ? (
          <div className="border-line bg-card flex flex-col items-center gap-4 rounded-3xl border p-8 text-center">
            <span className="bg-lime-pale text-lime-deep flex size-14 items-center justify-center rounded-full">
              <CheckCircle2Icon className="size-7" />
            </span>
            <div>
              <h2 className="font-display text-ink text-xl font-semibold">
                Conta criada com sucesso!
              </h2>
              <p className="text-ink-soft mt-2 text-[14.5px]">
                Sua conta de cidadão está pronta. Agora é só entrar e começar a registrar demandas
                do seu bairro.
              </p>
            </div>
            <Link
              href="/login"
              className="bg-primary text-primary-foreground hover:bg-primary/80 inline-flex h-12 w-full items-center justify-center rounded-full px-6 text-[15px] font-medium transition-all"
            >
              Entrar na conta
            </Link>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="grid gap-5 sm:grid-cols-2">
            <div className="flex flex-col gap-2 sm:col-span-2">
              <Label htmlFor="name">Nome completo</Label>
              <Input
                id="name"
                type="text"
                autoComplete="name"
                placeholder="Maria da Silva"
                className="bg-card h-12 rounded-2xl px-4 text-[15px]"
                aria-invalid={Boolean(errors.name)}
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
              {errors.name ? (
                <p className="text-destructive text-sm font-medium">{errors.name}</p>
              ) : null}
            </div>

            <div className="flex flex-col gap-2 sm:col-span-2">
              <Label htmlFor="email">E-mail</Label>
              <Input
                id="email"
                type="email"
                autoComplete="email"
                placeholder="voce@exemplo.com.br"
                className="bg-card h-12 rounded-2xl px-4 text-[15px]"
                aria-invalid={Boolean(errors.email)}
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
              {errors.email ? (
                <p className="text-destructive text-sm font-medium">{errors.email}</p>
              ) : null}
            </div>

            <div className="grid grid-cols-[minmax(0,1fr)_90px] gap-4 sm:col-span-2">
              <div className="flex flex-col gap-2">
                <Label htmlFor="city">Cidade</Label>
                <Input
                  id="city"
                  type="text"
                  autoComplete="address-level2"
                  placeholder="São Paulo"
                  className="bg-card h-12 rounded-2xl px-4 text-[15px]"
                  aria-invalid={Boolean(errors.city)}
                  value={city}
                  onChange={(e) => setCity(e.target.value)}
                />
                {errors.city ? (
                  <p className="text-destructive text-sm font-medium">{errors.city}</p>
                ) : null}
              </div>
              <div className="flex flex-col gap-2">
                <Label htmlFor="state">Estado</Label>
                <Input
                  id="state"
                  type="text"
                  autoComplete="address-level1"
                  placeholder="SP"
                  maxLength={2}
                  className="bg-card h-12 rounded-2xl px-4 text-[15px]"
                  aria-invalid={Boolean(errors.state)}
                  value={state}
                  onChange={(e) => setState(e.target.value.toUpperCase())}
                />
                {errors.state ? (
                  <p className="text-destructive text-sm font-medium">{errors.state}</p>
                ) : null}
              </div>
            </div>

            <div className="flex flex-col gap-2">
              <Label htmlFor="password">Senha</Label>
              <div className="relative">
                <Input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  autoComplete="new-password"
                  placeholder="Entre 6 e 72 caracteres"
                  minLength={PASSWORD_MIN_LENGTH}
                  maxLength={PASSWORD_MAX_LENGTH}
                  className="bg-card h-12 rounded-2xl px-4 pr-12 text-[15px]"
                  aria-invalid={Boolean(errors.password)}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                />
                <button
                  type="button"
                  aria-label={showPassword ? "Ocultar senha" : "Mostrar senha"}
                  onClick={() => setShowPassword((v) => !v)}
                  className="text-ink-faint hover:bg-ink/5 hover:text-ink absolute top-1/2 right-3 -translate-y-1/2 rounded-lg p-1.5 transition-colors"
                >
                  {showPassword ? (
                    <EyeOffIcon className="size-4.5" />
                  ) : (
                    <EyeIcon className="size-4.5" />
                  )}
                </button>
              </div>
              {errors.password ? (
                <p className="text-destructive text-sm font-medium">{errors.password}</p>
              ) : null}
            </div>

            <div className="flex flex-col gap-2">
              <Label htmlFor="confirmPassword">Confirmar senha</Label>
              <Input
                id="confirmPassword"
                type={showPassword ? "text" : "password"}
                autoComplete="new-password"
                placeholder="Repita a senha"
                minLength={PASSWORD_MIN_LENGTH}
                maxLength={PASSWORD_MAX_LENGTH}
                className="bg-card h-12 rounded-2xl px-4 text-[15px]"
                aria-invalid={Boolean(errors.confirmPassword)}
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
              />
              {errors.confirmPassword ? (
                <p className="text-destructive text-sm font-medium">{errors.confirmPassword}</p>
              ) : null}
            </div>

            <Button
              type="submit"
              disabled={loading}
              className="mt-1 h-12 rounded-full text-[15px] sm:col-span-2"
            >
              {loading ? <Loader2Icon className="animate-spin" /> : null}
              {loading ? "Criando conta…" : "Criar conta de cidadão"}
            </Button>

            <p className="text-ink-soft text-center text-[14px] sm:col-span-2">
              Já tem uma conta?{" "}
              <Link href="/login" className="text-lime-deep hover:text-pine font-semibold">
                Entrar
              </Link>
            </p>
          </form>
        )}
      </AuthShell>
    </AuthenticatedRedirect>
  );
}

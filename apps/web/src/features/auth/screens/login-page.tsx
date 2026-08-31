"use client";

import { EyeIcon, EyeOffIcon, Loader2Icon } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { type FormEvent, useState } from "react";

import { showToast } from "@/components/shared/toast";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { AuthShell } from "@/features/auth/components/auth-shell";
import { AuthenticatedRedirect } from "@/features/auth/components/authenticated-redirect";
import { useSession } from "@/features/auth/components/session-provider";
import { apiLogin } from "@/lib/api";
import { apiErrorMessage, apiFieldErrors, type FieldErrors, invalidEmail } from "@/lib/forms";

const PASSWORD_MIN_LENGTH = 6;
const PASSWORD_MAX_LENGTH = 72;

export default function LoginPage() {
  const router = useRouter();
  const { refresh } = useSession();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [errors, setErrors] = useState<FieldErrors>({});
  const [loading, setLoading] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const nextErrors: FieldErrors = {};

    if (!email.trim()) nextErrors.email = "Informe seu e-mail.";
    else if (invalidEmail(email)) nextErrors.email = "Informe um e-mail válido.";
    if (!password) nextErrors.password = "Informe sua senha.";
    if (password.length < PASSWORD_MIN_LENGTH || password.length > PASSWORD_MAX_LENGTH) {
      nextErrors.password = "A senha precisa ter entre 6 e 72 caracteres.";
    }
    setErrors(nextErrors);
    if (Object.keys(nextErrors).length) return;

    setLoading(true);
    const result = await apiLogin(email, password);
    setLoading(false);

    if (!result.ok) {
      setErrors(apiFieldErrors(result.error));
      showToast(apiErrorMessage(result.error, "Não foi possível entrar."), "error");
      return;
    }

    showToast("Login realizado. Bem-vindo(a) ao Cidadon!");
    await refresh();
    router.replace(result.data?.role === "citizen" ? "/demands" : "/office");
    router.refresh();
  }

  return (
    <AuthenticatedRedirect>
      <AuthShell
        eyebrow="Boas-vindas de volta"
        title={
          <>
            Entre para <em className="text-lime-deep not-italic">acompanhar</em> suas demandas
          </>
        }
        subtitle="Acesse sua conta para ver cada protocolo, comentar e apoiar as ruas da sua cidade."
      >
        <form onSubmit={handleSubmit} className="flex flex-col gap-5">
          <div className="flex flex-col gap-2">
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

          <div className="flex flex-col gap-2">
            <div className="flex items-center justify-between">
              <Label htmlFor="password">Senha</Label>
              <Link href="/" className="text-lime-deep hover:text-pine text-[13px] font-semibold">
                Esqueci minha senha
              </Link>
            </div>
            <div className="relative">
              <Input
                id="password"
                type={showPassword ? "text" : "password"}
                autoComplete="current-password"
                placeholder="•••••••"
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

          <Button type="submit" disabled={loading} className="mt-1 h-12 rounded-full text-[15px]">
            {loading ? <Loader2Icon className="animate-spin" /> : null}
            {loading ? "Entrando…" : "Entrar"}
          </Button>

          <p className="text-ink-soft text-center text-[14px]">
            Ainda não tem conta?{" "}
            <Link href="/register" className="text-lime-deep hover:text-pine font-semibold">
              Criar conta de cidadão
            </Link>
            <span className="mx-1">ou</span>
            <Link
              href="/register/councillor"
              className="text-lime-deep hover:text-pine font-semibold"
            >
              criar conta de vereador
            </Link>
          </p>
        </form>
      </AuthShell>
    </AuthenticatedRedirect>
  );
}

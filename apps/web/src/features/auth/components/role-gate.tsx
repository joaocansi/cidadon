"use client";

import { useRouter } from "next/navigation";
import { type ReactNode, useEffect } from "react";

import { useSession } from "@/features/auth/components/session-provider";
import type { UserRole } from "@/lib/api";

export function RoleGate({ allowed, children }: { allowed: UserRole[]; children: ReactNode }) {
  const router = useRouter();
  const { user, status } = useSession();
  const allowedKey = allowed.join(",");
  useEffect(() => {
    if (status === "anonymous") router.replace("/entrar");
    else if (user && !allowedKey.split(",").includes(user.role))
      router.replace(user.role === "citizen" ? "/demandas" : "/gabinete");
  }, [allowedKey, router, status, user]);
  if (status === "loading" || !user || !allowed.includes(user.role))
    return (
      <div className="bg-paper text-ink-soft grid min-h-dvh place-items-center text-sm">
        Verificando seu acesso…
      </div>
    );
  return <>{children}</>;
}

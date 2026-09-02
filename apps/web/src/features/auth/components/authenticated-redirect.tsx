"use client";

import { useRouter } from "next/navigation";
import { type ReactNode, useEffect } from "react";

import { useSession } from "@/features/auth/components/session-provider";

const destination = (role: string) => (role === "citizen" ? "/demandas" : "/gabinete");

export function AuthenticatedRedirect({ children }: { children: ReactNode }) {
  const router = useRouter();
  const { user, status } = useSession();
  useEffect(() => {
    if (user) router.replace(destination(user.role));
  }, [router, user]);
  if (status === "loading" || user)
    return (
      <div className="bg-paper text-ink-soft grid min-h-dvh place-items-center text-sm">
        Carregando…
      </div>
    );
  return <>{children}</>;
}

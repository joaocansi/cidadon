"use client";

import { Building2Icon, ClipboardListIcon, PlusIcon } from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import type { ReactNode } from "react";

import { ProfileMenu } from "@/features/auth/components/profile-menu";
import { useSession } from "@/features/auth/components/session-provider";
import { NotificationBell } from "@/features/notifications/components/notification-bell";
import { cn } from "@/lib/utils";

export function CitizenShell({
  title,
  subtitle,
  children,
}: {
  title: string;
  subtitle?: string;
  children: ReactNode;
}) {
  const pathname = usePathname();
  const { user } = useSession();
  const links =
    user?.role === "citizen"
      ? [
          { href: "/demands", label: "Demandas", icon: ClipboardListIcon },
          { href: "/offices", label: "Gabinetes", icon: Building2Icon },
        ]
      : [{ href: "/offices", label: "Gabinetes", icon: Building2Icon }];
  const home = user?.role === "citizen" ? "/demands" : user ? "/office" : "/";
  const eyebrow =
    user?.role === "citizen" ? "Área do cidadão" : user ? "Visão pública" : "Diretório público";
  return (
    <div className="bg-paper text-ink min-h-screen">
      <header className="border-line/80 bg-paper/92 sticky top-0 z-40 border-b backdrop-blur-xl">
        <div className="mx-auto flex h-14 max-w-[1440px] items-center justify-between gap-3 px-4 sm:h-16 sm:px-6 lg:px-8">
          <Link
            href={home}
            className="font-display flex shrink-0 items-center gap-2 text-lg font-semibold"
          >
            <span className="bg-lime text-pine grid size-7 place-items-center rounded-md">
              <ClipboardListIcon className="size-4" />
            </span>
            <span className="hidden sm:inline">Cidadon</span>
          </Link>
          <nav className="flex items-center gap-1" aria-label="Navegação do cidadão">
            {links.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className={cn(
                  "text-ink-soft hover:bg-ink/5 hover:text-ink flex h-9 items-center gap-2 rounded-lg px-2.5 text-sm font-semibold transition",
                  pathname.startsWith(link.href) && "bg-lime-pale text-pine",
                )}
              >
                <link.icon className="size-4" />
                <span className="hidden sm:inline">{link.label}</span>
              </Link>
            ))}
          </nav>
          <div className="flex items-center gap-2">
            {user ? <NotificationBell /> : null}
            {user?.role === "citizen" ? (
              <Link
                href="/demands/new"
                className="bg-pine text-paper inline-flex h-9 items-center gap-2 rounded-lg px-3 text-sm font-semibold"
              >
                <PlusIcon className="size-4" />
                <span className="hidden sm:inline">Nova demanda</span>
              </Link>
            ) : user ? (
              <Link
                href="/office"
                className="bg-pine text-paper inline-flex h-9 items-center rounded-lg px-3 text-sm font-semibold"
              >
                Minha área
              </Link>
            ) : (
              <Link
                href="/login"
                className="bg-pine text-paper inline-flex h-9 items-center rounded-lg px-3 text-sm font-semibold"
              >
                Entrar
              </Link>
            )}
            {user ? <ProfileMenu /> : null}
          </div>
        </div>
      </header>
      <main className="mx-auto w-full max-w-[1440px] px-4 py-5 sm:px-6 sm:py-6 lg:px-8">
        <div className="mb-5 flex flex-wrap items-end justify-between gap-3">
          <div>
            <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">
              {eyebrow}
            </p>
            <h1 className="font-display mt-0.5 text-2xl font-semibold sm:text-3xl">{title}</h1>
            {subtitle ? <p className="text-ink-soft mt-1 max-w-2xl text-sm">{subtitle}</p> : null}
          </div>
        </div>
        {children}
      </main>
    </div>
  );
}

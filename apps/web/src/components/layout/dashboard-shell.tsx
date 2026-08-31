"use client";

import {
  Building2Icon,
  ClipboardListIcon,
  Clock3Icon,
  ExternalLinkIcon,
  MenuIcon,
  UsersIcon,
  XIcon,
} from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { type ReactNode, useState } from "react";

import { Button } from "@/components/ui/button";
import { ProfileMenu } from "@/features/auth/components/profile-menu";
import { useSession } from "@/features/auth/components/session-provider";
import { NotificationBell } from "@/features/notifications/components/notification-bell";
import { cn } from "@/lib/utils";

export function DashboardShell({
  title,
  subtitle,
  actions,
  officeId,
  children,
}: {
  title: string;
  subtitle?: string;
  actions?: ReactNode;
  officeId?: number;
  children: ReactNode;
}) {
  const pathname = usePathname();
  const { user } = useSession();
  const [mobileNavigationOpen, setMobileNavigationOpen] = useState(false);
  const navItems = [
    { label: "Demandas", href: "/office", icon: ClipboardListIcon },
    { label: "Acompanhamento", icon: Clock3Icon, pending: true },
    ...(user?.role === "councillor"
      ? [
          { label: "Equipe", href: "/office/team", icon: UsersIcon },
          { label: "Gabinete", href: "/office/profile", icon: Building2Icon },
        ]
      : []),
    {
      label: "Perfil público",
      href: officeId ? `/offices/${officeId}` : "/offices",
      icon: Building2Icon,
    },
  ];

  const navigation = navItems.map((item) => {
    const href = item.href ?? "";
    const active = href.startsWith("/offices")
      ? pathname.startsWith("/offices")
      : href === "/office"
        ? pathname === "/office"
        : href.includes("#")
          ? false
          : pathname === href;
    const isPublicProfile = item.label === "Perfil público";
    if (item.pending)
      return (
        <span
          key={item.label}
          aria-disabled="true"
          className="text-paper/35 flex h-10 shrink-0 cursor-not-allowed items-center gap-2.5 rounded-xl px-3 text-sm font-semibold"
        >
          <item.icon className="size-4" />
          {item.label}
          <span className="ml-auto text-[10px] font-bold tracking-wide uppercase">Em breve</span>
        </span>
      );
    return (
      <Link
        key={item.label}
        href={href}
        target={isPublicProfile ? "_blank" : undefined}
        rel={isPublicProfile ? "noreferrer" : undefined}
        aria-current={active ? "page" : undefined}
        onClick={() => setMobileNavigationOpen(false)}
        className={cn(
          "flex h-10 shrink-0 items-center gap-2.5 rounded-xl px-3 text-sm font-semibold transition",
          active ? "bg-lime text-pine" : "text-paper/70 hover:bg-paper/10 hover:text-paper",
        )}
      >
        <item.icon className="size-4" />
        {item.label}
        {isPublicProfile ? (
          <ExternalLinkIcon
            className="ml-auto size-3.5 opacity-75"
            aria-label="Abrir perfil público em nova aba"
          />
        ) : null}
      </Link>
    );
  });

  return (
    <div className="bg-paper text-ink min-h-dvh lg:grid lg:grid-cols-[248px_minmax(0,1fr)]">
      <aside className="bg-pine text-paper hidden lg:sticky lg:top-0 lg:flex lg:h-dvh lg:flex-col">
        <div className="border-paper/10 border-b px-5 py-5">
          <Link href="/office" className="inline-flex items-center gap-3">
            <span className="bg-lime text-pine grid size-9 place-items-center rounded-xl">
              <ClipboardListIcon className="size-4" />
            </span>
            <span>
              <span className="font-display block text-xl leading-none font-semibold">Cidadon</span>
              <span className="text-paper/50 mt-1 block text-[11px] font-semibold tracking-wide uppercase">
                Gabinete
              </span>
            </span>
          </Link>
        </div>
        <nav className="flex flex-1 flex-col gap-1.5 p-3" aria-label="Navegação do gabinete">
          {navigation}
        </nav>
        <div className="border-paper/10 text-paper/50 border-t p-5 text-xs leading-5">
          Atuação pública, organizada e rastreável.
        </div>
      </aside>
      <div className="min-w-0">
        <header className="border-line/80 bg-paper/92 sticky top-0 z-30 border-b backdrop-blur-xl">
          <div className="flex h-14 items-center justify-between gap-3 px-4 sm:h-16 sm:px-6 lg:px-8">
            <div className="min-w-0">
              <p className="text-lime-deep text-[10px] font-semibold tracking-wide uppercase">
                Painel operacional
              </p>
              <h1 className="font-display truncate text-lg leading-tight font-semibold sm:text-xl">
                {title}
              </h1>
              {subtitle ? (
                <p className="text-ink-soft hidden truncate text-xs xl:block">{subtitle}</p>
              ) : null}
            </div>
            <div className="flex shrink-0 items-center gap-2">
              {actions}
              <NotificationBell />
              <Button
                type="button"
                size="icon"
                variant="outline"
                className="size-9 lg:hidden"
                aria-label="Abrir navegação"
                aria-expanded={mobileNavigationOpen}
                aria-controls="mobile-office-navigation"
                onClick={() => setMobileNavigationOpen(true)}
              >
                <MenuIcon />
              </Button>
              <ProfileMenu />
            </div>
          </div>
        </header>
        {mobileNavigationOpen ? (
          <>
            <button
              type="button"
              aria-label="Fechar navegação"
              className="bg-pine/35 fixed inset-0 z-40 backdrop-blur-[1px] lg:hidden"
              onClick={() => setMobileNavigationOpen(false)}
            />
            <aside
              id="mobile-office-navigation"
              className="bg-pine text-paper fixed inset-y-0 left-0 z-50 flex w-[min(20rem,calc(100vw-2.5rem))] flex-col shadow-2xl lg:hidden"
            >
              <div className="border-paper/10 flex items-center justify-between border-b px-5 py-5">
                <Link
                  href="/office"
                  onClick={() => setMobileNavigationOpen(false)}
                  className="inline-flex items-center gap-3"
                >
                  <span className="bg-lime text-pine grid size-9 place-items-center rounded-xl">
                    <ClipboardListIcon className="size-4" />
                  </span>
                  <span>
                    <span className="font-display block text-xl leading-none font-semibold">
                      Cidadon
                    </span>
                    <span className="text-paper/50 mt-1 block text-[11px] font-semibold tracking-wide uppercase">
                      Gabinete
                    </span>
                  </span>
                </Link>
                <Button
                  type="button"
                  size="icon"
                  variant="ghost"
                  className="text-paper hover:bg-paper/10 hover:text-paper"
                  aria-label="Fechar navegação"
                  onClick={() => setMobileNavigationOpen(false)}
                >
                  <XIcon />
                </Button>
              </div>
              <nav className="flex flex-1 flex-col gap-1.5 p-3" aria-label="Navegação do gabinete">
                {navigation}
              </nav>
              <div className="border-paper/10 text-paper/50 border-t p-5 text-xs leading-5">
                Atuação pública, organizada e rastreável.
              </div>
            </aside>
          </>
        ) : null}
        <main className="mx-auto w-full max-w-[1440px] px-4 py-5 sm:px-6 sm:py-6 lg:px-8">
          {children}
        </main>
      </div>
    </div>
  );
}

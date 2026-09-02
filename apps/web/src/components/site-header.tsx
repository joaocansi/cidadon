"use client";

import { LogOutIcon, MenuIcon, XIcon } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import { showToast } from "@/components/shared/toast";
import { buttonVariants } from "@/components/ui/button";
import { useSession } from "@/features/auth/components/session-provider";
import { cn } from "@/lib/utils";

const NAV_LINKS = [
  { href: "/#como-funciona", label: "Como funciona" },
  { href: "/#cidadaos", label: "Para cidadãos" },
  { href: "/#vereadores", label: "Para vereadores" },
  { href: "/demandas", label: "Demandas" },
  { href: "/#faq", label: "Perguntas frequentes" },
];

export function SiteHeader() {
  const [scrolled, setScrolled] = useState(false);
  const [open, setOpen] = useState(false);
  const { user, signOut } = useSession();
  const router = useRouter();

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 8);
    onScroll();
    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  useEffect(() => {
    document.body.style.overflow = open ? "hidden" : "";
    return () => {
      document.body.style.overflow = "";
    };
  }, [open]);

  async function logout() {
    if (await signOut()) {
      showToast("Você saiu da sua conta.");
      setOpen(false);
      router.replace("/");
      router.refresh();
    }
  }
  const areaHref = user?.role === "citizen" ? "/demandas" : "/gabinete";

  return (
    <header
      className={cn(
        "bg-paper/85 sticky top-0 z-40 border-b backdrop-blur-xl transition-[border-color,box-shadow]",
        scrolled ? "border-line" : "border-transparent",
      )}
      style={scrolled ? { boxShadow: "0 4px 16px -12px rgba(30,36,23,0.3)" } : undefined}
    >
      <div className="mx-auto flex h-16 max-w-(--container) items-center justify-between px-5 sm:px-7">
        <Link
          href="/"
          className="font-display text-ink flex items-center gap-2.5 text-[22px] font-bold"
        >
          <span className="logo-mark" />
          cidadon
        </Link>

        <nav className="hidden items-center gap-8 lg:flex">
          {NAV_LINKS.map((link) => (
            <Link
              key={`desktop-${link.href}`}
              href={link.href}
              className="text-ink-soft hover:text-ink text-[15px] font-semibold transition-colors"
            >
              {link.label}
            </Link>
          ))}
        </nav>

        <div className="flex items-center gap-4">
          {user ? (
            <>
              <Link
                href={areaHref}
                className={cn(
                  buttonVariants({ variant: "default", size: "sm" }),
                  "hidden h-9 rounded-full px-5 lg:inline-flex",
                )}
              >
                Minha área
              </Link>
              <button
                type="button"
                onClick={() => void logout()}
                className="text-ink-soft hover:text-ink hidden h-9 items-center gap-2 text-[15px] font-semibold lg:inline-flex"
              >
                <LogOutIcon className="size-4" />
                Sair
              </button>
            </>
          ) : (
            <>
              <Link
                href="/entrar"
                className="text-ink-soft hover:text-ink hidden text-[15px] font-semibold transition-colors lg:inline-block"
              >
                Entrar
              </Link>
              <Link
                href="/cadastro"
                className={cn(
                  buttonVariants({ variant: "default", size: "sm" }),
                  "h-9 rounded-full px-5",
                )}
              >
                Criar conta
              </Link>
            </>
          )}
          <button
            type="button"
            aria-label={open ? "Fechar menu" : "Abrir menu"}
            aria-expanded={open}
            onClick={() => setOpen((o) => !o)}
            className="hover:bg-ink/5 flex size-10 items-center justify-center rounded-[10px] transition-colors lg:hidden"
          >
            {open ? (
              <XIcon className="text-ink size-5" />
            ) : (
              <MenuIcon className="text-ink size-5" />
            )}
          </button>
        </div>
      </div>

      <div
        className={cn(
          "bg-paper fixed inset-x-0 top-16 bottom-0 z-40 overflow-y-auto px-5 pt-3 transition-all duration-300 sm:px-7 lg:hidden",
          open
            ? "pointer-events-auto translate-y-0 opacity-100"
            : "pointer-events-none -translate-y-2 opacity-0",
        )}
      >
        {NAV_LINKS.map((link) => (
          <Link
            key={`mobile-${link.href}`}
            href={link.href}
            onClick={() => setOpen(false)}
            className="border-line-soft text-ink block border-b py-4 text-xl font-semibold"
          >
            {link.label}
          </Link>
        ))}
        <div className="mt-6 flex flex-col gap-3">
          {user ? (
            <>
              <Link
                href={areaHref}
                onClick={() => setOpen(false)}
                className={cn(buttonVariants({ variant: "default" }), "h-12 w-full rounded-full")}
              >
                Minha área
              </Link>
              <button
                type="button"
                onClick={() => void logout()}
                className={cn(buttonVariants({ variant: "outline" }), "h-12 w-full rounded-full")}
              >
                Sair
              </button>
            </>
          ) : (
            <>
              <Link
                href="/entrar"
                onClick={() => setOpen(false)}
                className={cn(buttonVariants({ variant: "outline" }), "h-12 w-full rounded-full")}
              >
                Entrar
              </Link>
              <Link
                href="/cadastro"
                onClick={() => setOpen(false)}
                className={cn(buttonVariants({ variant: "default" }), "h-12 w-full rounded-full")}
              >
                Criar conta
              </Link>
            </>
          )}
        </div>
      </div>
    </header>
  );
}

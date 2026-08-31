"use client";

import { ChevronDownIcon, LogOutIcon } from "lucide-react";
import { useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";

import { showToast } from "@/components/shared/toast";
import { UserAvatar } from "@/components/shared/user-avatar";
import { useSession } from "@/features/auth/components/session-provider";
import { cn } from "@/lib/utils";

const labels = {
  citizen: "Cidadão",
  councillor: "Vereador(a)",
  office_member: "Membro do gabinete",
};

export function ProfileMenu() {
  const router = useRouter();
  const ref = useRef<HTMLDivElement>(null);
  const [open, setOpen] = useState(false);
  const { user, signOut } = useSession();

  useEffect(() => {
    const close = (event: MouseEvent) => {
      if (ref.current && !ref.current.contains(event.target as Node)) setOpen(false);
    };
    window.addEventListener("mousedown", close);
    return () => window.removeEventListener("mousedown", close);
  }, []);

  async function logout() {
    if (!(await signOut())) {
      showToast("Não foi possível encerrar sua sessão.", "error");
      return;
    }
    showToast("Você saiu da sua conta.");
    router.replace("/");
    router.refresh();
  }

  return (
    <div className="relative" ref={ref}>
      <button
        type="button"
        onClick={() => setOpen((value) => !value)}
        aria-expanded={open}
        aria-label="Abrir menu de perfil"
        className="border-line bg-card text-ink hover:border-lime/60 flex h-10 items-center gap-2 rounded-full border p-1 pr-2.5 text-sm font-semibold transition"
      >
        <UserAvatar name={user?.name} imageUrl={user?.image_url} className="bg-lime size-8" />
        <span className="hidden xl:inline">{user?.name?.split(" ")[0] || "Perfil"}</span>
        <ChevronDownIcon className={cn("size-3.5 transition-transform", open && "rotate-180")} />
      </button>
      {open ? (
        <div className="border-line bg-card text-ink absolute top-[calc(100%+8px)] right-0 z-50 w-64 rounded-xl border p-1.5 shadow-[0_14px_32px_-16px_rgba(0,0,0,.45)]">
          <div className="border-line-soft border-b px-3 py-2.5">
            <p className="truncate text-sm font-semibold">{user?.name}</p>
            <p className="text-ink-soft mt-0.5 truncate text-xs">{user?.email}</p>
            <p className="text-lime-deep mt-1 text-xs font-semibold">
              {user ? labels[user.role] : ""}
            </p>
          </div>
          <button
            type="button"
            onClick={() => void logout()}
            className="hover:bg-destructive/10 hover:text-destructive mt-1 flex h-9 w-full items-center gap-2 rounded-lg px-3 text-left text-sm font-semibold transition"
          >
            <LogOutIcon className="size-4" />
            Sair da conta
          </button>
        </div>
      ) : null}
    </div>
  );
}

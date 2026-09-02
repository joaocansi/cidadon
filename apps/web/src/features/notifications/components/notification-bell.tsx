"use client";

import { BellIcon, CheckIcon } from "lucide-react";
import Link from "next/link";
import { useEffect, useRef, useState } from "react";

import { showToast } from "@/components/shared/toast";
import { useSession } from "@/features/auth/components/session-provider";
import { apiNotifications, apiReadNotifications, type Notification } from "@/lib/api";

const copy: Record<string, string> = {
  demand_created: "Nova demanda direcionada ao gabinete",
  demand_claimed: "Uma demanda foi assumida",
  demand_execution_started: "A execução foi iniciada",
  demand_confirmation_requested: "Uma demanda aguarda validação",
  demand_confirmed: "A conclusão foi confirmada",
  demand_reopened: "Uma demanda foi reaberta",
  demand_commented: "Há uma nova atualização",
  demand_automatically_completed: "Uma demanda foi concluída pelo prazo",
  demand_milestone: "O gabinete publicou um novo marco",
};
export function NotificationBell() {
  const { status, user } = useSession();
  const [items, setItems] = useState<Notification[]>([]);
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);
  const knownIDs = useRef(new Set<number>());
  const latestID = useRef(0);
  useEffect(() => {
    if (status !== "authenticated" || !user) return;
    let active = true;
    let stream: EventSource | undefined;
    void (async () => {
      const result = await apiNotifications();
      if (!active) return;
      const initial = result.ok ? (result.data ?? []) : [];
      knownIDs.current = new Set(initial.map((item) => item.id));
      latestID.current = initial.reduce((latest, item) => Math.max(latest, item.id), 0);
      setItems(initial);

      stream = new EventSource(`/api/notifications/stream?after=${latestID.current}`);
      stream.addEventListener("notification", (event) => {
        try {
          const next = JSON.parse(event.data) as Notification;
          if (knownIDs.current.has(next.id)) return;
          knownIDs.current.add(next.id);
          latestID.current = Math.max(latestID.current, next.id);
          setItems((current) => [next, ...current]);
          showToast(copy[next.type] ?? "Há uma atualização em uma demanda.");
        } catch {
          // Invalid stream payloads should not interrupt the authenticated UI.
        }
      });
      stream.addEventListener("error", () => {
        // EventSource reconnects by itself. Previously seen notification IDs are
        // kept across reconnects, so a replay can never show a duplicate toast.
      });
    })();
    return () => {
      active = false;
      stream?.close();
    };
  }, [status, user]);
  useEffect(() => {
    const close = (event: MouseEvent) => {
      if (ref.current && !ref.current.contains(event.target as Node)) setOpen(false);
    };
    window.addEventListener("mousedown", close);
    return () => window.removeEventListener("mousedown", close);
  }, []);
  const unread = items.filter((item) => !item.read_at).length;
  async function markAllRead() {
    const result = await apiReadNotifications();
    if (result.ok)
      setItems((current) =>
        current.map((item) => ({ ...item, read_at: item.read_at ?? new Date().toISOString() })),
      );
  }
  return (
    <div className="relative" ref={ref}>
      <button
        type="button"
        onClick={() => setOpen((value) => !value)}
        aria-label="Abrir notificações"
        aria-expanded={open}
        className="text-ink-soft hover:bg-ink/5 hover:text-ink relative grid size-9 place-items-center rounded-lg transition"
      >
        <BellIcon className="size-4" />
        {unread ? (
          <span className="bg-destructive absolute top-0.5 right-0.5 grid min-w-4 place-items-center rounded-full px-1 text-[10px] leading-4 font-bold text-white">
            {unread > 9 ? "9+" : unread}
          </span>
        ) : null}
      </button>
      {open ? (
        <div className="border-line bg-card absolute top-[calc(100%+8px)] right-0 z-50 w-[min(22rem,calc(100vw-2rem))] overflow-hidden rounded-xl border shadow-[0_18px_40px_-18px_rgba(0,0,0,.4)]">
          <div className="border-line-soft flex items-center justify-between border-b px-4 py-3">
            <div>
              <p className="text-sm font-semibold">Notificações</p>
              <p className="text-ink-soft text-xs">
                {unread ? `${unread} não lida${unread > 1 ? "s" : ""}` : "Tudo em dia"}
              </p>
            </div>
            {unread ? (
              <button
                type="button"
                onClick={() => void markAllRead()}
                className="text-pine inline-flex items-center gap-1 text-xs font-semibold hover:underline"
              >
                <CheckIcon className="size-3.5" />
                Ler todas
              </button>
            ) : null}
          </div>
          <div className="max-h-80 overflow-y-auto">
            {items.length ? (
              items.slice(0, 10).map((item) => (
                <Link
                  key={item.id}
                  href={`/demandas/${item.demand_id}`}
                  onClick={() => setOpen(false)}
                  className={`border-line-soft hover:bg-paper-2 block border-b px-4 py-3 transition ${item.read_at ? "" : "bg-lime-pale/35"}`}
                >
                  <p className="text-ink text-sm font-semibold">
                    {copy[item.type] ?? "Atualização em uma demanda"}
                  </p>
                  <p className="text-ink-soft mt-1 text-xs">
                    {new Date(item.created_at).toLocaleString("pt-BR")}
                  </p>
                </Link>
              ))
            ) : (
              <p className="text-ink-soft px-4 py-8 text-center text-sm">
                Nenhuma notificação por enquanto.
              </p>
            )}
          </div>
        </div>
      ) : null}
    </div>
  );
}

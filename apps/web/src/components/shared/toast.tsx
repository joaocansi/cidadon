"use client";

import { AlertCircleIcon, CheckCircle2Icon } from "lucide-react";
import { useEffect, useState } from "react";

import { cn } from "@/lib/utils";

const TOAST_EVENT = "cidadon:toast";

type Toast = { message: string; variant?: "success" | "error" };
export function showToast(message: string, variant: Toast["variant"] = "success") {
  if (typeof window === "undefined") return;
  window.dispatchEvent(
    new CustomEvent(TOAST_EVENT, { detail: { message, variant } satisfies Toast }),
  );
}

export function Toaster() {
  const [toast, setToast] = useState<Toast | null>(null);

  useEffect(() => {
    let timer: ReturnType<typeof setTimeout> | undefined;
    const handler = (event: Event) => {
      setToast((event as CustomEvent<Toast>).detail);
      clearTimeout(timer);
      timer = setTimeout(() => setToast(null), 3600);
    };
    window.addEventListener(TOAST_EVENT, handler);
    return () => {
      window.removeEventListener(TOAST_EVENT, handler);
      if (timer) clearTimeout(timer);
    };
  }, []);

  return (
    <div
      role="status"
      aria-live="polite"
      className={cn(
        "text-paper fixed bottom-7 left-1/2 z-[200] flex -translate-x-1/2 items-center gap-2.5 rounded-2xl px-5 py-3 text-sm font-semibold shadow-[0_14px_36px_-10px_rgba(0,0,0,0.4)] transition-all duration-300",
        toast?.variant === "error" ? "bg-destructive" : "bg-pine",
        toast ? "translate-y-0 opacity-100" : "pointer-events-none translate-y-5 opacity-0",
      )}
    >
      {toast?.variant === "error" ? (
        <AlertCircleIcon className="size-4.5" />
      ) : (
        <CheckCircle2Icon className="text-lime size-4.5" />
      )}
      <span>{toast?.message}</span>
    </div>
  );
}

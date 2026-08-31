import Link from "next/link";
import type { ReactNode } from "react";

import { cn } from "@/lib/utils";

export function AuthShell({
  eyebrow,
  title,
  subtitle,
  children,
  wide = false,
}: {
  eyebrow: ReactNode;
  title: ReactNode;
  subtitle?: ReactNode;
  children: ReactNode;
  wide?: boolean;
}) {
  return (
    <div className="bg-paper relative flex min-h-dvh flex-col overflow-hidden">
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0 -z-10 bg-[linear-gradient(var(--line-soft)_1px,transparent_1px),linear-gradient(90deg,var(--line-soft)_1px,transparent_1px)] [mask-image:radial-gradient(ellipse_90%_60%_at_50%_0%,black_40%,transparent_85%)] bg-[size:56px_56px] opacity-70"
      />
      <div className="wrap flex items-center justify-between pt-7">
        <Link
          href="/"
          className="font-display text-ink flex items-center gap-2.5 text-[22px] font-bold"
        >
          <span className="logo-mark" />
          cidadon
        </Link>
      </div>

      <main
        className={cn(
          "mx-auto flex w-full flex-1 flex-col justify-center px-5 py-10 sm:py-12",
          wide ? "max-w-2xl" : "max-w-md",
        )}
      >
        <div className="mb-8">
          <span className="eyebrow">
            <span className="dot" />
            {eyebrow}
          </span>
          <h1 className="font-display text-ink mt-5 text-[clamp(28px,4vw,38px)] leading-tight font-semibold tracking-tight">
            {title}
          </h1>
          {subtitle ? (
            <p className="text-ink-soft mt-3 text-[15.5px] leading-relaxed">{subtitle}</p>
          ) : null}
        </div>

        {children}
      </main>

      <footer className="wrap pb-7">
        <p className="text-ink-faint text-center text-[13px]">
          Cidadon — sua demanda tem um caminho até o gabinete.
        </p>
      </footer>
    </div>
  );
}

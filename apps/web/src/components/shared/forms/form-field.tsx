import type { ReactNode } from "react";

import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";

export function FormField({
  id,
  label,
  hint,
  error,
  children,
  className,
}: {
  id?: string;
  label: string;
  hint?: string;
  error?: string;
  children: ReactNode;
  className?: string;
}) {
  const messageId = id ? `${id}-message` : undefined;
  return (
    <div className={cn("flex min-w-0 flex-col gap-1.5", className)}>
      <Label htmlFor={id} className="text-ink text-sm font-semibold">
        {label}
      </Label>
      {children}
      {error ? (
        <p id={messageId} role="alert" className="text-destructive text-xs font-semibold">
          {error}
        </p>
      ) : hint ? (
        <p id={messageId} className="text-ink-faint text-xs leading-5">
          {hint}
        </p>
      ) : null}
    </div>
  );
}

"use client";

import { Loader2Icon, SendIcon, XIcon } from "lucide-react";
import { type FormEvent, useEffect, useState } from "react";

import { ImageAttachmentPicker } from "@/components/shared/image-attachment-picker";
import { showToast } from "@/components/shared/toast";
import { Button } from "@/components/ui/button";

export type TimelinePayload = { message: string; images: File[] };

type DemandTransitionDialogProps = {
  open: boolean;
  title: string;
  description: string;
  submitLabel: string;
  onClose: () => void;
  onSubmit: (payload: TimelinePayload) => Promise<boolean>;
};

// One dialog is shared by every operational status change. Keeping the
// mandatory message in this component prevents an action from bypassing the
// public operational record.
export function DemandTransitionDialog({
  open,
  title,
  description,
  submitLabel,
  onClose,
  onSubmit,
}: DemandTransitionDialogProps) {
  const [message, setMessage] = useState("");
  const [images, setImages] = useState<File[]>([]);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (!open) return;
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape" && !saving) onClose();
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [onClose, open, saving]);

  if (!open) return null;

  async function submit(event: FormEvent) {
    event.preventDefault();
    if (message.trim().length < 3) {
      showToast("Escreva uma justificativa com pelo menos três caracteres.", "error");
      return;
    }
    setSaving(true);
    const completed = await onSubmit({ message: message.trim(), images });
    setSaving(false);
    if (!completed) return;
    setMessage("");
    setImages([]);
    onClose();
  }

  return (
    <div
      role="dialog"
      aria-modal="true"
      aria-labelledby="demand-transition-title"
      className="bg-ink/45 fixed inset-0 z-[90] grid place-items-center p-4 backdrop-blur-sm"
      onMouseDown={() => !saving && onClose()}
    >
      <form
        onSubmit={submit}
        className="border-line bg-card w-full max-w-xl rounded-2xl border p-5 shadow-2xl sm:p-6"
        onMouseDown={(event) => event.stopPropagation()}
      >
        <div className="flex items-start justify-between gap-4">
          <div>
            <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">
              Atualização operacional
            </p>
            <h2 id="demand-transition-title" className="font-display mt-1 text-2xl font-semibold">
              {title}
            </h2>
            <p className="text-ink-soft mt-2 text-sm leading-6">{description}</p>
          </div>
          <Button
            type="button"
            variant="ghost"
            size="icon"
            disabled={saving}
            onClick={onClose}
            aria-label="Fechar"
          >
            <XIcon />
          </Button>
        </div>
        <label className="mt-5 block text-sm font-semibold">
          Mensagem
          <textarea
            value={message}
            onChange={(event) => setMessage(event.target.value)}
            placeholder="Explique o que mudou, qual é o próximo passo ou a razão desta atualização."
            className="field-textarea mt-2 min-h-32"
            maxLength={2000}
            autoFocus
          />
        </label>
        <div className="mt-5 flex flex-wrap items-center justify-between gap-3">
          <ImageAttachmentPicker files={images} onChange={setImages} disabled={saving} />
          <Button type="submit" disabled={saving || message.trim().length < 3}>
            {saving ? <Loader2Icon className="animate-spin" /> : <SendIcon />}
            {submitLabel}
          </Button>
        </div>
      </form>
    </div>
  );
}

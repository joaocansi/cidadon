"use client";

import { ImagePlusIcon, Loader2Icon, XIcon } from "lucide-react";
import { type ChangeEvent, useEffect, useMemo } from "react";

import { showToast } from "@/components/shared/toast";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

const acceptedTypes = ["image/jpeg", "image/png", "image/webp"];
const maxFileSize = 2 * 1024 * 1024;

export function ImageAttachmentPicker({
  files,
  onChange,
  disabled,
}: {
  files: File[];
  onChange: (files: File[]) => void;
  disabled?: boolean;
}) {
  const previews = useMemo(() => files.map((file) => URL.createObjectURL(file)), [files]);
  useEffect(() => () => previews.forEach(URL.revokeObjectURL), [previews]);

  function select(event: ChangeEvent<HTMLInputElement>) {
    const selected = Array.from(event.target.files ?? []);
    event.target.value = "";
    if (!selected.length) return;
    if (files.length + selected.length > 5) {
      showToast("Você pode anexar no máximo cinco imagens.", "error");
      return;
    }
    if (selected.some((file) => !acceptedTypes.includes(file.type) || file.size > maxFileSize)) {
      showToast("Use imagens JPG, PNG ou WebP de até 2 MB.", "error");
      return;
    }
    onChange([...files, ...selected]);
  }

  return (
    <div className="space-y-3">
      {files.length ? (
        <div className="flex flex-wrap gap-2">
          {files.map((file, index) => (
            <div key={`${file.name}-${file.lastModified}-${index}`} className="relative">
              {/* Object URLs are transient previews only; uploads use the original File. */}
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img
                src={previews[index]}
                alt={`Prévia do anexo ${index + 1}`}
                className="border-line size-16 rounded-lg border object-cover"
              />
              <Button
                type="button"
                size="icon"
                variant="default"
                className="absolute -top-2 -right-2 size-5 rounded-full"
                disabled={disabled}
                onClick={() => onChange(files.filter((_, fileIndex) => fileIndex !== index))}
                aria-label={`Remover anexo ${index + 1}`}
              >
                <XIcon className="size-3" />
              </Button>
            </div>
          ))}
        </div>
      ) : null}
      <label className="text-ink-soft hover:text-ink inline-flex cursor-pointer items-center gap-2 text-sm font-semibold">
        {disabled ? (
          <Loader2Icon className="size-4 animate-spin" />
        ) : (
          <ImagePlusIcon className="size-4" />
        )}
        Anexar imagens
        <span className="text-ink-faint text-xs font-medium">Opcional</span>
        <Input
          type="file"
          accept="image/jpeg,image/png,image/webp"
          multiple
          className="sr-only"
          disabled={disabled}
          onChange={select}
        />
      </label>
    </div>
  );
}

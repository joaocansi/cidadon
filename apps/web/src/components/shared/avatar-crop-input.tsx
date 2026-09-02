"use client";

import { CameraIcon, Loader2Icon, RotateCcwIcon, XIcon } from "lucide-react";
import { type ChangeEvent, useCallback, useEffect, useMemo, useState } from "react";
import { createPortal } from "react-dom";
import Cropper, { type Area } from "react-easy-crop";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

const maxSourceSize = 10 * 1024 * 1024;

export function AvatarCropInput({
  id,
  value,
  onChange,
  invalid,
  onInvalid,
}: {
  id: string;
  value: File | null;
  onChange: (file: File | null) => void;
  invalid?: boolean;
  onInvalid?: (message: string) => void;
}) {
  const [source, setSource] = useState<string | null>(null);
  const [crop, setCrop] = useState({ x: 0, y: 0 });
  const [zoom, setZoom] = useState(1);
  const [area, setArea] = useState<Area | null>(null);
  const [processing, setProcessing] = useState(false);
  const [editorOpen, setEditorOpen] = useState(false);
  const preview = useMemo(() => (value ? URL.createObjectURL(value) : null), [value]);

  useEffect(
    () => () => {
      if (preview) URL.revokeObjectURL(preview);
    },
    [preview],
  );
  useEffect(
    () => () => {
      if (source) URL.revokeObjectURL(source);
    },
    [source],
  );

  function select(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    event.target.value = "";
    if (!file) return;
    if (
      !["image/jpeg", "image/png", "image/webp"].includes(file.type) ||
      file.size > maxSourceSize
    ) {
      onInvalid?.("Use uma imagem JPG, PNG ou WebP de até 10 MB.");
      return;
    }
    if (source) URL.revokeObjectURL(source);
    setSource(URL.createObjectURL(file));
    setCrop({ x: 0, y: 0 });
    setZoom(1);
    setArea(null);
    setEditorOpen(true);
  }

  async function applyCrop() {
    if (!source || !area) return;
    setProcessing(true);
    try {
      const file = await cropToWebP(source, area);
      onChange(file);
      URL.revokeObjectURL(source);
      setSource(null);
      setEditorOpen(false);
    } finally {
      setProcessing(false);
    }
  }

  const discardCrop = useCallback(() => {
    if (source) URL.revokeObjectURL(source);
    setSource(null);
    setEditorOpen(false);
  }, [source]);

  useEffect(() => {
    if (!editorOpen) return;
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") discardCrop();
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [editorOpen, discardCrop]);

  return (
    <>
      <div className="border-line bg-paper-2 flex flex-wrap items-center gap-3 rounded-xl border p-3 sm:flex-nowrap">
        <div className="border-line bg-paper relative grid size-14 shrink-0 place-items-center overflow-hidden rounded-full border">
          {preview ? (
            // eslint-disable-next-line @next/next/no-img-element
            <img src={preview} alt="Prévia da foto de perfil" className="size-full object-cover" />
          ) : (
            <CameraIcon className="text-ink-faint size-6" />
          )}
        </div>
        <div className="min-w-0 flex-1">
          <p className="text-sm font-semibold">
            {value ? "Foto selecionada" : "Adicione sua foto"}
          </p>
          <p className="text-ink-faint mt-0.5 text-xs">JPG, PNG ou WebP. Até 10 MB.</p>
        </div>
        <label className="text-ink-soft hover:text-ink inline-flex cursor-pointer items-center gap-2 text-sm font-semibold">
          <CameraIcon className="size-4" />
          {value ? "Ajustar foto" : "Selecionar foto"}
          <Input
            id={id}
            type="file"
            accept="image/jpeg,image/png,image/webp"
            className="sr-only"
            aria-invalid={invalid}
            onChange={select}
          />
        </label>
        {value ? (
          <Button type="button" variant="ghost" size="sm" onClick={() => onChange(null)}>
            <XIcon /> Remover
          </Button>
        ) : null}
      </div>
      {source && editorOpen && typeof document !== "undefined"
        ? createPortal(
            <div
              role="dialog"
              aria-modal="true"
              aria-labelledby={`${id}-crop-title`}
              className="bg-ink/50 fixed inset-0 z-[110] grid place-items-center p-4 backdrop-blur-sm"
              onMouseDown={discardCrop}
            >
              <section
                className="bg-card text-ink w-full max-w-xl rounded-2xl p-5 shadow-2xl sm:p-6"
                onMouseDown={(event) => event.stopPropagation()}
              >
                <div className="mb-5 flex items-start justify-between gap-4">
                  <div>
                    <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">
                      Foto de perfil
                    </p>
                    <h2
                      id={`${id}-crop-title`}
                      className="font-display mt-1 text-2xl font-semibold"
                    >
                      Ajuste o enquadramento
                    </h2>
                    <p className="text-ink-soft mt-1 text-sm">
                      Arraste a imagem e use o zoom para posicioná-la no círculo.
                    </p>
                  </div>
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    onClick={discardCrop}
                    aria-label="Cancelar ajuste da foto"
                  >
                    <XIcon />
                  </Button>
                </div>
                <div className="relative h-[min(52vh,24rem)] overflow-hidden rounded-xl bg-black/80">
                  <Cropper
                    image={source}
                    crop={crop}
                    zoom={zoom}
                    aspect={1}
                    cropShape="round"
                    showGrid={false}
                    onCropChange={setCrop}
                    onZoomChange={setZoom}
                    onCropComplete={(_, pixels) => setArea(pixels)}
                  />
                </div>
                <label className="mt-4 flex items-center gap-3 text-xs font-semibold">
                  Zoom
                  <input
                    type="range"
                    min="1"
                    max="3"
                    step="0.05"
                    value={zoom}
                    onChange={(event) => setZoom(Number(event.target.value))}
                    className="accent-lime-deep w-full"
                  />
                </label>
                <div className="mt-5 flex justify-end gap-2">
                  <Button type="button" variant="ghost" size="sm" onClick={discardCrop}>
                    <RotateCcwIcon /> Cancelar
                  </Button>
                  <Button
                    type="button"
                    size="sm"
                    disabled={processing || !area}
                    onClick={() => void applyCrop()}
                  >
                    {processing ? <Loader2Icon className="animate-spin" /> : null}
                    Usar esta foto
                  </Button>
                </div>
              </section>
            </div>,
            document.body,
          )
        : null}
    </>
  );
}

async function cropToWebP(source: string, crop: Area) {
  const image = await loadImage(source);
  const canvas = document.createElement("canvas");
  canvas.width = 512;
  canvas.height = 512;
  const context = canvas.getContext("2d");
  if (!context) throw new Error("Canvas indisponível");
  context.drawImage(image, crop.x, crop.y, crop.width, crop.height, 0, 0, 512, 512);
  const blob = await new Promise<Blob>((resolve, reject) =>
    canvas.toBlob(
      (value) => (value ? resolve(value) : reject(new Error("Não foi possível criar a foto."))),
      "image/webp",
      0.9,
    ),
  );
  return new File([blob], "foto-perfil.webp", { type: "image/webp" });
}

function loadImage(source: string) {
  return new Promise<HTMLImageElement>((resolve, reject) => {
    const image = new Image();
    image.onload = () => resolve(image);
    image.onerror = reject;
    image.src = source;
  });
}

"use client";

import { ChevronLeftIcon, ChevronRightIcon, XIcon } from "lucide-react";
import { useEffect, useState } from "react";

import { cn } from "@/lib/utils";

type DemandImageGalleryProps = {
  images?: string[] | null;
  altPrefix: string;
  className?: string;
  imageClassName?: string;
};

export function DemandImageGallery({
  images,
  altPrefix,
  className,
  imageClassName,
}: DemandImageGalleryProps) {
  const [activeIndex, setActiveIndex] = useState<number>();
  const gallery = images ?? [];
  const activeImage = activeIndex === undefined ? undefined : gallery[activeIndex];

  useEffect(() => {
    if (activeIndex === undefined) return;
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") setActiveIndex(undefined);
      if (event.key === "ArrowLeft")
        setActiveIndex((index) => previousIndex(index, gallery.length));
      if (event.key === "ArrowRight") setActiveIndex((index) => nextIndex(index, gallery.length));
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [activeIndex, gallery.length]);

  if (!gallery.length) return null;
  return (
    <>
      <div className={cn("grid min-w-0 grid-cols-2 gap-2 sm:grid-cols-3", className)}>
        {gallery.map((image, index) => (
          <button
            key={image}
            type="button"
            className="border-line group relative min-w-0 overflow-hidden rounded-lg border text-left focus-visible:outline-2 focus-visible:outline-offset-2"
            onClick={() => setActiveIndex(index)}
            aria-label={`Ampliar ${altPrefix} ${index + 1}`}
          >
            {/* Local data URLs are intentionally rendered without Next image optimization. */}
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={image}
              alt={`${altPrefix} ${index + 1}`}
              className={cn(
                "aspect-video w-full object-cover transition group-hover:scale-[1.03]",
                imageClassName,
              )}
            />
          </button>
        ))}
      </div>
      {activeImage ? (
        <div
          role="dialog"
          aria-modal="true"
          aria-label="Galeria de imagens"
          className="bg-ink/90 fixed inset-0 z-[100] grid place-items-center p-4 backdrop-blur-sm"
          onMouseDown={() => setActiveIndex(undefined)}
        >
          <div
            className="relative flex h-full w-full max-w-6xl items-center justify-center"
            onMouseDown={(event) => event.stopPropagation()}
          >
            <button
              type="button"
              onClick={() => setActiveIndex(undefined)}
              className="bg-card/15 text-card hover:bg-card/25 absolute top-0 right-0 z-10 grid size-10 place-items-center rounded-full"
              aria-label="Fechar galeria"
            >
              <XIcon />
            </button>
            {gallery.length > 1 ? (
              <button
                type="button"
                onClick={() => setActiveIndex((index) => previousIndex(index, gallery.length))}
                className="bg-card/15 text-card hover:bg-card/25 absolute left-0 z-10 grid size-11 place-items-center rounded-full"
                aria-label="Imagem anterior"
              >
                <ChevronLeftIcon />
              </button>
            ) : null}
            {/* Local data URLs are intentionally rendered without Next image optimization. */}
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={activeImage}
              alt={`${altPrefix} ${(activeIndex ?? 0) + 1}`}
              className="max-h-[84vh] max-w-[calc(100%-5rem)] rounded-xl object-contain shadow-2xl"
            />
            {gallery.length > 1 ? (
              <button
                type="button"
                onClick={() => setActiveIndex((index) => nextIndex(index, gallery.length))}
                className="bg-card/15 text-card hover:bg-card/25 absolute right-0 z-10 grid size-11 place-items-center rounded-full"
                aria-label="Próxima imagem"
              >
                <ChevronRightIcon />
              </button>
            ) : null}
            <p className="text-card/80 absolute bottom-1 text-sm">
              {(activeIndex ?? 0) + 1} de {gallery.length}
            </p>
          </div>
        </div>
      ) : null}
    </>
  );
}

function previousIndex(index: number | undefined, length: number) {
  return ((index ?? 0) - 1 + length) % length;
}

function nextIndex(index: number | undefined, length: number) {
  return ((index ?? 0) + 1) % length;
}

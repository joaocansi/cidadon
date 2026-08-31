"use client";

import Image from "next/image";

import { cn } from "@/lib/utils";

export function UserAvatar({
  name,
  imageUrl,
  className,
}: {
  name?: string;
  imageUrl?: string;
  className?: string;
}) {
  const initials = (name || "Usuário")
    .split(/\s+/)
    .slice(0, 2)
    .map((part) => part[0])
    .join("")
    .toUpperCase();

  return (
    <span
      className={cn(
        "bg-lime-pale text-pine ring-pine/10 relative grid size-10 shrink-0 place-items-center overflow-hidden rounded-full text-xs font-bold ring-1",
        className,
      )}
      aria-label={name ? `Foto de ${name}` : "Foto de perfil"}
    >
      {imageUrl ? (
        <Image
          src={imageUrl}
          alt=""
          fill
          sizes="80px"
          className="object-cover"
          unoptimized
          loader={({ src }) => src}
        />
      ) : (
        initials
      )}
    </span>
  );
}

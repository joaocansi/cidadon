"use client";

import {
  ArrowLeftIcon,
  Building2Icon,
  ExternalLinkIcon,
  Globe2Icon,
  MailIcon,
  MapPinIcon,
  PhoneIcon,
} from "lucide-react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useEffect, useState } from "react";

import { CitizenShell } from "@/components/layout/citizen-shell";
import { UserAvatar } from "@/components/shared/user-avatar";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { apiGetOffice, type OfficeContact, type PublicOffice } from "@/lib/api";
import { apiErrorMessage } from "@/lib/forms";

export default function OfficeProfilePage() {
  const params = useParams<{ id: string }>();
  const [office, setOffice] = useState<PublicOffice | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    void apiGetOffice(params.id).then((result) => {
      if (result.ok && result.data) setOffice(result.data);
      else setError(apiErrorMessage(result.error, "Gabinete não encontrado."));
    });
  }, [params.id]);

  const channels = office
    ? [...office.contacts, ...office.social_networks].sort((a, b) => a.position - b.position)
    : [];

  return (
    <CitizenShell
      title={office ? `Gabinete ${office.councillor_name}` : "Perfil do gabinete"}
      subtitle={office ? `Representação parlamentar em ${office.city}/${office.state}.` : undefined}
    >
      <Link
        href="/offices"
        className="text-lime-deep hover:text-pine mb-5 inline-flex items-center gap-2 text-sm font-semibold"
      >
        <ArrowLeftIcon className="size-4" />
        Voltar para gabinetes
      </Link>
      {error ? (
        <Card className="border-destructive/30">
          <CardContent className="text-destructive p-8 text-center text-sm">{error}</CardContent>
        </Card>
      ) : !office ? (
        <div className="text-ink-soft grid min-h-64 place-items-center text-sm">
          Carregando perfil…
        </div>
      ) : (
        <div className="grid gap-5 lg:grid-cols-[minmax(0,1fr)_380px]">
          <Card className="border-line">
            <CardContent className="p-6 sm:p-8">
              <div className="flex flex-col gap-5 sm:flex-row sm:items-center">
                <UserAvatar
                  name={office.councillor_name}
                  imageUrl={office.image_url}
                  className="size-20 text-lg"
                />
                <div>
                  <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">
                    {office.party || "Gabinete parlamentar"}
                  </p>
                  <h2 className="font-display mt-1 text-3xl font-semibold">
                    {office.councillor_name}
                  </h2>
                  <p className="text-ink-soft mt-2 flex items-center gap-1.5 text-sm">
                    <MapPinIcon className="text-lime-deep size-4" />
                    {office.city} · {office.state}
                  </p>
                </div>
              </div>
              <div className="border-line-soft mt-8 grid gap-4 border-t pt-6 sm:grid-cols-2">
                <div className="bg-paper-2/60 rounded-xl p-4">
                  <p className="text-ink-faint text-xs font-semibold tracking-wide uppercase">
                    Sobre
                  </p>
                  <p className="text-ink-soft mt-2 text-sm leading-6">
                    {office.description ||
                      "Este é o perfil público do gabinete para atendimento da população e acompanhamento de sua atuação regional."}
                  </p>
                </div>
                <div className="bg-lime-pale/70 rounded-xl p-4">
                  <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">
                    Transparência
                  </p>
                  <p className="text-pine mt-2 text-sm leading-6">
                    As demandas atribuídas ao gabinete são acompanhadas com status e histórico
                    dentro da plataforma.
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card className="border-line">
            <CardHeader className="border-line-soft border-b">
              <CardTitle className="font-display flex items-center gap-2 text-lg">
                <Building2Icon className="text-lime-deep size-5" />
                Canais de atendimento
              </CardTitle>
            </CardHeader>
            <CardContent>
              {channels.length ? (
                <div className="space-y-2">
                  {channels.map((channel, index) => (
                    <Channel key={`${channel.type}-${channel.value}-${index}`} channel={channel} />
                  ))}
                </div>
              ) : (
                <div className="border-line rounded-xl border border-dashed p-6 text-center">
                  <Globe2Icon className="text-lime-deep mx-auto size-6" />
                  <p className="mt-2 text-sm font-semibold">Nenhum canal público cadastrado</p>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      )}
    </CitizenShell>
  );
}

function Channel({ channel }: { channel: OfficeContact }) {
  const href = channelHref(channel);
  const type = channel.type.toLowerCase();
  const Icon =
    type.includes("email") || type.includes("e-mail")
      ? MailIcon
      : type.includes("tel") || type.includes("whats")
        ? PhoneIcon
        : Globe2Icon;
  const content = (
    <>
      <Icon className="text-lime-deep size-4 shrink-0" />
      <span className="min-w-0">
        <span className="text-ink-faint block text-xs font-semibold tracking-wide uppercase">
          {channel.type}
        </span>
        <span className="text-ink block truncate text-sm font-semibold">{channel.value}</span>
      </span>
      {href ? <ExternalLinkIcon className="text-ink-faint ml-auto size-4 shrink-0" /> : null}
    </>
  );
  return href ? (
    <a
      href={href}
      target={href.startsWith("http") ? "_blank" : undefined}
      rel="noreferrer"
      className="border-line-soft hover:border-lime/70 hover:bg-lime-pale/30 flex items-center gap-3 rounded-xl border p-3 transition"
    >
      {content}
    </a>
  ) : (
    <div className="border-line-soft flex items-center gap-3 rounded-xl border p-3">{content}</div>
  );
}

function channelHref(channel: OfficeContact) {
  const value = channel.value.trim();
  const type = channel.type.toLowerCase();
  if (/^(https?:|mailto:|tel:)/i.test(value)) return value;
  if (type.includes("email") || type.includes("e-mail")) return `mailto:${value}`;
  if (type.includes("whats")) return `https://wa.me/${value.replace(/\D/g, "")}`;
  if (type.includes("tel") || type.includes("celular"))
    return `tel:${value.replace(/[^\d+]/g, "")}`;
  if (type.includes("instagram")) return `https://instagram.com/${value.replace(/^@/, "")}`;
  if (type.includes("facebook")) return `https://facebook.com/${value.replace(/^@/, "")}`;
  if (type.includes("linkedin")) return `https://linkedin.com/in/${value.replace(/^@/, "")}`;
  if (/^[\w.-]+\.[a-z]{2,}/i.test(value)) return `https://${value}`;
  return undefined;
}

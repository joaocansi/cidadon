"use client";

import { Clock3Icon, MailIcon, Trash2Icon, UsersIcon } from "lucide-react";
import type { FormEvent } from "react";

import { UserAvatar } from "@/components/shared/user-avatar";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import type { ManagedOffice, OfficeMemberRequest } from "@/lib/api";

const dateFormatter = new Intl.DateTimeFormat("pt-BR", { dateStyle: "medium", timeStyle: "short" });

export function OfficeTeamPanel({
  members,
  requests,
  email,
  inviting,
  removingId,
  cancellingId,
  error,
  onEmailChange,
  onInvite,
  onRemove,
  onCancel,
}: {
  members: ManagedOffice["members"];
  requests: OfficeMemberRequest[];
  email: string;
  inviting: boolean;
  removingId?: number;
  cancellingId?: number;
  error?: string;
  onEmailChange: (email: string) => void;
  onInvite: (event: FormEvent<HTMLFormElement>) => void;
  onRemove: (member: ManagedOffice["members"][number]) => void;
  onCancel: (request: OfficeMemberRequest) => void;
}) {
  return (
    <div className="space-y-6">
      <section className="space-y-3">
        <div>
          <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">
            Colaboração
          </p>
          <h2 className="font-display mt-1 text-xl font-semibold">Equipe do gabinete</h2>
          <p className="text-ink-soft mt-1 text-sm">
            Convide pessoas, acompanhe os acessos pendentes e revogue acessos quando necessário.
          </p>
        </div>
        <div className="grid gap-4 lg:grid-cols-[minmax(0,1.25fr)_minmax(320px,.75fr)]">
          <Card className="border-line">
            <CardHeader>
              <CardTitle className="font-display flex items-center gap-2">
                <UsersIcon className="text-lime-deep size-5" />
                Membros ativos
              </CardTitle>
            </CardHeader>
            <CardContent>
              {members.length ? (
                <div className="grid gap-2 sm:grid-cols-2">
                  {members.map((member) => (
                    <div
                      key={member.user_id}
                      className="border-line-soft bg-paper/45 flex items-center gap-3 rounded-xl border p-3"
                    >
                      <UserAvatar name={member.name} imageUrl={member.image_url} />
                      <div className="min-w-0 flex-1">
                        <p className="truncate text-sm font-semibold">{member.name}</p>
                        <p className="text-ink-soft truncate text-xs">{member.email}</p>
                      </div>
                      <Button
                        type="button"
                        size="icon"
                        variant="ghost"
                        disabled={removingId === member.user_id}
                        onClick={() => onRemove(member)}
                        aria-label={`Remover ${member.name}`}
                      >
                        <Trash2Icon className="text-destructive size-4" />
                      </Button>
                    </div>
                  ))}
                </div>
              ) : (
                <Empty
                  title="Sua equipe ainda está vazia"
                  text="Convide a primeira pessoa usando o formulário ao lado."
                />
              )}
            </CardContent>
          </Card>
          <Card className="border-line">
            <CardHeader>
              <CardTitle className="font-display flex items-center gap-2">
                <MailIcon className="text-lime-deep size-5" />
                Convidar membro
              </CardTitle>
            </CardHeader>
            <CardContent>
              <form className="space-y-3" onSubmit={onInvite}>
                <div>
                  <Label htmlFor="member-email">E-mail da pessoa</Label>
                  <Input
                    id="member-email"
                    type="email"
                    value={email}
                    onChange={(event) => onEmailChange(event.target.value)}
                    placeholder="membro@gabinete.gov.br"
                    aria-invalid={Boolean(error)}
                    className="bg-paper mt-1.5 h-10"
                  />
                  {error ? (
                    <p className="text-destructive mt-1.5 text-xs font-semibold">{error}</p>
                  ) : null}
                </div>
                <Button type="submit" className="h-10 w-full" disabled={inviting}>
                  <MailIcon />
                  {inviting ? "Enviando…" : "Enviar convite"}
                </Button>
              </form>
              <p className="text-ink-soft mt-3 text-xs leading-5">
                O convite expira em 72 horas. Um novo envio substitui o link anterior para o mesmo
                e-mail.
              </p>
            </CardContent>
          </Card>
        </div>
      </section>
      <section className="space-y-3">
        <div>
          <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">Convites</p>
          <h2 className="font-display mt-1 text-xl font-semibold">Envios pendentes</h2>
        </div>
        <Card className="border-line">
          <CardContent className="p-0">
            {requests.length ? (
              <div className="divide-line-soft divide-y">
                {requests.map((request) => (
                  <div key={request.id} className="flex flex-wrap items-center gap-3 p-4 sm:px-5">
                    <span className="bg-lime-pale text-lime-deep grid size-9 place-items-center rounded-full">
                      <Clock3Icon className="size-4" />
                    </span>
                    <div className="min-w-0 flex-1">
                      <p className="truncate text-sm font-semibold">{request.email}</p>
                      <p className="text-ink-soft mt-0.5 text-xs">
                        Enviado em {dateFormatter.format(new Date(request.created_at))} · expira em{" "}
                        {dateFormatter.format(new Date(request.expires_at))}
                      </p>
                    </div>
                    <Button
                      type="button"
                      size="sm"
                      variant="outline"
                      disabled={cancellingId === request.id}
                      onClick={() => onCancel(request)}
                    >
                      <Trash2Icon className="text-destructive" />
                      {cancellingId === request.id ? "Cancelando…" : "Cancelar"}
                    </Button>
                  </div>
                ))}
              </div>
            ) : (
              <div className="p-6">
                <Empty
                  title="Nenhum convite pendente"
                  text="Os convites em aberto aparecerão aqui até serem aceitos, cancelados ou expirarem."
                />
              </div>
            )}
          </CardContent>
        </Card>
      </section>
    </div>
  );
}

function Empty({ title, text }: { title: string; text: string }) {
  return (
    <div className="border-line rounded-xl border border-dashed p-6 text-center">
      <p className="font-semibold">{title}</p>
      <p className="text-ink-soft mt-1 text-sm">{text}</p>
    </div>
  );
}

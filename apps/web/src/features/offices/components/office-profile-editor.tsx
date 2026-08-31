"use client";

import {
  ArrowDownIcon,
  ArrowUpIcon,
  Building2Icon,
  PlusIcon,
  SaveIcon,
  Trash2Icon,
} from "lucide-react";
import type { FormEvent } from "react";

import { FormField } from "@/components/shared/forms/form-field";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import type { OfficeContact } from "@/lib/api";

export function OfficeProfileEditor({
  description,
  city,
  state,
  contacts,
  socials,
  saving,
  errors,
  onDescriptionChange,
  onCityChange,
  onStateChange,
  onContactsChange,
  onSocialsChange,
  onSave,
}: {
  description: string;
  city: string;
  state: string;
  contacts: OfficeContact[];
  socials: OfficeContact[];
  saving: boolean;
  errors?: { city?: string; state?: string };
  onDescriptionChange: (value: string) => void;
  onCityChange: (value: string) => void;
  onStateChange: (value: string) => void;
  onContactsChange: (items: OfficeContact[]) => void;
  onSocialsChange: (items: OfficeContact[]) => void;
  onSave: (event: FormEvent<HTMLFormElement>) => void;
}) {
  return (
    <section id="perfil" className="scroll-mt-24 space-y-3">
      <div>
        <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">
          Presença pública
        </p>
        <h2 className="font-display mt-1 text-xl font-semibold">Contatos e redes sociais</h2>
      </div>
      <Card className="border-line">
        <CardHeader className="border-line-soft border-b">
          <CardTitle className="font-display flex items-center gap-2">
            <Building2Icon className="text-lime-deep size-5" />
            Perfil público do gabinete
          </CardTitle>
          <p className="text-ink-soft text-sm">
            A ordem abaixo será a mesma exibida para os cidadãos.
          </p>
        </CardHeader>
        <CardContent>
          <form className="grid gap-7 xl:grid-cols-2" onSubmit={onSave}>
            <div className="space-y-4 xl:col-span-2">
              <FormField
                id="office-description"
                label="Sobre o gabinete"
                hint="Explique a atuação, prioridades e como a população pode utilizar este canal."
              >
                <textarea
                  id="office-description"
                  className="field-textarea min-h-32"
                  value={description}
                  maxLength={1500}
                  placeholder="Ex.: Gabinete voltado à mobilidade, zeladoria e escuta ativa dos bairros."
                  onChange={(event) => onDescriptionChange(event.target.value)}
                />
              </FormField>
              <div className="grid gap-4 sm:grid-cols-[minmax(0,1fr)_100px]">
                <FormField id="office-city" label="Cidade de atuação" error={errors?.city}>
                  <Input
                    id="office-city"
                    value={city}
                    placeholder="São Paulo"
                    aria-invalid={Boolean(errors?.city)}
                    className={
                      errors?.city
                        ? "border-destructive focus-visible:ring-destructive/25"
                        : undefined
                    }
                    onChange={(event) => onCityChange(event.target.value)}
                  />
                </FormField>
                <FormField id="office-state" label="UF" error={errors?.state}>
                  <Input
                    id="office-state"
                    value={state}
                    maxLength={2}
                    placeholder="SP"
                    aria-invalid={Boolean(errors?.state)}
                    className={
                      errors?.state
                        ? "border-destructive focus-visible:ring-destructive/25"
                        : undefined
                    }
                    onChange={(event) => onStateChange(event.target.value.toUpperCase())}
                  />
                </FormField>
              </div>
            </div>
            <ChannelEditor title="Contatos" items={contacts} onChange={onContactsChange} />
            <ChannelEditor title="Redes sociais" items={socials} onChange={onSocialsChange} />
            <div className="border-line-soft flex justify-end border-t pt-5 xl:col-span-2">
              <Button type="submit" className="h-10 px-5" disabled={saving}>
                <SaveIcon />
                {saving ? "Salvando…" : "Salvar perfil"}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </section>
  );
}

function ChannelEditor({
  title,
  items,
  onChange,
}: {
  title: string;
  items: OfficeContact[];
  onChange: (items: OfficeContact[]) => void;
}) {
  function move(index: number, delta: number) {
    const target = index + delta;
    if (target < 0 || target >= items.length) return;
    const next = [...items];
    [next[index], next[target]] = [next[target], next[index]];
    onChange(next);
  }

  return (
    <div className="min-w-0">
      <div className="mb-3 flex items-center justify-between gap-3">
        <Label>{title}</Label>
        <Button
          type="button"
          size="sm"
          variant="outline"
          onClick={() => onChange([...items, { type: "", value: "", position: items.length }])}
        >
          <PlusIcon />
          Adicionar
        </Button>
      </div>
      <div className="space-y-2">
        {items.length ? (
          items.map((item, index) => (
            <div
              key={`${title}-${index}`}
              className="border-line-soft bg-paper/45 grid min-w-0 gap-2 rounded-xl border p-2 sm:grid-cols-[120px_minmax(0,1fr)_auto]"
            >
              <Input
                value={item.type}
                placeholder="Tipo"
                aria-label={`Tipo do ${title.toLowerCase()}`}
                onChange={(event) =>
                  onChange(
                    items.map((current, itemIndex) =>
                      itemIndex === index ? { ...current, type: event.target.value } : current,
                    ),
                  )
                }
              />
              <Input
                value={item.value}
                placeholder="Contato, usuário ou link"
                aria-label={`Valor do ${title.toLowerCase()}`}
                onChange={(event) =>
                  onChange(
                    items.map((current, itemIndex) =>
                      itemIndex === index ? { ...current, value: event.target.value } : current,
                    ),
                  )
                }
              />
              <div className="flex justify-end">
                <Button
                  type="button"
                  size="icon"
                  variant="ghost"
                  disabled={index === 0}
                  onClick={() => move(index, -1)}
                  aria-label="Mover para cima"
                >
                  <ArrowUpIcon />
                </Button>
                <Button
                  type="button"
                  size="icon"
                  variant="ghost"
                  disabled={index === items.length - 1}
                  onClick={() => move(index, 1)}
                  aria-label="Mover para baixo"
                >
                  <ArrowDownIcon />
                </Button>
                <Button
                  type="button"
                  size="icon"
                  variant="ghost"
                  onClick={() => onChange(items.filter((_, itemIndex) => itemIndex !== index))}
                  aria-label="Remover canal"
                >
                  <Trash2Icon className="text-destructive" />
                </Button>
              </div>
            </div>
          ))
        ) : (
          <div className="border-line text-ink-soft rounded-xl border border-dashed p-5 text-center text-sm">
            Nenhum canal adicionado.
          </div>
        )}
      </div>
    </div>
  );
}

"use client";

import { Loader2Icon, MapPinIcon, SendIcon } from "lucide-react";
import dynamic from "next/dynamic";
import { useRouter } from "next/navigation";
import { type FormEvent, useEffect, useState } from "react";

import { FormField } from "@/components/shared/forms/form-field";
import { ImageAttachmentPicker } from "@/components/shared/image-attachment-picker";
import { showToast } from "@/components/shared/toast";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { apiCreateDemand, apiListOffices, type CreateDemandInput } from "@/lib/api";
import { apiErrorMessage, type FieldErrors } from "@/lib/forms";

const DemandMap = dynamic(
  () => import("@/features/demands/maps/demand-map").then((module) => module.DemandMap),
  { ssr: false },
);

const categories = [
  "Iluminação pública",
  "Pavimentação",
  "Coleta de lixo",
  "Saúde",
  "Segurança",
  "Transporte",
  "Praças e lazer",
];
const empty: CreateDemandInput = {
  title: "",
  description: "",
  category: categories[0],
  street: "",
  number: "",
  neighborhood: "",
  city: "",
  state: "",
  latitude: 0,
  longitude: 0,
};

export function DemandCreateForm() {
  const router = useRouter();
  const [form, setForm] = useState<CreateDemandInput>(empty);
  const [errors, setErrors] = useState<FieldErrors>({});
  const [offices, setOffices] = useState<
    Array<{ office_id: number; councillor_name: string; party: string }>
  >([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const timeout = window.setTimeout(async () => {
      if (form.city.trim().length < 2 || form.state.length !== 2) {
        setOffices([]);
        return;
      }
      const result = await apiListOffices(form.city.trim(), form.state);
      setOffices(result.ok ? (result.data ?? []) : []);
    }, 350);
    return () => window.clearTimeout(timeout);
  }, [form.city, form.state]);

  function update<K extends keyof CreateDemandInput>(key: K, value: CreateDemandInput[K]) {
    setForm((current) => ({ ...current, [key]: value }));
    setErrors((current) => ({
      ...current,
      [key]: "",
      ...(key === "latitude" || key === "longitude" ? { location: "" } : {}),
    }));
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const next: FieldErrors = {};
    if (form.title.trim().length < 5) next.title = "Use ao menos 5 caracteres.";
    if (form.description.trim().length < 10)
      next.description = "Descreva a situação com ao menos 10 caracteres.";
    if (!form.street.trim()) next.street = "Informe a rua.";
    if (!form.neighborhood.trim()) next.neighborhood = "Informe o bairro.";
    if (form.city.trim().length < 2) next.city = "Informe a cidade.";
    if (form.state.length !== 2) next.state = "Informe a UF com 2 letras.";
    if (form.latitude === 0 && form.longitude === 0) next.location = "Marque o local no mapa.";
    setErrors(next);
    if (Object.keys(next).length) return;

    setLoading(true);
    const result = await apiCreateDemand({
      ...form,
      title: form.title.trim(),
      description: form.description.trim(),
      street: form.street.trim(),
      number: form.number?.trim(),
      neighborhood: form.neighborhood.trim(),
      city: form.city.trim(),
      state: form.state.toUpperCase(),
    });
    setLoading(false);
    if (!result.ok) {
      showToast(apiErrorMessage(result.error, "Não foi possível registrar a demanda."), "error");
      return;
    }
    showToast(`Demanda registrada com sucesso. Protocolo ${result.data?.protocol}.`);
    router.replace("/demandas");
  }

  return (
    <form
      className="grid min-w-0 gap-5 xl:grid-cols-[minmax(360px,.78fr)_minmax(520px,1.22fr)]"
      onSubmit={submit}
    >
      <div className="space-y-5">
        <Card className="border-line">
          <CardHeader className="border-line-soft border-b">
            <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">
              1. O que aconteceu?
            </p>
            <CardTitle className="font-display text-lg">Detalhes da situação</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <FormField id="title" label="Título" error={errors.title}>
              <Input
                id="title"
                value={form.title}
                aria-invalid={Boolean(errors.title)}
                placeholder="Ex.: Poste apagado há duas semanas"
                onChange={(event) => update("title", event.target.value)}
              />
            </FormField>
            <FormField id="category" label="Categoria">
              <select
                id="category"
                className="field-select"
                value={form.category}
                onChange={(event) => update("category", event.target.value)}
              >
                {categories.map((category) => (
                  <option key={category}>{category}</option>
                ))}
              </select>
            </FormField>
            <FormField id="description" label="Descrição" error={errors.description}>
              <textarea
                id="description"
                className="field-textarea min-h-36"
                aria-invalid={Boolean(errors.description)}
                value={form.description}
                placeholder="Explique o problema, há quanto tempo ocorre e o impacto para a região."
                onChange={(event) => update("description", event.target.value)}
              />
            </FormField>
          </CardContent>
        </Card>

        <Card className="border-line">
          <CardHeader className="border-line-soft border-b">
            <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">2. Onde?</p>
            <CardTitle className="font-display text-lg">Endereço e encaminhamento</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-3 sm:grid-cols-[minmax(0,1fr)_100px]">
              <FormField id="street" label="Rua" error={errors.street}>
                <Input
                  id="street"
                  value={form.street}
                  aria-invalid={Boolean(errors.street)}
                  placeholder="Rua das Palmeiras"
                  onChange={(event) => update("street", event.target.value)}
                />
              </FormField>
              <FormField id="number" label="Número">
                <Input
                  id="number"
                  value={form.number}
                  placeholder="245"
                  onChange={(event) => update("number", event.target.value)}
                />
              </FormField>
            </div>
            <FormField id="neighborhood" label="Bairro" error={errors.neighborhood}>
              <Input
                id="neighborhood"
                value={form.neighborhood}
                aria-invalid={Boolean(errors.neighborhood)}
                placeholder="Jardim Bela Vista"
                onChange={(event) => update("neighborhood", event.target.value)}
              />
            </FormField>
            <div className="grid gap-3 sm:grid-cols-[minmax(0,1fr)_90px]">
              <FormField id="city" label="Cidade" error={errors.city}>
                <Input
                  id="city"
                  value={form.city}
                  aria-invalid={Boolean(errors.city)}
                  placeholder="São Paulo"
                  onChange={(event) => update("city", event.target.value)}
                />
              </FormField>
              <FormField id="state" label="UF" error={errors.state}>
                <Input
                  id="state"
                  maxLength={2}
                  value={form.state}
                  aria-invalid={Boolean(errors.state)}
                  placeholder="SP"
                  onChange={(event) => update("state", event.target.value.toUpperCase())}
                />
              </FormField>
            </div>
            <FormField
              id="office"
              label="Direcionar a um gabinete"
              hint="Se você não escolher, a plataforma encaminhará aos gabinetes compatíveis com a região."
            >
              <select
                id="office"
                className="field-select"
                value={form.directed_office_id ?? ""}
                onChange={(event) =>
                  update(
                    "directed_office_id",
                    event.target.value ? Number(event.target.value) : undefined,
                  )
                }
              >
                <option value="">Encaminhamento automático pela região</option>
                {offices.map((office) => (
                  <option key={office.office_id} value={office.office_id}>
                    {office.councillor_name}
                    {office.party ? ` · ${office.party}` : ""}
                  </option>
                ))}
              </select>
            </FormField>
            <FormField
              id="images"
              label="Imagens (opcional)"
              error={errors.images}
              hint="Até 5 imagens JPG, PNG ou WebP, com no máximo 2 MB cada."
            >
              <ImageAttachmentPicker
                files={form.images ?? []}
                disabled={loading}
                onChange={(files) => update("images", files)}
              />
            </FormField>
          </CardContent>
        </Card>
      </div>

      <Card className="border-line min-w-0 self-start xl:sticky xl:top-20">
        <CardHeader className="border-line-soft border-b">
          <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">
            3. Marque o ponto
          </p>
          <CardTitle className="font-display text-lg">Local exato no mapa</CardTitle>
          <p className="text-ink-soft text-sm">
            Clique no mapa onde a situação precisa de atenção.
          </p>
        </CardHeader>
        <CardContent className="space-y-4">
          <DemandMap
            className="border-line h-[min(68vh,680px)] min-h-[520px] w-full overflow-hidden rounded-xl border"
            location={form.latitude || form.longitude ? [form.latitude, form.longitude] : undefined}
            onLocationChange={([latitude, longitude]) => {
              update("latitude", latitude);
              update("longitude", longitude);
            }}
          />
          {errors.location ? (
            <p role="alert" className="text-destructive text-xs font-semibold">
              {errors.location}
            </p>
          ) : form.latitude || form.longitude ? (
            <p className="text-ink-soft flex items-center gap-2 text-xs">
              <MapPinIcon className="text-lime-deep size-4" />
              Ponto selecionado: {form.latitude.toFixed(5)}, {form.longitude.toFixed(5)}
            </p>
          ) : null}
          <Button type="submit" disabled={loading} className="h-11 w-full px-5">
            {loading ? <Loader2Icon className="animate-spin" /> : <SendIcon />}
            {loading ? "Registrando…" : "Registrar demanda"}
          </Button>
        </CardContent>
      </Card>
    </form>
  );
}

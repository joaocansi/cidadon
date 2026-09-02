"use client";

import { Building2Icon, Loader2Icon, MapPinIcon, SearchIcon } from "lucide-react";
import Link from "next/link";
import { type FormEvent, useEffect, useState } from "react";

import { CitizenShell } from "@/components/layout/citizen-shell";
import { FormField } from "@/components/shared/forms/form-field";
import { showToast } from "@/components/shared/toast";
import { UserAvatar } from "@/components/shared/user-avatar";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { apiSearchOffices, type PublicOffice } from "@/lib/api";
import { apiErrorMessage } from "@/lib/forms";

export default function OfficesPage() {
  const [filters, setFilters] = useState({ q: "", city: "", state: "" });
  const [offices, setOffices] = useState<PublicOffice[]>([]);
  const [loading, setLoading] = useState(true);

  async function search(values = filters) {
    setLoading(true);
    const result = await apiSearchOffices({
      q: values.q.trim(),
      city: values.city.trim(),
      state: values.state.trim().toUpperCase(),
    });
    setLoading(false);
    if (!result.ok) {
      showToast(apiErrorMessage(result.error, "Não foi possível buscar os gabinetes."), "error");
      return;
    }
    setOffices(result.data ?? []);
  }

  useEffect(() => {
    void apiSearchOffices().then((result) => {
      setLoading(false);
      if (result.ok) setOffices(result.data ?? []);
      else
        showToast(
          apiErrorMessage(result.error, "Não foi possível carregar os gabinetes."),
          "error",
        );
    });
  }, []);

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    void search();
  }

  return (
    <CitizenShell
      title="Encontre um gabinete"
      subtitle="Pesquise pelo nome do vereador, partido ou região e conheça os canais públicos de atendimento."
    >
      <Card className="border-line mb-6">
        <CardContent className="p-4 sm:p-5">
          <form
            className="grid items-end gap-3 md:grid-cols-[minmax(220px,1.5fr)_1fr_80px_auto]"
            onSubmit={submit}
          >
            <FormField id="office-query" label="Vereador ou partido">
              <Input
                id="office-query"
                value={filters.q}
                onChange={(event) =>
                  setFilters((current) => ({ ...current, q: event.target.value }))
                }
                placeholder="Ex.: Ana Souza ou Partido Verde"
              />
            </FormField>
            <FormField id="office-city" label="Cidade">
              <Input
                id="office-city"
                value={filters.city}
                onChange={(event) =>
                  setFilters((current) => ({ ...current, city: event.target.value }))
                }
                placeholder="São Paulo"
              />
            </FormField>
            <FormField id="office-state" label="UF">
              <Input
                id="office-state"
                value={filters.state}
                onChange={(event) =>
                  setFilters((current) => ({ ...current, state: event.target.value.toUpperCase() }))
                }
                maxLength={2}
                placeholder="SP"
              />
            </FormField>
            <Button type="submit" className="h-10" disabled={loading}>
              <SearchIcon />
              Buscar
            </Button>
          </form>
        </CardContent>
      </Card>

      <div className="mb-3 flex items-center justify-between">
        <p className="text-ink-soft text-sm font-semibold">
          {loading
            ? "Buscando…"
            : `${offices.length} ${offices.length === 1 ? "gabinete encontrado" : "gabinetes encontrados"}`}
        </p>
      </div>

      {loading ? (
        <div className="text-ink-soft grid min-h-64 place-items-center text-sm">
          <span className="flex items-center gap-2">
            <Loader2Icon className="size-4 animate-spin" />
            Buscando gabinetes…
          </span>
        </div>
      ) : offices.length ? (
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {offices.map((office) => (
            <Link key={office.office_id} href={`/gabinetes/${office.slug}`} className="group">
              <Card className="border-line group-hover:border-lime/70 h-full transition duration-200 group-hover:-translate-y-0.5 group-hover:shadow-[var(--shadow-card)]">
                <CardContent className="p-5">
                  <div className="flex items-start gap-3">
                    <UserAvatar
                      name={office.councillor_name}
                      imageUrl={office.image_url}
                      className="size-12"
                    />
                    <div className="min-w-0">
                      <p className="text-lime-deep text-xs font-semibold tracking-wide uppercase">
                        Gabinete parlamentar
                      </p>
                      <h2 className="font-display truncate text-lg font-semibold">
                        {office.councillor_name}
                      </h2>
                      <p className="text-ink-soft text-sm">
                        {office.party || "Partido não informado"}
                      </p>
                    </div>
                  </div>
                  <div className="border-line-soft mt-5 flex items-center justify-between border-t pt-4">
                    <p className="text-ink-soft flex items-center gap-1.5 text-sm">
                      <MapPinIcon className="text-lime-deep size-4" />
                      {office.city} · {office.state}
                    </p>
                    <Building2Icon className="text-ink-faint group-hover:text-lime-deep size-4 transition" />
                  </div>
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      ) : (
        <Card className="border-line border-dashed">
          <CardContent className="p-10 text-center">
            <Building2Icon className="text-lime-deep mx-auto size-8" />
            <p className="mt-3 font-semibold">Nenhum gabinete encontrado</p>
            <p className="text-ink-soft mt-1 text-sm">
              Tente remover algum filtro ou pesquisar por outro nome.
            </p>
          </CardContent>
        </Card>
      )}
    </CitizenShell>
  );
}

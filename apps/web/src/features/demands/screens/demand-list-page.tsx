"use client";

import {
  ClipboardListIcon,
  ListIcon,
  Loader2Icon,
  MapIcon,
  RotateCcwIcon,
  SearchIcon,
  SlidersHorizontalIcon,
} from "lucide-react";
import { type FormEvent, useEffect, useState } from "react";

import { CitizenShell } from "@/components/layout/citizen-shell";
import { FormField } from "@/components/shared/forms/form-field";
import { showToast } from "@/components/shared/toast";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { RoleGate } from "@/features/auth/components/role-gate";
import { DemandCard } from "@/features/demands/components/demand-card";
import { DemandMapExplorer } from "@/features/demands/components/demand-map-explorer";
import { getDemandStatusLabel } from "@/features/demands/components/demand-status-badge";
import {
  apiListDemands,
  apiListMyDemands,
  type Demand,
  type DemandFilters,
  type DemandStatus,
} from "@/lib/api";
import { apiErrorMessage } from "@/lib/forms";

const categories = [
  "Iluminação pública",
  "Pavimentação",
  "Coleta de lixo",
  "Saúde",
  "Segurança",
  "Transporte",
  "Praças e lazer",
];
const statuses: DemandStatus[] = [
  "registered",
  "under_review",
  "in_progress",
  "awaiting_confirmation",
  "completed",
];
const emptyFilters: DemandFilters = { city: "", state: "", neighborhood: "", category: "" };

export default function DemandsPage() {
  const [demands, setDemands] = useState<Demand[]>([]);
  const [draftFilters, setDraftFilters] = useState<DemandFilters>(emptyFilters);
  const [appliedFilters, setAppliedFilters] = useState<DemandFilters>(emptyFilters);
  const [history, setHistory] = useState(false);
  const [regionView, setRegionView] = useState<"map" | "list">("map");
  const [filtersVisible, setFiltersVisible] = useState(false);
  const [selected, setSelected] = useState<Demand>();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let active = true;
    const request = history ? apiListMyDemands() : apiListDemands(appliedFilters);
    void request.then((result) => {
      if (!active) return;
      setLoading(false);
      if (!result.ok) {
        showToast(apiErrorMessage(result.error, "Não foi possível carregar as demandas."), "error");
        return;
      }
      setDemands(result.data ?? []);
    });
    return () => {
      active = false;
    };
  }, [appliedFilters, history]);

  function applyFilters(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setAppliedFilters({
      ...draftFilters,
      city: draftFilters.city?.trim(),
      neighborhood: draftFilters.neighborhood?.trim(),
      state: draftFilters.state?.trim().toUpperCase(),
    });
  }

  function clearFilters() {
    setLoading(true);
    setDraftFilters(emptyFilters);
    setAppliedFilters(emptyFilters);
    setSelected(undefined);
  }

  return (
    <RoleGate allowed={["citizen"]}>
      <CitizenShell
        title={history ? "Meu histórico" : "Demandas da região"}
        subtitle={
          history
            ? "Acompanhe cada protocolo que você registrou."
            : "Consulte situações registradas pela comunidade e seus encaminhamentos."
        }
      >
        <div className="border-line bg-card mb-5 inline-flex rounded-xl border p-1">
          <Button
            variant={history ? "ghost" : "default"}
            onClick={() => {
              setLoading(true);
              setHistory(false);
            }}
          >
            Demandas da região
          </Button>
          <Button
            variant={history ? "default" : "ghost"}
            onClick={() => {
              setLoading(true);
              setHistory(true);
            }}
          >
            Meu histórico
          </Button>
        </div>

        <Card className="border-line">
          <CardHeader className="border-line-soft flex-row items-center justify-between gap-3 border-b">
            <div>
              <p className="text-lime-deep flex items-center gap-2 text-xs font-semibold tracking-wide uppercase">
                <SearchIcon className="size-4" />
                Protocolos
              </p>
              <CardTitle className="font-display mt-1 text-xl">
                {history ? "Demandas que você registrou" : "Acompanhe sua região"}
              </CardTitle>
              {!loading ? (
                <p className="text-ink-soft mt-1 text-sm">
                  {demands.length}{" "}
                  {demands.length === 1 ? "demanda encontrada" : "demandas encontradas"}
                </p>
              ) : null}
            </div>
            {!history ? (
              <div className="flex items-center gap-1">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setFiltersVisible((current) => !current)}
                  aria-expanded={filtersVisible}
                >
                  <SlidersHorizontalIcon />
                  Filtros
                </Button>
                <Button variant="ghost" size="sm" onClick={clearFilters}>
                  <RotateCcwIcon />
                  <span className="hidden sm:inline">Limpar</span>
                </Button>
              </div>
            ) : null}
          </CardHeader>

          {!history ? (
            <>
              <CardContent className="border-line-soft flex flex-wrap items-center justify-between gap-3 border-b px-4 py-3 sm:px-5">
                <div
                  className="border-line bg-paper inline-flex rounded-lg border p-1"
                  role="tablist"
                >
                  <Button
                    type="button"
                    size="sm"
                    variant={regionView === "map" ? "default" : "ghost"}
                    onClick={() => setRegionView("map")}
                    role="tab"
                    aria-selected={regionView === "map"}
                  >
                    <MapIcon />
                    Mapa
                  </Button>
                  <Button
                    type="button"
                    size="sm"
                    variant={regionView === "list" ? "default" : "ghost"}
                    onClick={() => setRegionView("list")}
                    role="tab"
                    aria-selected={regionView === "list"}
                  >
                    <ListIcon />
                    Lista
                  </Button>
                </div>
                <p className="text-ink-soft text-xs">
                  Clique em um ponto para consultar a demanda sem sair do mapa.
                </p>
              </CardContent>
              {filtersVisible ? (
                <CardContent className="border-line-soft bg-paper/40 border-b p-4">
                  <form
                    className="grid items-end gap-3 md:grid-cols-2 xl:grid-cols-[1fr_80px_1fr_1fr_1fr_auto]"
                    onSubmit={applyFilters}
                  >
                    <Filter
                      label="Cidade"
                      value={draftFilters.city ?? ""}
                      placeholder="São Paulo"
                      onChange={(value) =>
                        setDraftFilters((current) => ({ ...current, city: value }))
                      }
                    />
                    <Filter
                      label="UF"
                      value={draftFilters.state ?? ""}
                      placeholder="SP"
                      maxLength={2}
                      onChange={(value) =>
                        setDraftFilters((current) => ({ ...current, state: value.toUpperCase() }))
                      }
                    />
                    <Filter
                      label="Bairro"
                      value={draftFilters.neighborhood ?? ""}
                      placeholder="Jardim Bela Vista"
                      onChange={(value) =>
                        setDraftFilters((current) => ({ ...current, neighborhood: value }))
                      }
                    />
                    <FormField id="category-filter" label="Categoria">
                      <select
                        id="category-filter"
                        className="field-select"
                        value={draftFilters.category ?? ""}
                        onChange={(event) =>
                          setDraftFilters((current) => ({
                            ...current,
                            category: event.target.value,
                          }))
                        }
                      >
                        <option value="">Todas</option>
                        {categories.map((category) => (
                          <option key={category}>{category}</option>
                        ))}
                      </select>
                    </FormField>
                    <FormField id="status-filter" label="Status">
                      <select
                        id="status-filter"
                        className="field-select"
                        value={draftFilters.status ?? ""}
                        onChange={(event) =>
                          setDraftFilters((current) => ({
                            ...current,
                            status: (event.target.value as DemandStatus) || undefined,
                          }))
                        }
                      >
                        <option value="">Todos</option>
                        {statuses.map((status) => (
                          <option key={status} value={status}>
                            {getDemandStatusLabel(status)}
                          </option>
                        ))}
                      </select>
                    </FormField>
                    <Button type="submit" className="h-10">
                      <SearchIcon />
                      Buscar
                    </Button>
                  </form>
                </CardContent>
              ) : null}
            </>
          ) : null}

          <CardContent className="p-4 sm:p-5">
            {loading ? (
              <div className="text-ink-soft grid min-h-64 place-items-center text-sm">
                <span className="flex items-center gap-2">
                  <Loader2Icon className="size-4 animate-spin" />
                  Carregando demandas…
                </span>
              </div>
            ) : demands.length ? (
              !history && regionView === "map" ? (
                <DemandMapExplorer
                  demands={demands}
                  selected={selected}
                  onSelect={setSelected}
                  mapClassName="h-[min(70dvh,680px)] min-h-[500px] w-full rounded-xl"
                />
              ) : (
                <div className="grid gap-4 lg:grid-cols-2">
                  {demands.map((demand) => (
                    <DemandCard key={demand.id} demand={demand} />
                  ))}
                </div>
              )
            ) : (
              <div className="p-8 text-center">
                <ClipboardListIcon className="text-lime-deep mx-auto size-8" />
                <p className="mt-3 font-semibold">Nenhuma demanda encontrada</p>
                <p className="text-ink-soft mt-1 text-sm">
                  {history
                    ? "As demandas que você criar aparecerão aqui."
                    : "Ajuste os filtros para ampliar sua busca."}
                </p>
              </div>
            )}
          </CardContent>
        </Card>
      </CitizenShell>
    </RoleGate>
  );
}

function Filter({
  label,
  value,
  placeholder,
  maxLength,
  onChange,
}: {
  label: string;
  value: string;
  placeholder?: string;
  maxLength?: number;
  onChange: (value: string) => void;
}) {
  const id = `${label.toLowerCase()}-filter`;
  return (
    <FormField id={id} label={label}>
      <Input
        id={id}
        value={value}
        placeholder={placeholder}
        maxLength={maxLength}
        onChange={(event) => onChange(event.target.value)}
      />
    </FormField>
  );
}

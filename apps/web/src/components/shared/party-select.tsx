"use client";

import { CheckIcon, ChevronDownIcon, Loader2Icon, SearchIcon } from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { apiListParties, type PoliticalParty } from "@/lib/api";
import { cn } from "@/lib/utils";

export function PartySelect({
  id,
  value,
  onChange,
  invalid,
  placeholder = "Busque pela sigla ou nome",
}: {
  id: string;
  value: string;
  onChange: (value: string) => void;
  invalid?: boolean;
  placeholder?: string;
}) {
  const [parties, setParties] = useState<PoliticalParty[]>([]);
  const [query, setQuery] = useState("");
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(true);
  const root = useRef<HTMLDivElement>(null);
  const selected = parties.find((party) => party.sigla === value);

  useEffect(() => {
    let active = true;
    void apiListParties().then((result) => {
      if (active && result.ok && result.data) setParties(result.data);
      if (active) setLoading(false);
    });
    return () => {
      active = false;
    };
  }, []);

  useEffect(() => {
    function close(event: MouseEvent) {
      if (!root.current?.contains(event.target as Node)) setOpen(false);
    }
    document.addEventListener("mousedown", close);
    return () => document.removeEventListener("mousedown", close);
  }, []);

  const filtered = useMemo(() => {
    const term = query.trim().toLocaleLowerCase("pt-BR");
    if (!term) return parties;
    return parties.filter(
      (party) =>
        party.sigla.toLocaleLowerCase("pt-BR").includes(term) ||
        party.nome.toLocaleLowerCase("pt-BR").includes(term),
    );
  }, [parties, query]);

  function select(party: PoliticalParty) {
    onChange(party.sigla);
    setQuery("");
    setOpen(false);
  }

  return (
    <div ref={root} className="relative">
      <Button
        id={id}
        type="button"
        variant="outline"
        role="combobox"
        aria-expanded={open}
        aria-controls={`${id}-options`}
        aria-invalid={invalid}
        disabled={loading}
        className={cn(
          "bg-card hover:bg-card h-11 w-full justify-between rounded-2xl px-3 font-normal",
          invalid && "border-destructive ring-destructive/20",
        )}
        onClick={() => setOpen((current) => !current)}
      >
        <span className={cn("truncate", !selected && "text-muted-foreground")}>
          {selected ? `${selected.sigla} · ${selected.nome}` : value || placeholder}
        </span>
        {loading ? (
          <Loader2Icon className="size-4 animate-spin" />
        ) : (
          <ChevronDownIcon className="size-4" />
        )}
      </Button>
      {open ? (
        <div
          id={`${id}-options`}
          role="listbox"
          className="bg-popover border-line absolute z-50 mt-2 w-full overflow-hidden rounded-2xl border p-2 shadow-xl"
        >
          <div className="relative mb-2">
            <SearchIcon className="text-ink-soft pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2" />
            <Input
              autoFocus
              value={query}
              placeholder="Buscar partido"
              className="h-10 pl-9"
              onChange={(event) => setQuery(event.target.value)}
            />
          </div>
          <div className="max-h-56 overflow-y-auto">
            {filtered.length ? (
              filtered.map((party) => (
                <button
                  key={party.id}
                  type="button"
                  role="option"
                  aria-selected={party.sigla === value}
                  className="hover:bg-muted flex w-full items-center justify-between rounded-xl px-3 py-2 text-left text-sm"
                  onClick={() => select(party)}
                >
                  <span className="min-w-0 truncate">
                    <strong>{party.sigla}</strong>{" "}
                    <span className="text-ink-soft">· {party.nome}</span>
                  </span>
                  {party.sigla === value ? (
                    <CheckIcon className="text-lime-deep size-4 shrink-0" />
                  ) : null}
                </button>
              ))
            ) : (
              <p className="text-ink-soft px-3 py-4 text-sm">Nenhum partido encontrado.</p>
            )}
          </div>
        </div>
      ) : null}
    </div>
  );
}

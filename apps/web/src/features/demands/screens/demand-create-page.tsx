"use client";

import { ArrowLeftIcon } from "lucide-react";
import Link from "next/link";

import { CitizenShell } from "@/components/layout/citizen-shell";
import { RoleGate } from "@/features/auth/components/role-gate";
import { DemandCreateForm } from "@/features/demands/components/demand-create-form";

export default function NewDemandPage() {
  return (
    <RoleGate allowed={["citizen"]}>
      <CitizenShell
        title="Nova demanda"
        subtitle="Descreva a situação e marque o ponto exato no mapa para encaminharmos ao gabinete adequado."
      >
        <Link
          href="/demands"
          className="text-lime-deep mb-5 inline-flex items-center gap-2 text-sm font-semibold"
        >
          <ArrowLeftIcon className="size-4" />
          Voltar para demandas
        </Link>
        <DemandCreateForm />
      </CitizenShell>
    </RoleGate>
  );
}

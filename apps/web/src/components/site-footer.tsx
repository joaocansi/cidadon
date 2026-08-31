import Link from "next/link";

import { ProtocolTicker } from "@/components/protocol-ticker";

const PLATFORM_LINKS = [
  { href: "#o-que-e", label: "O que é" },
  { href: "#como-funciona", label: "Como funciona" },
  { href: "#faq", label: "Perguntas frequentes" },
];

const CITIZEN_LINKS = [
  { href: "#cidadaos", label: "Registrar demanda" },
  { href: "#cidadaos", label: "Acompanhar protocolo" },
  { href: "#comecar", label: "Criar conta" },
];

const ALDERMAN_LINKS = [
  { href: "#vereadores", label: "Abrir gabinete" },
  { href: "#vereadores", label: "Adicionar equipe" },
  { href: "#comecar", label: "Cadastrar mandato" },
];

export function SiteFooter() {
  return (
    <footer className="site-footer">
      <div className="wrap">
        <div className="foot-grid">
          <div className="foot-brand">
            <Link
              href="#topo"
              className="font-display text-paper flex items-center gap-2.5 text-[22px] font-bold"
            >
              <span className="logo-mark" />
              cidadon
            </Link>
            <p>
              Uma plataforma para transformar o pedido do bairro em protocolo acompanhável — e a
              resposta do gabinete em algo visível para todo mundo.
            </p>
          </div>
          <div className="foot-col">
            <h5>Plataforma</h5>
            {PLATFORM_LINKS.map((link) => (
              <a key={`platform-${link.href}-${link.label}`} href={link.href}>
                {link.label}
              </a>
            ))}
          </div>
          <div className="foot-col">
            <h5>Cidadão</h5>
            {CITIZEN_LINKS.map((link) => (
              <a key={`citizen-${link.href}-${link.label}`} href={link.href}>
                {link.label}
              </a>
            ))}
          </div>
          <div className="foot-col">
            <h5>Vereador</h5>
            {ALDERMAN_LINKS.map((link) => (
              <a key={`alderman-${link.href}-${link.label}`} href={link.href}>
                {link.label}
              </a>
            ))}
          </div>
        </div>
        <div className="foot-bottom">
          <p>© 2026 Cidadon. Feito para aproximar bairro e gabinete.</p>
          <ProtocolTicker />
        </div>
      </div>
    </footer>
  );
}

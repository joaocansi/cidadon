import {
  BellIcon,
  EyeIcon,
  MapPinIcon,
  MessageSquareIcon,
  ShieldCheckIcon,
  ThumbsUpIcon,
  UsersIcon,
} from "lucide-react";
import Link from "next/link";

import { Reveal } from "@/components/reveal";
import { buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export function Hero() {
  return (
    <section className="hero">
      <div className="wrap hero-grid">
        <div>
          <span className="eyebrow">
            <span className="dot" />
            Feito para a rotina do bairro
          </span>
          <h1>
            Sua rua tem uma demanda?
            <br />
            Agora ela tem <em>um caminho</em> até o gabinete.
          </h1>
          <p className="lead">
            O Cidadon liga quem vive o problema a quem pode resolver. Registre, acompanhe cada etapa
            como um protocolo de verdade e veja seu vereador responder — em público, com data e
            progresso.
          </p>
          <div className="hero-ctas">
            <Link
              href="/cadastro"
              className={cn(
                buttonVariants({ variant: "default" }),
                "h-12 rounded-full px-6 text-[15px]",
              )}
            >
              Sou cidadão, quero registrar
            </Link>
            <Link
              href="/cadastro/vereador"
              className={cn(
                buttonVariants({ variant: "outline" }),
                "h-12 rounded-full px-6 text-[15px]",
              )}
            >
              Sou vereador, quero um gabinete
            </Link>
          </div>
          <p className="hero-note">
            <ShieldCheckIcon className="icon" />
            Gratuito para começar. Sem burocracia para se cadastrar.
          </p>
        </div>

        <Reveal>
          <div className="corkboard">
            <div className="frame" />

            <div className="pin-card ticket">
              <div className="push-pin" />
              <div className="ticket-head">
                <div>
                  <p className="proto">PROTOCOLO Nº 2026-04821</p>
                  <p className="ticket-title">Buraco na calçada</p>
                  <p className="ticket-addr">
                    <MapPinIcon className="icon" />
                    Rua das Palmeiras, 245 · Jardim Bela Vista
                  </p>
                </div>
                <span className="tag-pill andamento">Em andamento</span>
              </div>
              <div className="status-row">
                <div className="status-line">
                  <div className="fill" />
                </div>
                <div className="status-step done">
                  <div className="sdot" />
                  <div className="slabel">Registrada</div>
                </div>
                <div className="status-step done">
                  <div className="sdot" />
                  <div className="slabel">Vista pelo gabinete</div>
                </div>
                <div className="status-step active">
                  <div className="sdot" />
                  <div className="slabel">Em andamento</div>
                </div>
                <div className="status-step">
                  <div className="sdot" />
                  <div className="slabel">Resolvida</div>
                </div>
              </div>
              <div className="ticket-foot">
                <span className="support-count">
                  <ThumbsUpIcon className="icon" />
                  38 vizinhos apoiam
                </span>
                <span className="text-ink-faint text-xs">atualizado há 2 dias</span>
              </div>
            </div>

            <div className="stamp">
              <span>PEDIDO</span>
              <span className="big">VISTO</span>
              <span>PELO GABINETE</span>
            </div>

            <div className="pin-card sticky-note">Notificamos o vereador em segundos ⚡</div>

            <div className="pin-card member-note">
              <p className="mtitle">Gabinete do Vereador</p>
              <div className="avatars">
                <div className="avatar">RM</div>
                <div className="avatar n2">TL</div>
                <div className="avatar n3">JC</div>
                <div className="avatar plus">+2</div>
              </div>
            </div>
          </div>
        </Reveal>
      </div>

      <div className="wrap chip-row">
        <span className="chip">
          <BellIcon className="icon" />
          Notificação automática ao vereador
        </span>
        <span className="chip">
          <EyeIcon className="icon" />
          Progresso público e rastreável
        </span>
        <span className="chip">
          <MessageSquareIcon className="icon" />
          Comentários da comunidade
        </span>
        <span className="chip">
          <UsersIcon className="icon" />
          Gabinete com equipe própria
        </span>
      </div>
    </section>
  );
}

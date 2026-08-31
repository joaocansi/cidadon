import {
  EyeIcon,
  HomeIcon,
  MapPinIcon,
  MessageSquareIcon,
  SendIcon,
  ThumbsUpIcon,
} from "lucide-react";

import { Reveal } from "@/components/reveal";

const CHECKS = [
  {
    icon: SendIcon,
    title: "Registre com sua rua e bairro",
    text: "Conte o que está acontecendo, escolha a categoria e envie direto para quem representa sua região.",
  },
  {
    icon: EyeIcon,
    title: "Acompanhe cada etapa",
    text: "Veja quando o gabinete visualizou, quando entrou em andamento e quando foi resolvida.",
  },
  {
    icon: MessageSquareIcon,
    title: 'Comente e dê um "eu apoio"',
    text: "Mostre que não é só você — vizinhos podem apoiar e comentar na mesma demanda.",
  },
  {
    icon: HomeIcon,
    title: "Veja o que já foi pedido na sua região",
    text: "Antes de registrar de novo, confira se alguém do seu bairro já abriu a mesma demanda.",
  },
];

export function Citizens() {
  return (
    <section id="cidadaos">
      <div className="wrap split">
        <div>
          <span className="section-tag">Para quem vive o bairro</span>
          <h2 className="mt-4 mb-3">Reclamar sozinho cansa. Registrar e acompanhar, não.</h2>
          <p className="text-ink-soft text-base">
            Em poucos minutos, sua demanda sai do “alguém devia fazer algo” e vira um caso com
            endereço, categoria e um responsável do outro lado.
          </p>
          <div className="check-list">
            {CHECKS.map((check) => (
              <div key={check.title} className="check-item">
                <div className="bullet">
                  <check.icon className="icon" />
                </div>
                <div>
                  <h4>{check.title}</h4>
                  <p>{check.text}</p>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="visual-col">
          <Reveal>
            <div className="demand-preview">
              <div className="dp-top">
                <span className="tag-pill">Iluminação pública</span>
                <span className="text-ink-faint text-xs">Protocolo 2026-03110</span>
              </div>
              <h4>Poste apagado há duas semanas</h4>
              <p className="meta">
                <MapPinIcon className="icon" />
                Av. dos Ipês, 88 · Vila Esperança
              </p>
              <div className="comment-block">
                <b>Marina S.:</b> “Mesma coisa aqui na esquina, já são 3 postes.”
              </div>
              <div className="dp-foot">
                <span className="support-btn">
                  <ThumbsUpIcon className="icon" />
                  27 apoios
                </span>
                <span className="text-ink-faint flex items-center gap-1.5 text-xs">
                  <MessageSquareIcon className="icon h-3.5 w-3.5" />9 comentários
                </span>
              </div>
              <div className="demand-list-mini">
                <div className="row">
                  <span>
                    <span className="status-dot" />
                    Coleta atrasada · Rua Rio Verde
                  </span>
                  <span className="text-ink-faint">Resolvida</span>
                </div>
                <div className="row">
                  <span>
                    <span className="status-dot" />
                    Praça sem manutenção · Bairro Alto
                  </span>
                  <span className="text-ink-faint">Em andamento</span>
                </div>
              </div>
            </div>
          </Reveal>
        </div>
      </div>
    </section>
  );
}

import { BellIcon, MapPinIcon, TargetIcon, ThumbsUpIcon } from "lucide-react";

import { Reveal } from "@/components/reveal";

const FEATURES = [
  {
    icon: MapPinIcon,
    title: "Solicitação por região",
    text: "A demanda chega automaticamente até o vereador que representa a rua ou o bairro do cidadão.",
  },
  {
    icon: TargetIcon,
    title: "Rastreamento de verdade",
    text: "Cada demanda tem um status visível: registrada, vista, em andamento ou resolvida — com data de cada mudança.",
  },
  {
    icon: ThumbsUpIcon,
    title: "Comentário e apoio",
    text: "Outros moradores comentam e apoiam uma demanda, mostrando ao gabinete que o problema é coletivo.",
  },
  {
    icon: BellIcon,
    title: "Aviso na hora certa",
    text: "O vereador e o gabinete são notificados assim que uma nova demanda chega para a região deles.",
  },
];

export function Features() {
  return (
    <section id="o-que-e">
      <div className="wrap">
        <div className="section-head">
          <span className="section-tag">O que é o Cidadon</span>
          <h2>Um mural de bairro, só que com protocolo, prazo e resposta.</h2>
          <p>
            Toda demanda registrada vira um caso acompanhável: tem endereço, categoria, histórico de
            atualizações e um responsável do outro lado. Nada se perde num grupo de WhatsApp ou numa
            fila de e-mail.
          </p>
        </div>
        <div className="feature-grid">
          {FEATURES.map((feature, i) => (
            <Reveal key={feature.title} delay={i * 80}>
              <div className="feature-card">
                <div className="feature-icon">
                  <feature.icon className="icon" />
                </div>
                <h3>{feature.title}</h3>
                <p>{feature.text}</p>
              </div>
            </Reveal>
          ))}
        </div>
      </div>
    </section>
  );
}

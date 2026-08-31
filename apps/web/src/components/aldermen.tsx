import { BellIcon, BuildingIcon, CheckCircle2Icon, PlusIcon, UsersIcon } from "lucide-react";

import { Reveal } from "@/components/reveal";

const CHECKS = [
  {
    icon: BuildingIcon,
    title: "Crie seu gabinete oficial",
    text: "Cadastre-se como vereador e abra o gabinete que vai representar sua região dentro do Cidadon.",
  },
  {
    icon: UsersIcon,
    title: "Adicione membros da sua equipe",
    text: "Cada membro pode atuar como moderador, respondendo demandas e atualizando o progresso em seu nome.",
  },
  {
    icon: BellIcon,
    title: "Seja notificado a cada nova demanda",
    text: "Assim que um morador da sua região registra um pedido, seu gabinete recebe o aviso na hora.",
  },
  {
    icon: CheckCircle2Icon,
    title: "Mostre progresso, não promessa",
    text: "Atualize o status, comente o andamento e construa confiança com quem acompanha em público.",
  },
];

export function Aldermen() {
  return (
    <section id="vereadores" className="panel-dark">
      <div className="wrap split">
        <div className="visual-col">
          <Reveal>
            <div className="gabinete-card">
              <div className="gab-head">
                <div className="gab-avatar">GV</div>
                <div>
                  <p className="gname">Gabinete da Vereadora Ana Fontoura</p>
                  <p className="grole">Zona Norte · 4 membros ativos</p>
                </div>
                <span className="bell-flash ml-auto">
                  <BellIcon className="icon" />
                  <span className="ping" />
                </span>
              </div>
              <div className="gab-members">
                <span className="mchip">
                  <span className="ma">AF</span>Titular
                </span>
                <span className="mchip">
                  <span className="ma">RM</span>Moderador
                </span>
                <span className="mchip">
                  <span className="ma">TL</span>Moderador
                </span>
                <span className="mchip add">
                  <PlusIcon className="icon h-3.5 w-3.5" />
                  Adicionar membro
                </span>
              </div>
              <div className="inbox-row">
                <span className="inbox-dot" />
                <div className="itext">
                  <p className="ititle">Buraco na calçada — Jardim Bela Vista</p>
                  <p className="isub">nova demanda · há 4 minutos</p>
                </div>
                <span className="inbox-badge">Nova</span>
              </div>
              <div className="inbox-row">
                <span
                  className="inbox-dot"
                  style={{
                    background: "var(--amber)",
                    boxShadow: "0 0 0 3px rgba(227,162,58,0.25)",
                  }}
                />
                <div className="itext">
                  <p className="ititle">Poste apagado — Vila Esperança</p>
                  <p className="isub">27 apoios · em andamento</p>
                </div>
                <span
                  className="inbox-badge"
                  style={{ background: "rgba(227,162,58,0.18)", color: "var(--amber)" }}
                >
                  Prioridade
                </span>
              </div>
              <div className="inbox-row">
                <span className="inbox-dot" style={{ background: "rgba(245,247,234,0.3)" }} />
                <div className="itext">
                  <p className="ititle">Coleta atrasada — Rua Rio Verde</p>
                  <p className="isub">resolvida há 2 dias</p>
                </div>
              </div>
            </div>
          </Reveal>
        </div>

        <div>
          <span className="section-tag">Para quem representa o bairro</span>
          <h2 className="text-paper mt-4 mb-3">
            Um gabinete de verdade, com equipe, dentro da plataforma.
          </h2>
          <p className="text-paper/70 text-base">
            Ao criar sua conta, você monta seu gabinete digital: convida sua equipe, organiza as
            demandas da sua região e mostra progresso para quem votou em você.
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
      </div>
    </section>
  );
}

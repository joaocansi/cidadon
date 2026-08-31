import { Reveal } from "@/components/reveal";

const STEPS = [
  {
    number: "01",
    marked: true,
    title: "O cidadão registra a demanda",
    text: "Endereço, bairro, categoria e uma descrição do problema — leva menos de dois minutos.",
  },
  {
    number: "02",
    marked: true,
    title: "O gabinete é notificado na hora",
    text: "O vereador e a equipe responsável pela região recebem o aviso assim que a demanda chega.",
  },
  {
    number: "03",
    marked: false,
    title: "A equipe modera e responde",
    text: "Um membro do gabinete assume o caso, comenta e atualiza o status conforme o andamento real.",
  },
  {
    number: "04",
    marked: false,
    title: "A comunidade acompanha e apoia",
    text: 'Vizinhos comentam, dão "eu apoio" e ajudam a mostrar o tamanho real do problema.',
  },
  {
    number: "05",
    marked: true,
    title: "A demanda é resolvida e fica registrada",
    text: "O histórico completo continua público — prova de trabalho para o gabinete, resposta para o bairro.",
  },
];

export function HowItWorks() {
  return (
    <section id="como-funciona">
      <div className="wrap">
        <div className="section-head center">
          <span className="section-tag">Como funciona</span>
          <h2>Da rua até a resposta, em cinco passos</h2>
          <p>
            É uma sequência real — cada demanda passa exatamente por essas etapas, do primeiro
            registro até o histórico público de resolução.
          </p>
        </div>
        <div className="steps">
          {STEPS.map((step) => (
            <Reveal key={step.number}>
              <div className={step.marked ? "step-item marked" : "step-item"}>
                <div className="step-num">{step.number}</div>
                <div className="step-body">
                  <h4>{step.title}</h4>
                  <p>{step.text}</p>
                </div>
              </div>
            </Reveal>
          ))}
        </div>
      </div>
    </section>
  );
}

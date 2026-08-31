import { PlusIcon } from "lucide-react";

import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";

const FAQS = [
  {
    value: "faq-1",
    question: "Preciso pagar para usar o Cidadon?",
    answer:
      "Não. Criar sua conta, registrar demandas e acompanhar o progresso é gratuito para cidadãos. Vereadores também abrem o gabinete sem custo para começar.",
  },
  {
    value: "faq-2",
    question: "Meu vereador ainda não está no Cidadon. E agora?",
    answer:
      "Você pode registrar a demanda normalmente. Assim que o vereador da sua região criar o gabinete, ele encontra o histórico completo da sua área esperando por ele.",
  },
  {
    value: "faq-3",
    question: "Minha demanda fica pública para qualquer pessoa?",
    answer:
      "O andamento e as atualizações de status ficam visíveis para a comunidade — é isso que permite comentar e apoiar. Seus dados de contato pessoal não aparecem publicamente.",
  },
  {
    value: "faq-4",
    question: "Posso apoiar demandas de outros bairros?",
    answer:
      "Sim. Você pode comentar e apoiar qualquer demanda registrada na plataforma, mesmo fora da sua região.",
  },
  {
    value: "faq-5",
    question: "Como sei que minha demanda foi vista?",
    answer:
      'O status muda de "Registrada" para "Vista pelo gabinete" assim que alguém da equipe abre o caso — e você recebe um aviso quando isso acontece.',
  },
];

export function Faq() {
  return (
    <section id="faq" className="bg-paper-2">
      <div className="wrap">
        <div className="section-head center">
          <span className="section-tag">Perguntas frequentes</span>
          <h2>Antes de começar</h2>
        </div>

        <Accordion
          defaultValue={["faq-1"]}
          className="faq-list overflow-visible rounded-none border-0 bg-transparent"
        >
          {FAQS.map((faq) => (
            <AccordionItem
              key={faq.value}
              value={faq.value}
              className="faq-item border-line bg-transparent data-open:bg-transparent"
            >
              <AccordionTrigger className="accordion-trigger faq-summary text-ink rounded-none px-1 py-5.5 text-base font-semibold hover:no-underline [&_[data-slot=accordion-trigger-icon]]:hidden">
                {faq.question}
                <span className="accordion-plus">
                  <PlusIcon />
                </span>
              </AccordionTrigger>
              <AccordionContent className="text-ink-soft px-1 pb-6 text-[15px]">
                {faq.answer}
              </AccordionContent>
            </AccordionItem>
          ))}
        </Accordion>
      </div>
    </section>
  );
}

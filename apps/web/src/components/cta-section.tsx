import Link from "next/link";

import { buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export function CtaSection() {
  return (
    <section id="comecar">
      <div className="wrap">
        <div className="cta-final">
          <span className="eyebrow">
            <span className="dot" />
            Comece agora, sem burocracia
          </span>
          <h2>Sua rua já tem uma demanda esperando para virar protocolo.</h2>
          <p>
            Cadastre-se como cidadão para registrar a sua, ou como vereador para abrir o gabinete da
            sua região.
          </p>
          <div className="hero-ctas">
            <Link
              href="/cadastro"
              className={cn(buttonVariants(), "h-12 rounded-full px-6 text-[15px]")}
            >
              Criar conta de cidadão
            </Link>
            <Link
              href="/cadastro/vereador"
              className={cn(
                buttonVariants({ variant: "outline" }),
                "bg-card h-12 rounded-full px-6 text-[15px]",
              )}
            >
              Abrir gabinete de vereador
            </Link>
          </div>
        </div>
      </div>
    </section>
  );
}

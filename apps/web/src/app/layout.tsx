import "./globals.css";

import type { Metadata } from "next";
import type { ReactNode } from "react";

import { Toaster } from "@/components/shared/toast";
import { SessionProvider } from "@/features/auth/components/session-provider";

export const metadata: Metadata = {
  title: "Cidadon — sua demanda tem um caminho até o gabinete",
  description:
    "Cidadon conecta moradores e vereadores. Registre uma demanda do seu bairro, acompanhe cada etapa e veja o gabinete responder de verdade.",
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="pt-BR" className="h-full antialiased" data-scroll-behavior="smooth">
      <body className="flex min-h-full flex-col">
        <SessionProvider>
          {children}
          <Toaster />
        </SessionProvider>
      </body>
    </html>
  );
}

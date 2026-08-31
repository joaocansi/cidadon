import { Aldermen } from "@/components/aldermen";
import { Citizens } from "@/components/citizens";
import { CtaSection } from "@/components/cta-section";
import { Faq } from "@/components/faq";
import { Features } from "@/components/features";
import { Hero } from "@/components/hero";
import { HowItWorks } from "@/components/how-it-works";
import { SiteFooter } from "@/components/site-footer";
import { SiteHeader } from "@/components/site-header";
import { AuthenticatedRedirect } from "@/features/auth/components/authenticated-redirect";

export default function Home() {
  return (
    <AuthenticatedRedirect>
      <SiteHeader />
      <main id="topo">
        <Hero />
        <Features />
        <Citizens />
        <Aldermen />
        <HowItWorks />
        <Faq />
        <CtaSection />
      </main>
      <SiteFooter />
    </AuthenticatedRedirect>
  );
}

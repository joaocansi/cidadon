import { mkdir } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

import { Document, Page, renderToFile, StyleSheet, Text, View } from "@react-pdf/renderer";
import React from "react";

const h = React.createElement;
const here = path.dirname(fileURLToPath(import.meta.url));
const repositoryRoot = path.resolve(here, "../../..");
const outputPath = path.join(repositoryRoot, "docs", "arquitetura-cidadon.pdf");

const colors = {
  ink: "#16312A",
  soft: "#53645E",
  lime: "#B8DA43",
  limePale: "#EEF6CA",
  paper: "#FBFBF5",
  line: "#D8DED3",
  blue: "#E7F0FA",
  purple: "#EEE9FB",
  white: "#FFFFFF",
};

const styles = StyleSheet.create({
  page: {
    backgroundColor: colors.paper,
    color: colors.ink,
    fontFamily: "Helvetica",
    fontSize: 9.4,
    padding: 42,
  },
  cover: { backgroundColor: colors.ink, color: colors.white, padding: 48 },
  coverMark: {
    color: colors.lime,
    fontSize: 10,
    fontFamily: "Helvetica-Bold",
    letterSpacing: 1.5,
    marginTop: 18,
  },
  coverTitle: {
    fontSize: 33,
    fontFamily: "Helvetica-Bold",
    lineHeight: 1.05,
    marginTop: 18,
    maxWidth: 360,
  },
  coverSubtitle: { color: "#D8E5DF", fontSize: 14, lineHeight: 1.45, marginTop: 18, maxWidth: 350 },
  coverMeta: { color: "#B7C8C0", fontSize: 9, position: "absolute", left: 48, bottom: 48 },
  rule: { backgroundColor: colors.lime, height: 5, width: 70, marginTop: 26 },
  eyebrow: {
    color: "#68820B",
    fontFamily: "Helvetica-Bold",
    fontSize: 8,
    letterSpacing: 1.1,
    textTransform: "uppercase",
    marginBottom: 5,
  },
  title: {
    color: colors.ink,
    fontFamily: "Helvetica-Bold",
    fontSize: 21,
    lineHeight: 1.15,
    marginBottom: 14,
  },
  subtitle: { color: colors.soft, fontSize: 10.5, lineHeight: 1.45, marginBottom: 18 },
  section: { marginBottom: 18 },
  heading: { color: colors.ink, fontFamily: "Helvetica-Bold", fontSize: 13, marginBottom: 7 },
  text: { color: colors.soft, fontSize: 9.2, lineHeight: 1.45 },
  bullet: { color: colors.soft, fontSize: 9.2, lineHeight: 1.45, marginBottom: 3, paddingLeft: 10 },
  footer: {
    bottom: 24,
    color: "#83918B",
    fontSize: 8,
    left: 42,
    position: "absolute",
    right: 42,
    textAlign: "right",
  },
  table: {
    borderColor: colors.line,
    borderWidth: 1,
    borderRadius: 5,
    overflow: "hidden",
    marginBottom: 12,
  },
  tableHeader: { backgroundColor: colors.ink, color: colors.white, flexDirection: "row" },
  tableHeaderCell: { fontFamily: "Helvetica-Bold", fontSize: 8.4, padding: 6 },
  tableRow: { borderColor: colors.line, borderTopWidth: 1, flexDirection: "row" },
  tableLabel: {
    backgroundColor: "#F2F4EC",
    borderColor: colors.line,
    borderRightWidth: 1,
    color: colors.ink,
    fontFamily: "Helvetica-Bold",
    fontSize: 8.3,
    padding: 6,
    width: "25%",
  },
  tableValue: { color: colors.soft, fontSize: 8.3, lineHeight: 1.35, padding: 6, width: "75%" },
  twoColumn: { flexDirection: "row", gap: 12 },
  card: {
    backgroundColor: colors.white,
    borderColor: colors.line,
    borderRadius: 7,
    borderWidth: 1,
    flexGrow: 1,
    padding: 10,
  },
  cardTitle: { color: colors.ink, fontFamily: "Helvetica-Bold", fontSize: 10, marginBottom: 5 },
  diagram: {
    backgroundColor: colors.white,
    borderColor: colors.line,
    borderRadius: 8,
    borderWidth: 1,
    marginBottom: 16,
    padding: 14,
  },
  diagramRow: {
    alignItems: "center",
    flexDirection: "row",
    justifyContent: "center",
    marginBottom: 10,
  },
  diagramBox: {
    alignItems: "center",
    borderRadius: 6,
    borderWidth: 1,
    minHeight: 42,
    justifyContent: "center",
    paddingHorizontal: 8,
    width: 112,
  },
  boxClient: { backgroundColor: colors.blue, borderColor: "#9DC1E3" },
  boxWeb: { backgroundColor: colors.purple, borderColor: "#B9ABEB" },
  boxAPI: { backgroundColor: colors.limePale, borderColor: "#B7D15E" },
  boxInfra: { backgroundColor: "#F3F4EE", borderColor: colors.line },
  boxText: { color: colors.ink, fontFamily: "Helvetica-Bold", fontSize: 8.6, textAlign: "center" },
  arrow: { color: "#81928A", fontSize: 17, marginHorizontal: 8 },
  diagramCaption: { color: colors.soft, fontSize: 8.3, lineHeight: 1.35, textAlign: "center" },
  useCase: { marginBottom: 13 },
  useCaseTitle: {
    backgroundColor: colors.limePale,
    borderColor: "#CFE18D",
    borderRadius: 5,
    borderWidth: 1,
    color: colors.ink,
    fontFamily: "Helvetica-Bold",
    fontSize: 10.5,
    padding: 7,
  },
  note: {
    backgroundColor: "#F4F7EA",
    borderColor: "#D7E8A4",
    borderLeftWidth: 3,
    borderRadius: 4,
    color: colors.soft,
    fontSize: 8.7,
    lineHeight: 1.4,
    marginBottom: 12,
    padding: 8,
  },
});

function PageFrame({ children, number }) {
  return h(
    Page,
    { size: "A4", style: styles.page },
    children,
    h(Text, { fixed: true, style: styles.footer }, `Cidadon · Arquitetura do Projeto · ${number}`),
  );
}

function Section({ eyebrow, title, children }) {
  return h(
    View,
    { style: styles.section },
    eyebrow && h(Text, { style: styles.eyebrow }, eyebrow),
    h(Text, { style: styles.heading }, title),
    children,
  );
}

function KeyValueTable({ rows }) {
  return h(
    View,
    { style: styles.table, wrap: false },
    rows.map(([label, value]) =>
      h(
        View,
        { key: label, style: styles.tableRow },
        h(Text, { style: styles.tableLabel }, label),
        h(Text, { style: styles.tableValue }, value),
      ),
    ),
  );
}

function MatrixTable({ headers, rows, widths }) {
  return h(
    View,
    { style: styles.table },
    h(
      View,
      { style: styles.tableHeader },
      headers.map((header, index) =>
        h(
          Text,
          { key: header, style: { ...styles.tableHeaderCell, width: widths[index] } },
          header,
        ),
      ),
    ),
    rows.map((row, rowIndex) =>
      h(
        View,
        { key: `row-${rowIndex}`, style: styles.tableRow },
        row.map((cell, index) =>
          h(
            Text,
            {
              key: `${rowIndex}-${index}`,
              style: {
                ...styles.tableValue,
                width: widths[index],
                ...(index < row.length - 1
                  ? { borderColor: colors.line, borderRightWidth: 1 }
                  : {}),
              },
            },
            cell,
          ),
        ),
      ),
    ),
  );
}

function ArchitectureDiagram() {
  const Box = ({ name, tone }) =>
    h(View, { style: [styles.diagramBox, tone] }, h(Text, { style: styles.boxText }, name));
  return h(
    View,
    { style: styles.diagram },
    h(
      View,
      { style: styles.diagramRow },
      h(Box, { name: "Navegador", tone: styles.boxClient }),
      h(Text, { style: styles.arrow }, "→"),
      h(Box, { name: "Next.js\napps/web", tone: styles.boxWeb }),
      h(Text, { style: styles.arrow }, "→"),
      h(Box, { name: "API Go\ncmd/api", tone: styles.boxAPI }),
    ),
    h(
      View,
      { style: styles.diagramRow },
      h(Box, { name: "Worker Go\ncmd/worker", tone: styles.boxAPI }),
      h(Text, { style: styles.arrow }, "→"),
      h(Box, { name: "Casos de uso\napplication", tone: styles.boxAPI }),
      h(Text, { style: styles.arrow }, "→"),
      h(Box, { name: "PostgreSQL\n& Mailpit", tone: styles.boxInfra }),
    ),
    h(
      Text,
      { style: styles.diagramCaption },
      "A API e o worker compartilham os mesmos casos de uso. A API expõe HTTP, SSE, OpenAPI YAML e a referência Scalar.",
    ),
  );
}

const useCases = [
  [
    "UC-01 · Cadastro e acesso do cidadão",
    [
      ["Atores", "Cidadão, API."],
      ["Objetivo", "Criar a conta e acessar a área destinada ao cidadão."],
      ["Pré-requisitos", "E-mail ainda não cadastrado; nome, senha, cidade e UF válidos."],
      [
        "Fluxo principal",
        "1. O cidadão envia seus dados. 2. A API valida e cria o perfil. 3. Efetua login. 4. A API cria cookies HTTP-only. 5. O frontend redireciona para a área cidadã.",
      ],
      [
        "Exceções",
        "E-mail duplicado, dados inválidos ou credenciais incorretas retornam códigos estáveis para mensagens específicas no frontend.",
      ],
      ["Pós-condições", "Conta e sessão ativas; o usuário pode registrar e acompanhar demandas."],
    ],
  ],
  [
    "UC-02 · Cadastro de vereador e configuração do gabinete",
    [
      ["Atores", "Vereador, API."],
      ["Objetivo", "Criar o perfil parlamentar e tornar público o gabinete."],
      ["Pré-requisitos", "E-mail livre; nome, partido, foto, cidade e UF informados."],
      [
        "Fluxo principal",
        "1. A API cria usuário e vereador. 2. O vereador autentica-se. 3. Cria ou edita o gabinete. 4. Define descrição, região, contatos e redes sociais.",
      ],
      [
        "Exceções",
        "Usuários sem papel de vereador não podem editar configurações exclusivas do gabinete.",
      ],
      ["Pós-condições", "O gabinete pode ser pesquisado pelo nome do vereador, partido ou região."],
    ],
  ],
  [
    "UC-03 · Convite e ingresso de membro do gabinete",
    [
      ["Atores", "Vereador, convidado, SMTP/Mailpit, API."],
      ["Objetivo", "Incluir uma pessoa na equipe autenticada do gabinete."],
      ["Pré-requisitos", "Vereador autenticado, gabinete existente e e-mail do convidado válido."],
      [
        "Fluxo principal",
        "1. API cria token expirável. 2. Envia e-mail HTML identificando o gabinete. 3. Convidado abre o link. 4. Informa nome, foto e senha. 5. API valida token e cria o membro.",
      ],
      [
        "Exceções",
        "Convite expirado, cancelado ou já utilizado é rejeitado; e-mail existente retorna conflito.",
      ],
      ["Pós-condições", "Membro acessa as demandas e notificações do gabinete."],
    ],
  ],
  [
    "UC-04 · Criar e direcionar uma demanda",
    [
      ["Atores", "Cidadão, API, gabinetes candidatos."],
      ["Objetivo", "Registrar uma necessidade regional e encaminhá-la ao gabinete adequado."],
      [
        "Pré-requisitos",
        "Cidadão autenticado; detalhes, endereço, categoria e coordenadas no mapa válidos.",
      ],
      [
        "Fluxo principal",
        "1. Cidadão informa dados e imagens opcionais. 2. Escolhe gabinete ou deixa o direcionamento automático. 3. API persiste a demanda e cria evento. 4. Algoritmo classifica gabinetes por atuação e raio progressivo. 5. Destinatários recebem notificação.",
      ],
      [
        "Exceções",
        "Sem gabinete no raio inicial, a busca amplia o alcance e recorre à cidade/UF; imagens inválidas são recusadas.",
      ],
      ["Pós-condições", "Demanda fica registrada e visível nos mapas/listas autorizados."],
    ],
  ],
  [
    "UC-05 · Assumir e executar uma demanda",
    [
      ["Atores", "Vereador ou membro do gabinete; cidadão autor."],
      ["Objetivo", "Definir um responsável único e registrar o atendimento."],
      ["Pré-requisitos", "Demanda atribuída ao gabinete e em estado permitido."],
      [
        "Fluxo principal",
        "1. Primeiro gabinete assume e torna-se responsável. 2. Estado vai para em análise. 3. Gabinete inicia execução. 4. Publica atualização. 5. Solicita validação e inicia prazo de 120 horas.",
      ],
      [
        "Exceções",
        "Outro gabinete não pode atuar após a assunção; transições fora da ordem são rejeitadas.",
      ],
      ["Pós-condições", "Timeline registra as ações e envolvidos recebem notificação."],
    ],
  ],
  [
    "UC-06 · Confirmar, reabrir ou concluir automaticamente",
    [
      ["Atores", "Cidadão autor, gabinete responsável, worker."],
      ["Objetivo", "Concluir a demanda após validação do cidadão ou após o prazo."],
      [
        "Pré-requisitos",
        "Demanda em aguardando validação; confirmação/reabertura somente pelo autor.",
      ],
      [
        "Fluxo principal",
        "1. Autor confirma e a demanda é concluída; ou 2. reabre com conteúdo e ela retorna para em análise; ou 3. o worker conclui ao final de 120 horas.",
      ],
      ["Exceções", "Usuário não autor, estado inválido ou reabertura sem conteúdo são recusados."],
      ["Pós-condições", "Evento e notificação registram a decisão ou a expiração."],
    ],
  ],
  [
    "UC-07 · Conversar, responder e moderar",
    [
      ["Atores", "Cidadão autenticado, vereador, membro do gabinete, autor do comentário."],
      ["Objetivo", "Manter conversa pública organizada, com anexos e moderação."],
      ["Pré-requisitos", "Autenticação; texto e/ou até cinco imagens válidas."],
      [
        "Fluxo principal",
        "1. Usuário publica comentário raiz. 2. Resposta recebe parent_id e é exibida abaixo do pai. 3. Gabinete aparece com selo e avatar. 4. Imagens abrem galeria contextual. 5. Autor pode excluir seu comentário e respostas diretas.",
      ],
      [
        "Exceções",
        "Não autor não pode excluir. Somente gabinete responsável pode ocultar comentário denunciado.",
      ],
      [
        "Pós-condições",
        "Exclusão ou moderação gera evento auditável; conteúdo removido deixa de aparecer na conversa.",
      ],
    ],
  ],
  [
    "UC-08 · Receber notificações em tempo real",
    [
      ["Atores", "Cidadão, vereador, membro do gabinete, API."],
      ["Objetivo", "Informar participantes sem depender de recarga da página."],
      ["Pré-requisitos", "Usuário autenticado e envolvido diretamente na demanda."],
      [
        "Fluxo principal",
        "1. Caso de uso persiste notificação, excluindo o autor da ação. 2. Frontend recupera não lidas. 3. Mantém SSE em /notifications/stream. 4. Evento novo atualiza o sino e exibe toast.",
      ],
      [
        "Exceções",
        "Se SSE cair, a central persistente recupera as notificações na próxima consulta.",
      ],
      ["Pós-condições", "Usuário marca notificações como lidas individualmente ou em lote."],
    ],
  ],
];

function ArchitectureDocument() {
  return h(
    Document,
    { title: "Arquitetura do Cidadon", author: "Cidadon" },
    h(
      Page,
      { size: "A4", style: styles.cover },
      h(Text, { style: styles.coverMark }, "CIDADON · DOCUMENTAÇÃO TÉCNICA"),
      h(Text, { style: styles.coverTitle }, "Arquitetura e cenários de uso"),
      h(View, { style: styles.rule }),
      h(
        Text,
        { style: styles.coverSubtitle },
        "Visão estrutural do monorepo, componentes, responsabilidades e fluxos de negócio da plataforma de cidadania.",
      ),
      h(
        Text,
        { style: styles.coverMeta },
        "Versão 1.0 · Gerado automaticamente a partir do repositório",
      ),
    ),
    h(
      PageFrame,
      { number: "01" },
      h(Text, { style: styles.eyebrow }, "VISÃO DE SOLUÇÃO"),
      h(Text, { style: styles.title }, "Arquitetura do sistema"),
      h(
        Text,
        { style: styles.subtitle },
        "O Cidadon organiza interface, API e processamento assíncrono em processos independentes que compartilham regras de negócio e persistência.",
      ),
      h(ArchitectureDiagram),
      h(
        Section,
        { title: "Responsabilidades do monorepo" },
        h(MatrixTable, {
          headers: ["Área", "Responsabilidade", "Tecnologias"],
          widths: ["25%", "48%", "27%"],
          rows: [
            [
              "apps/web",
              "Interface para cidadãos, vereadores e membros; rotas, componentes por feature, Leaflet e SSE.",
              "Next.js, React, TypeScript, shadcn",
            ],
            [
              "apps/backend",
              "API HTTP, regras de negócio, worker, autenticação e contrato OpenAPI.",
              "Go, Gin, Gorm",
            ],
            [
              "infra",
              "Serviços locais necessários para desenvolvimento e teste de e-mail.",
              "PostgreSQL, Mailpit, Docker Compose",
            ],
            [
              "docs",
              "Documentação funcional, arquitetura e artefatos de entrega.",
              "Markdown, PDF",
            ],
          ],
        }),
      ),
      h(
        Section,
        { title: "Interfaces públicas" },
        h(MatrixTable, {
          headers: ["URL", "Finalidade"],
          widths: ["38%", "62%"],
          rows: [
            ["http://localhost:3000", "Aplicação web."],
            ["http://localhost:8080", "API HTTP e stream de notificações."],
            ["http://localhost:8080/openapi.yaml", "Contrato OpenAPI legível por ferramentas."],
            ["http://localhost:8080/docs", "Referência interativa moderna com Scalar."],
            ["http://localhost:8025", "Caixa de e-mail local do Mailpit."],
          ],
        }),
      ),
    ),
    h(
      PageFrame,
      { number: "02" },
      h(Text, { style: styles.eyebrow }, "ORGANIZAÇÃO INTERNA"),
      h(Text, { style: styles.title }, "Camadas, dados e segurança"),
      h(
        Section,
        { title: "Camadas do backend" },
        h(MatrixTable, {
          headers: ["Camada", "Função", "Exemplos"],
          widths: ["20%", "48%", "32%"],
          rows: [
            [
              "Domain",
              "Entidades e portas; não conhece HTTP, banco ou SMTP.",
              "Demand, Office, User, contratos de repositório",
            ],
            [
              "Application",
              "Orquestra os casos de uso e valida regras de autorização/transição.",
              "Assumir demanda, comentar, convite",
            ],
            [
              "Adapters",
              "Traduz HTTP, PostgreSQL e serviços externos para as portas.",
              "Handlers Gin, Gorm, SMTP, JWT",
            ],
            [
              "Platform",
              "Configuração, conexão e infraestrutura transversal.",
              ".env, transações, banco",
            ],
            ["Bootstrap", "Compõe dependências, migrações e rotas.", "cmd/api e cmd/worker"],
          ],
        }),
      ),
      h(
        Section,
        { title: "Dados e regras centrais" },
        h(MatrixTable, {
          headers: ["Recurso", "Regra"],
          widths: ["29%", "71%"],
          rows: [
            [
              "Demanda",
              "Registrada → em análise → em execução → aguardando validação → concluída.",
            ],
            ["Responsabilidade", "O primeiro gabinete a assumir torna-se o responsável único."],
            [
              "Evento",
              "Registro imutável de transições, moderação, exclusão e conclusão automática.",
            ],
            [
              "Comentário",
              "Texto e/ou imagens, um nível de resposta visual, autor pode excluir comentário e respostas diretas.",
            ],
            ["Notificação", "Persistente, direcionada apenas a envolvidos e entregue por SSE."],
          ],
        }),
      ),
      h(
        Section,
        { title: "Autorização" },
        h(
          View,
          { style: styles.twoColumn },
          h(
            View,
            { style: styles.card },
            h(Text, { style: styles.cardTitle }, "Cidadão"),
            h(
              Text,
              { style: styles.text },
              "Registra demandas, acompanha região, comenta e confirma ou reabre a própria demanda.",
            ),
          ),
          h(
            View,
            { style: styles.card },
            h(Text, { style: styles.cardTitle }, "Vereador"),
            h(
              Text,
              { style: styles.text },
              "Configura gabinete e equipe; atua nas demandas atribuídas.",
            ),
          ),
          h(
            View,
            { style: styles.card },
            h(Text, { style: styles.cardTitle }, "Membro"),
            h(
              Text,
              { style: styles.text },
              "Atua nas demandas do gabinete, respeitando permissões exclusivas do vereador.",
            ),
          ),
        ),
      ),
      h(
        Text,
        { style: styles.note },
        "A autenticação usa cookies HTTP-only com JWT. O handler valida o papel e o caso de uso valida autoria, gabinete responsável e estado antes de alterar dados.",
      ),
    ),
    ...useCases.map(([title, rows], index) =>
      h(
        PageFrame,
        { key: title, number: String(index + 3).padStart(2, "0") },
        h(Text, { style: styles.eyebrow }, "CENÁRIO DE USO"),
        h(Text, { style: styles.title }, title),
        h(KeyValueTable, { rows }),
        index === 6 &&
          h(
            Text,
            { style: styles.note },
            "A conversa e a linha do tempo são deliberadamente separadas: comentários não são eventos de negócio. Eventos permanecem como registro auditável.",
          ),
        index === 7 &&
          h(
            Section,
            { title: "Validação e entrega" },
            h(
              Text,
              { style: styles.text },
              "O comando make ci executa formatação Go, go vet, testes Go, ESLint, Prettier, typecheck e validação do OpenAPI. A pipeline GitHub Actions executa essas validações e o build de produção do frontend.",
            ),
          ),
      ),
    ),
  );
}

await mkdir(path.dirname(outputPath), { recursive: true });
await renderToFile(h(ArchitectureDocument), outputPath);
console.log(`PDF gerado em ${outputPath}`);

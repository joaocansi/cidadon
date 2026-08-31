# Cidadon

Monorepo da plataforma de cidadania que conecta demandas da região a gabinetes.

```text
apps/
  backend/     API Go e worker, com Clean Architecture pragmática
  web/         Aplicação Next.js organizada por feature
infra/
  compose/     Serviços locais de desenvolvimento
  postgres/    Documentação do banco local
  mailpit/     Documentação do e-mail local
docs/          Requisitos de produto
```

## Desenvolvimento

```bash
make up
make down
make help
```

Endereços locais: web em http://localhost:3000, API em http://localhost:8080 e Mailpit em http://localhost:8025.

O contrato HTTP está em [apps/backend/docs/openapi.yaml](apps/backend/docs/openapi.yaml).

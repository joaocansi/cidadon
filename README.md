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

## Imagens

Por padrão, a API grava imagens em `.runtime/uploads` e as publica em `http://localhost:8080/media`. Para usar S3, configure `MEDIA_DRIVER=s3`, `MEDIA_PUBLIC_BASE_URL` (a URL pública do bucket) e as variáveis `MEDIA_S3_BUCKET`, `MEDIA_S3_REGION`, `MEDIA_S3_ACCESS_KEY_ID` e `MEDIA_S3_SECRET_ACCESS_KEY` em `apps/backend/.env`. `MEDIA_S3_ENDPOINT` é opcional para provedores S3 compatíveis. O bucket precisa permitir leitura pública dos objetos.

O worker migra automaticamente data URLs legadas para o storage configurado, em lotes e de forma idempotente.

O contrato HTTP está em [apps/backend/docs/openapi.yaml](apps/backend/docs/openapi.yaml).

A visão técnica e os cenários de uso estão em [docs/arquitetura.md](docs/arquitetura.md) e no [PDF diagramado](docs/arquitetura-cidadon.pdf).

# Cidadon Backend

Dois executáveis compartilham o núcleo de regras:

- `cmd/api`: API HTTP Gin.
- `cmd/worker`: conclusão automática de demandas em validação vencida.

As camadas estão em `internal/domain`, `internal/application`, `internal/adapters`, `internal/platform` e `internal/bootstrap`. Use `make api`, `make worker` e `make test` na raiz.

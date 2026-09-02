SHELL := /bin/bash

BACKEND_DIR := apps/backend
WEB_DIR := apps/web
COMPOSE_FILE := infra/compose/docker-compose.yml
RUNTIME_DIR := .runtime
GO_CACHE ?= /tmp/cidadon-go-cache

.PHONY: help install up down infra-up infra-down db mailpit api worker web test vet lint typecheck format format-check check backend-check openapi ci build fmt

help:
	@echo "Cidadon monorepo"
	@echo "  make up                       infraestrutura + API + worker + web"
	@echo "  make api | worker | web        inicia um processo específico"
	@echo "  make db | mailpit              inicia apenas infraestrutura"
	@echo "  make down                      encerra processos e containers"
	@echo "  make test | vet | lint | typecheck | format | check | ci | build | fmt"

install:
	cd $(WEB_DIR) && pnpm install --frozen-lockfile

infra-up:
	docker compose -f $(COMPOSE_FILE) up -d

infra-down:
	docker compose -f $(COMPOSE_FILE) down

db:
	docker compose -f $(COMPOSE_FILE) up -d postgres

mailpit:
	docker compose -f $(COMPOSE_FILE) up -d mailpit

api:
	@mkdir -p $(RUNTIME_DIR)
	@if test -f $(RUNTIME_DIR)/api.pid && kill -0 $$(cat $(RUNTIME_DIR)/api.pid) 2>/dev/null; then echo "API já está em execução"; else cd $(BACKEND_DIR) && GOCACHE=$(GO_CACHE) go build -o ../../$(RUNTIME_DIR)/cidadon-api ./cmd/api && (nohup ../../$(RUNTIME_DIR)/cidadon-api >../../$(RUNTIME_DIR)/api.log 2>&1 & echo $$! >../../$(RUNTIME_DIR)/api.pid); fi

worker:
	@mkdir -p $(RUNTIME_DIR)
	@if test -f $(RUNTIME_DIR)/worker.pid && kill -0 $$(cat $(RUNTIME_DIR)/worker.pid) 2>/dev/null; then echo "Worker já está em execução"; else cd $(BACKEND_DIR) && GOCACHE=$(GO_CACHE) go build -o ../../$(RUNTIME_DIR)/cidadon-worker ./cmd/worker && (nohup ../../$(RUNTIME_DIR)/cidadon-worker >../../$(RUNTIME_DIR)/worker.log 2>&1 & echo $$! >../../$(RUNTIME_DIR)/worker.pid); fi

web:
	@mkdir -p $(RUNTIME_DIR)
	@if test -f $(RUNTIME_DIR)/web.pid && kill -0 $$(cat $(RUNTIME_DIR)/web.pid) 2>/dev/null; then echo "Web já está em execução"; else (cd $(WEB_DIR); nohup pnpm dev >../../$(RUNTIME_DIR)/web.log 2>&1 & echo $$! >../../$(RUNTIME_DIR)/web.pid); fi

up: infra-up api worker web

down:
	@test ! -f $(RUNTIME_DIR)/api.pid || kill $$(cat $(RUNTIME_DIR)/api.pid) 2>/dev/null || true
	@test ! -f $(RUNTIME_DIR)/worker.pid || kill $$(cat $(RUNTIME_DIR)/worker.pid) 2>/dev/null || true
	@test ! -f $(RUNTIME_DIR)/web.pid || kill $$(cat $(RUNTIME_DIR)/web.pid) 2>/dev/null || true
	docker compose -f $(COMPOSE_FILE) down

test:
	cd $(BACKEND_DIR) && GOCACHE=$(GO_CACHE) go test ./...

vet:
	cd $(BACKEND_DIR) && GOCACHE=$(GO_CACHE) go vet ./...

lint:
	cd $(WEB_DIR) && pnpm lint

typecheck:
	cd $(WEB_DIR) && pnpm exec tsc --noEmit

format:
	cd $(WEB_DIR) && pnpm format

format-check:
	cd $(WEB_DIR) && pnpm format:check

check:
	cd $(WEB_DIR) && pnpm check

backend-check:
	@test -z "$$(gofmt -l $$(find $(BACKEND_DIR) -name '*.go' -type f))"
	$(MAKE) vet test

openapi:
	cd $(WEB_DIR) && pnpm --package=@redocly/cli@1.27.0 dlx redocly lint ../backend/docs/openapi.yaml

ci: backend-check check openapi

build: test
	cd $(WEB_DIR) && pnpm build

fmt:
	gofmt -w $(BACKEND_DIR)

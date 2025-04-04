# --- CONFIG ---

LOCAL_BIN     := $(CURDIR)/bin
MIGRATION_DIR := $(CURDIR)/migrations

BUF_VERSION           := v1.51.0
GOLANGCI_VERSION      := v1.64.5
GOOSE_VERSION         := v3.24.2

RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

define install_tool
	GOBIN=$(LOCAL_BIN) go install $(1)@$(2)
endef

# --- INSTALL TOOLS ---

.PHONY: install
install:
	mkdir -p $(LOCAL_BIN)
	go mod tidy
	$(call install_tool,github.com/bufbuild/buf/cmd/buf,$(BUF_VERSION))
	$(call install_tool,github.com/pressly/goose/v3/cmd/goose,$(GOOSE_VERSION))
	$(call install_tool,github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCI_VERSION))

# --- LINT ---

.PHONY: lint
lint:
	$(LOCAL_BIN)/golangci-lint run --fix

.PHONY: lint-proto
lint-proto:
	$(LOCAL_BIN)/buf lint

# --- PROTO GEN ---

.PHONY: gen-proto
gen-proto:
	$(LOCAL_BIN)/buf generate

# --- MIGRATIONS ---

.PHONY: migrate_status
migrate_status:
	$(LOCAL_BIN)/goose -dir $(MIGRATION_DIR) postgres "$(POSTGRES_DSN)" status

.PHONY: migrate_up
migrate_up:
	$(LOCAL_BIN)/goose -dir $(MIGRATION_DIR) postgres "$(POSTGRES_DSN)" up

.PHONY: migrate_down
migrate_down:
	$(LOCAL_BIN)/goose -dir $(MIGRATION_DIR) postgres "$(POSTGRES_DSN)" down

.PHONY: migrate_new
migrate_new:
ifndef name
	$(error migration name not specified. Use: make migrate_new name=your_migration)
endif
	$(LOCAL_BIN)/goose -dir $(MIGRATION_DIR) create $(name) sql

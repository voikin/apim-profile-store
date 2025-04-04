# --- CONFIG ---

LOCAL_BIN     := $(CURDIR)/bin
MIGRATION_DIR := $(CURDIR)/migrations

GRPC_GATEWAY_VERSION  := v2.25.1
GEN_GO_VERSION        := v1.31.0
GEN_GO_GRPC_VERSION   := v1.5.1
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
	# go mod tidy
	$(call install_tool,github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway,$(GRPC_GATEWAY_VERSION))
	$(call install_tool,github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2,$(GRPC_GATEWAY_VERSION))
	$(call install_tool,google.golang.org/protobuf/cmd/protoc-gen-go,$(GEN_GO_VERSION))
	$(call install_tool,google.golang.org/grpc/cmd/protoc-gen-go-grpc,$(GEN_GO_GRPC_VERSION))
	$(call install_tool,github.com/bufbuild/buf/cmd/buf,$(BUF_VERSION))
	$(call install_tool,github.com/pressly/goose/v3/cmd/goose,$(GOOSE_VERSION))
	$(call install_tool,github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCI_VERSION))

# --- LINT ---

.PHONY: lint
lint:
	$(LOCAL_BIN)/golangci-lint run --fix

.PHONY: lint-proto
lint-proto: update-buf
	PATH="$(PATH):$(LOCAL_BIN)" $(LOCAL_BIN)/buf lint

# --- PROTO GEN ---

.PHONY: update-buf
update-buf:
	$(LOCAL_BIN)/buf dep update

.PHONY: gen-proto
gen-proto: update-buf
	PATH="$(PATH):$(LOCAL_BIN)" $(LOCAL_BIN)/buf generate

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

# --- CONFIG ---

LOCAL_BIN     := $(CURDIR)/bin
MIGRATION_DIR := $(CURDIR)/migrations

GOLANGCI_VERSION := v1.64.5
GOOSE_VERSION    := v3.24.2
MINIMOCK_VERSION := v3.4.5

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
	$(call install_tool,github.com/pressly/goose/v3/cmd/goose,$(GOOSE_VERSION))
	$(call install_tool,github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCI_VERSION))
	$(call install_tool,github.com/gojuno/minimock/v3/cmd/minimock,$(MINIMOCK_VERSION))

# --- LINT ---

.PHONY: lint
lint:
	$(LOCAL_BIN)/golangci-lint run --fix

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

.PHONY: mock
mock:
	$(LOCAL_BIN)/minimock -i github.com/voikin/apim-profile-store/internal/usecase.PostgresRepo \
	-o ./internal/usecase/mocks \
	-s _mock.go
	$(LOCAL_BIN)/minimock -i github.com/voikin/apim-profile-store/internal/usecase.Neo4jRepo \
	-o ./internal/usecase/mocks \
	-s _mock.go
	$(LOCAL_BIN)/minimock -i github.com/voikin/apim-profile-store/internal/usecase.TrManager \
	-o ./internal/usecase/mocks \
	-s _mock.go
	
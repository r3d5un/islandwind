# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DATABASE
# ==================================================================================== #

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database="${ISLANDWIND_DB_CONNSTR}" up

## db/migrations/goto number=$1: target version to migrate to
.PHONY: db/migrations/goto
db/migrations/goto: confirm
	@echo 'Running down migrations...'
	migrate -path=./migrations -database="${ISLANDWIND_DB_CONNSTR}" goto ${number}

## db/migrations/down
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Running down migrations...'
	migrate -path=./migrations -database="${ISLANDWIND_DB_CONNSTR}" down

# ==================================================================================== #
# FORMATTING
# ==================================================================================== #
.PHONY: format/backend
format/backend:
	@echo 'Formatting code...'
	go fmt ./...
	golines . -w

# ==================================================================================== #
# TESTING
# ==================================================================================== #

.PHONY: test/backend/staticcheck
test/backend/staticCheck:
	@echo 'Performing static analysis'
	staticcheck ./...

.PHONY: test/backend/vet
test/backend/vet:
	@echo 'Vetting code...'
	go vet ./...

.PHONY: test/backend
test/backend: test/backend/staticCheck test/backend/vet
	@echo 'Running tests...'
	go test -race -vet=off ./...

.PHONY: test/backend/reload
test/backend/reload:
	@echo 'Running tests...'
	find . -name "*.go" | entr -c go test -race -vet=off ./...

# ==================================================================================== #
# TESTING
# ==================================================================================== #

.PHONY: build/backend/linux
build/backend/linux: test/backend
	mkdir -p ./build/bin
	GOOS=linux GOARCH=amd64 go build -o ./build/bin/api ./cmd/api/

.PHONY: build/backend/docker
build/backend/docker: test/backend
	docker build -t r3d5un/islandwind/backend -f ./build/api.Dockerfile .

# ==================================================================================== #
# RUNNERS
# ==================================================================================== #
.PHONY: run/backend
run/backend: format/backend test/backend
	@echo 'Running backend...'
	go run ./cmd/api/

.PHONY: run/backend/reload
run/backend/reload: format/backend test/backend
	@echo 'Running backend with live reload...'
	find . -name "*.go" | entr -c go run ./cmd/api/

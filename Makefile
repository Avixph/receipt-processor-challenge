#include .env

BINARY_NAME = api
BINARY_DIRECTORY = /tmp
BINARY_PATH = ${BINARY_DIRECTORY}/bin/${BINARY_NAME}
MAIN_PATH = ./cmd/api
MIGRATION_PATH = ./migrations

AIR = github.com/air-verse/air@latest
MIGRATE = github.com/golang-migrate/migrate/v4/cmd/migrate@latest
STATICCHECK = honnef.co/go/tools/cmd/staticcheck@latest
GOVULNCHECK = golang.org/x/vuln/cmd/govulncheck@latest

# ================================================================================== #
# HELPERS
# ================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [Y/N] ' && read ans && [ $${ans:-N} = Y ]

# ================================================================================== #
# Build
# ================================================================================== #

## build: build the cmd/api application
.PHONY: build
build:
	@echo 'Building cmd/api...'
	cd server && go build -o .${BINARY_PATH} ${MAIN_PATH}

## build/clean: Clean the cmd/api application executables
.PHONY: build/clean
build/clean:
	@echo 'Cleaning cmd/api executables...'
	cd server && rm -r -f .${BINARY_DIRECTORY} && go clean -cache

# ================================================================================== #
# DEVELOPMENT
# ================================================================================== #

## run: run the cmd/api application
.PHONY: run
run:
	@echo 'Running cmd/api...'
	cd server && go run ${MAIN_PATH}
#-db-dsn=${DATABASE_DSN}

## watch: run the cmd/api application with live reloading on file changes
.PHONY: watch
watch:
	cd server && go run ${AIR} \
		--build.cmd "cd .. make build" --build.bin "/tmp/bin/api" --build.delay "100" \
		--build.exclude_dir "assets, bin, internal, migrations, remote, testdata, vendor" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true" \
		--screen.clear_on_rebuild "true" \


# ================================================================================== #
# Docker
# ================================================================================== #

## docker/run: run the cmd/api application in a docker container
.PHONY: docker/run
docker/run:
	@echo 'Running cmd/api docker container...'
	docker compose up

## docker/watch: run the cmd/api application in a docker container with live reloading
.PHONY: docker/watch
docker/watch:
	@echo 'Running up Docker live...'
	docker compose watch

## docker/stop: stop running the cmd/api application in a docker container
.PHONY: docker/stop
docker/stop: confirm
	@echo 'Running down Docker...'
	docker compose down

# ================================================================================== #
# Quality Control
# ================================================================================== #

## tidy: tidy module dependencies, and format all .go files
.PHONY: tidy
tidy:
	@echo 'Tidying module dependencies...'
	cd server && go mod tidy
	@echo 'Formatting .go files...'
	cd server && go fmt ./...

## audit: run quality control checks
.PHONY: audit
audit: test
	@echo 'Checking module dependencies'
	cd server && go mod tidy -diff
	cd server && go mod verify
	@echo 'Vetting code...'
	cd server && go vet ./...
	cd server && go run ${STATICCHECK} -checks=all,-ST1000,-U1000 ./...
	cd server && go run ${GOVULNCHECK} ./...

## test: run all tests
.PHONY: test
test:
	@echo 'Running tests...'
	cd server && go test -v -race -buildvcs ./...
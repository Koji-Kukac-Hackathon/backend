LOCAL_UID=$(shell id -u)
LOCAL_GID=$(shell id -g)
OUTPUT_BINARY=bin/zgrabi-mjesto
DOCKER_COMPOSE=docker-compose --project-name 'zgrabi-mjesto'
DOCKER_COMPOSE_DEV=$(DOCKER_COMPOSE) \
	-f 'docker-compose.yml' \
	-f 'docker-compose.dev.yml' \
	-f 'docker-compose.override.yml'

PACKAGE=zgrabi-mjesto.hr
define LDFLAGS
-X '$(PACKAGE)/app/version.buildTimestamp=$(shell date -u '+%Y-%m-%dT%H:%M:%S%z')'
endef
LDFLAGS:=$(strip $(LDFLAGS))

.PHONY: build
build:
	CGO_ENABLED=0 \
	go \
	build \
	-a \
	-tags osusergo,netgo \
	-gcflags "all=-N -l" \
	-ldflags="-s -w -extldflags \"-static\" $(LDFLAGS)" \
	-o "${OUTPUT_BINARY}" \
	main.go

.PHONY: clean
clean:
	rm -rf $(OUTPUT_BINARY)

.PHONY: run
run: build
	./${OUTPUT_BINARY}

.PHONY: format
format:
	gofmt -e -l -s -w .

.PHONY: fmt
fmt: format

.PHONY: compact
compact: build
	upx --brute "${OUTPUT_BINARY}"

.PHONY: sync-deps
sync-deps:
	CGO_ENABLED=0 go mod download

.PHONY: $pull
$pull:
	git pull --rebase

.PHONY: dev/run
dev/run: dev/build
	./${OUTPUT_BINARY}

.PHONY: dev/build
dev/build:
	go \
	build \
	-tags osusergo,netgo \
	-ldflags="-s -w -extldflags \"-static\" $(LDFLAGS)" \
	-o "${OUTPUT_BINARY}" \
	main.go

# REQUIRED: go install github.com/cosmtrek/air@latest
.PHONY: dev/watch
dev/watch:
	air

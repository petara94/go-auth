BIN=./bin
EXPORTS = env PATH="$(PWD)/bin:$(PATH)" GOBIN="$(PWD)/bin"

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

CMD=go-auth
OUT=$(BIN)/$(GOOS)-$(GOARCH)/$(CMD)

GIT_BRANCH := $$(git rev-parse --abbrev-ref HEAD)
GIT_TAG := $(shell git describe --abbrev=0 --tags)
GIT_REV := $(shell git rev-parse --short HEAD)

prepare:
	GOBIN="$$PWD/bin" go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2
.PHONY: prepare

config:
	cp example.config.yml config.yml
.PHONY: config

vendor:
	go mod tidy
	go mod vendor
.PHONY: vendor

migrations.create:
	$(BIN)/migrate create -ext sql -dir migrations -seq $(M_NAME)
.PHONY: migrations.create

build:
	rm -rf $(OUT)
	$(EXPORTS) go build -ldflags "-X github.com/petara94/go-auth.Version=$(GIT_TAG)-$(GIT_REV)" -o $(OUT) ./cmd/$(CMD)



run: build
	./$(OUT)
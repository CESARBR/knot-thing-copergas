VENV := venv
PRE_COMMIT_HOOKS := .git/hooks/pre-commit

BIN := $(VENV)/bin
PIP := $(BIN)/pip
PRE_COMMIT := $(BIN)/pre-commit
PRE_COMMIT_VERSION := 2.11.1

GOCMD=go
GOSECCMD=gosec
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get -u -v

GOPATH := $(shell go env GOPATH)
OS := $(shell uname -s | awk '{print tolower($$0)}')
BINARY = app
GOARCH = amd64

LDFLAGS = -ldflags="$$(govvv -flags)"


.PHONY: bootstrap
bootstrap: venv \
	pre-commit-hooks \

.PHONY: clean-all
clean-all: clean \
	clean-bootstrap \
	clean-tools

.PHONY: venv
venv:
	python3 -m venv $(VENV)

.PHONY: pre-commit-hooks
pre-commit-hooks:
	$(PIP) install pre-commit==$(PRE_COMMIT_VERSION)
	$(PRE_COMMIT) install

.PHONY: tools
tools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/kisielk/errcheck@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/axw/gocov/gocov@latest
	go install github.com/matm/gocov-html@latest
	go install github.com/ahmetb/govvv@latest
	go install github.com/cespare/reflex@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: run
run: bin
	./$(BINARY)-$(OS)-$(GOARCH)

.PHONY: watch
watch:
	reflex -s -r '\.go$$' go run cmd/main.go

.PHONY: bin
bin:
	env CGO_ENABLED=0 GOOS=$(OS) GOARCH=${GOARCH} go build -a -installsuffix cgo ${LDFLAGS} -o ${BINARY}-$(OS)-${GOARCH} cmd/main.go ;

.PHONY: http-docs
http-docs:
	swag init -g pkg/server/server.go

.PHONY: sectest
sectest:
	$(GOSECCMD) -exclude-dir=docs -fmt=json ./...

.PHONY: lint
lint:
	golangci-lint run $(go list ./... | grep -v /vendor/) --timeout 10m

.PHONY: cover
cover:
	${GOCMD} test -coverprofile=coverage.out ./... && ${GOCMD} tool cover -html=coverage.out -o coverage.html

.SILENT: clean
.PHONY: clean
clean:
	$(GOCLEAN)
	@rm -f ${BINARY}-$(OS)-${GOARCH}
	@rm -f coverage.out
	@rm -f coverage.html
	@rm -f $(PLANTUML)

.PHONY: clean-bootstrap
clean-bootstrap:
	@echo "clean: undoing bootstrap..."
	@rm -rf $(VENV)
	@rm -f $(PRE_COMMIT_HOOKS)

.PHONY: clean-tools
clean-tools:
	@echo "clean: removing tools..."
	@rm -f $(GOPATH)/bin/goimports
	@rm -f $(GOPATH)/bin/errcheck
	@rm -f $(GOPATH)/bin/golangci-lint
	@rm -f $(GOPATH)/bin/gocov
	@rm -f $(GOPATH)/bin/gocov-html
	@rm -f $(GOPATH)/bin/govvv
	@rm -f $(GOPATH)/bin/reflex
	@rm -f $(GOPATH)/bin/gosec
	@rm -f $(GOPATH)/bin/swag

.DEFAULT_GOAL       := help
VERSION             := v0.0.0
TARGET_MAX_CHAR_NUM := 20

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)


GCP_ZONE ?= europe-west3-a
GCP_REGION ?= europe-west3
KEYRING_NAME ?= authz_keyring
KEY_NAME ?= authz_key
KMS_REGION ?= global
MAX_NODES ?= 15

MINIKUBE_VMDRIVER :=
ifeq ($(OS),Windows_NT)
    MINIKUBE_VMDRIVER += virtualbox
else
    OSNAME := $(shell uname -s)
    ifeq ($(OSNAME), Linux)
        MINIKUBE_VMDRIVER += kvm2
    else ifeq ($(OSNAME), Darwin)
        MINIKUBE_VMDRIVER += hyperkit
    endif
endif

.PHONY: help build prepare flu-web-run flu-mob-run clean

## Show help
help:
	@echo ''
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

proto-init:
	@type -p protoc > /dev/null 2>&1 || echo 'Please install protobuf first according to your package manager'
	@go get -u google.golang.org/protobuf/cmd/protoc-gen-go

modules-install:
	@go mod tidy
	@go mod vendor

## Compiles protobuf
proto-go: proto-init modules-install
	@protoc -I vendor/ -I api/ api/*.proto --go_out=pkg/api/



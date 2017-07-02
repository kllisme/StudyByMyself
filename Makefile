GO = godep go
ENTRY = main.go
RUNTIME_CONFIG_FILE = ./runtime.config
BINARY_NAME = ./binary

VERSION = $(shell git describe --tags)

env ?= development

# config
config:
	@NODE_ENV=${env} node -e "var cfg=require('config');cfg.version='${VERSION}';console.log(JSON.stringify(cfg, null, 2))" > ${RUNTIME_CONFIG_FILE}.json


# compile
binary: config
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${GO} build -o ${BINARY_NAME}

package: binary


# server
development: config
	@VERSION=${VERSION} ${GO} run ${ENTRY} -conf ${RUNTIME_CONFIG_FILE}
server:
	@VERSION=${VERSION} ${BINARY_NAME} -conf ${RUNTIME_CONFIG_FILE}


# alias
dev: development


.PHONY: config binary package development server dev

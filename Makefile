# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
all: help
.PHONY: help initial bank test
.PHONY: run down
.PHONY: spanner-up spanner-down spanner-init
.PHONY: spanner-execute spanner-migration-up spanner-migration-down spanner-migration-version spanner-migration-goto spanner-migration-force

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
TAG ?= :v$(shell date +%Y%m%d)
TAG_STAGE = :stage
IMG_PREFIX = pharmacy
IMG_IMPORTER = $(IMG_PREFIX)-importer
IMG_RESTFUL = $(IMG_PREFIX)-restful
IMG_TEST = $(IMG_PREFIX)-test
CTR_SPAN = spanner-emulator
CMD_TOOL = docker run --rm --network 'container:$(CTR_SPAN)' justdomepaul/tool:v1.0.0 dockerize $(1) $(2)

wire: ## generate wire service file
	docker run -ti --rm -v ${PWD}:/mnt justdomepaul/wire -c \
"go mod download && go mod tidy && \
wire ./cmd/importer && \
wire ./cmd/restful"

run: build spanner-up spanner-init spanner-migration-up-default import restful## run system

down:
	docker-compose down

build: ## build image
	docker build -f ./cmd/importer/Dockerfile --tag $(IMG_IMPORTER)$(TAG) --rm .
	docker image tag $(IMG_IMPORTER)$(TAG) $(IMG_IMPORTER)$(TAG_STAGE)
	docker build -f ./cmd/restful/Dockerfile --tag $(IMG_RESTFUL)$(TAG) --rm .
	docker image tag $(IMG_RESTFUL)$(TAG) $(IMG_RESTFUL)$(TAG_STAGE)

import: ## import data
	docker-compose up importer

restful: ## up restful api server
	docker-compose up -d restful

test: spanner-up spanner-init ## integration testing
	docker build -t $(IMG_TEST) -f test.Dockerfile .
	docker run --rm --network host $(IMG_TEST) go test -count=1 -cover ./...
	$(MAKE) spanner-down

spanner-generate: ## generate spanner schema file
	./spanner-migrate.sh generate

spanner-up: ## up spanner
	docker-compose up -d spanner-emulator

spanner-down: ## down spanner
	docker-compose kill spanner-emulator
	docker rm spanner-emulator

spanner-init: ## init spanner
	$(eval CURL = curl -s 'localhost:9020/v1/projects/test-project/instances')
	$(eval CURL += --data '{"instanceId": "test-instance"}')
	$(eval WAIT = -wait http://localhost:9020/v1/projects/foo/instances)
	$(call CMD_TOOL, $(WAIT), $(CURL))
	$(eval CURL = curl -s 'localhost:9020/v1/projects/test-project/instances/test-instance/databases')
	$(eval CURL += --data '{"createStatement": "CREATE DATABASE `test-database`"}')
	$(eval WAIT = -wait http://localhost:9020/v1/projects/foo/instances)
	$(call CMD_TOOL, $(WAIT), $(CURL))

spanner-execute: ## execute spanner
	./spanner.sh

spanner-migration-up-default: ## migrate spanner up default
	./spanner-migrate.sh upDefault

spanner-migration-up: ## migrate spanner up
	./spanner-migrate.sh up

spanner-migration-down: ## migrate spanner down
	./spanner-migrate.sh down

spanner-migration-version: ## migrate spanner version
	./spanner-migrate.sh version

spanner-migration-goto: ## migrate spanner goto
	./spanner-migrate.sh goto

spanner-migration-force: ## migrate spanner force
	./spanner-migrate.sh force

postgresql-generate: ## generate spanner schema file
	./postgresql-migrate.sh generate

postgresql-up: ## up spanner
	docker-compose up -d postgresql

postgresql-down: ## down spanner
	docker-compose kill postgresql
	docker rm postgresql

postgresql-migration-up: ## migrate spanner up
	./postgresql-migrate.sh up

postgresql-migration-down: ## migrate spanner down
	./postgresql-migrate.sh down

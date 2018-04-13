.DEFAULT_GOAL := help

## Setup
devel-deps:
	go get -u -v github.com/golang/dep/cmd/dep

## Build container test image
build-container-test-image:
	docker image build -t youyo/zabbix-aws-integration-agent:test .

## Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: help
.SILENT:

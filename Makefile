.DEFAULT_GOAL := help

## Setup
download-libs:
	go mod download

## Build container test image
build-container-test-image:
	docker image build -t youyo/zaia:test .

## Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: help
.SILENT:

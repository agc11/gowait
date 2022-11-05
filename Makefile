SHELL:=/bin/bash

coverage:
	@go test ./... -cover -race -coverprofile=coverage.out

coverage-html:
	@go tool cover -html=coverage.out

coverage-badge:
	./scripts/generate-badge.sh


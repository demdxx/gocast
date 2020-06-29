
.PHONY: fmt
fmt: ## Run formatting code
	@echo "Fix formatting"
	@gofmt -w ${GO_FMT_FLAGS} $$(go list -f "{{ .Dir }}" ./...); if [ "$${errors}" != "" ]; then echo "$${errors}"; fi

.PHONY: test
test: ## Run tests
	go test -race ./...

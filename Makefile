include dev.env
export

.PHONY: fmt
fmt:
	gofmt -l -w .

.PHONY: lint
lint:
	go vet ./...
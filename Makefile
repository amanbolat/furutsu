include dev.env
export

.PHONY: fmt
fmt:
	goimports -l -w .

.PHONY: lint
lint:
	go vet
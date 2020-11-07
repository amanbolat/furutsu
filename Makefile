include dev.env
export

.PHONY: fmt
fmt:
	gofmt -l -w .

.PHONY: lint
lint:
	go vet ./...

.PHONY: test
test.integration:
	docker-compose rm -s -f
	docker-compose --file docker-compose.local_test.yml up -d --force-recreate
	go test -count=1 -tags=integration ./integration_tests

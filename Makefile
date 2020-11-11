include dev.env
export

.PHONY: fmt
fmt:
	cd server && gofmt -l -w .

.PHONY: lint
vet:
	cd server && go vet ./...

.PHONY: test
test.integration:
	docker-compose rm -s -f
	docker-compose --file docker-compose.local_test.yml up -d --force-recreate
	go test -count=1 -tags=integration ./integration_tests

dc.run:
	docker-compose up -d

dc.clean:
	docker-compose stop
	docker-compose rm

init.data:
	docker run --volume "`pwd`/sql/init_data.sql:/init_data.sql"  --network furutsu_network  -it --rm jbergknoff/postgresql-client postgres://postgres:postgres@furutsu_db:5432/furutsu \
	-f /init_data.sql
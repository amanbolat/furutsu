include dev.env
export

fmt:
	cd server && gofmt -l -w .

vet:
	cd server && go vet ./...

lint:
	cd server && golangci-lint run

test:
	cd server && go test -count=1 ./...

test.integration:
	docker-compose rm -s -f
	docker-compose --file docker-compose.integration_test.yml up -d --force-recreate
	cd server && go test -v -count=1 -tags=integration ./integration_tests

dc.run:
	docker-compose up -d --build

dc.clean:
	docker-compose rm -f -s

init.data:
	docker run --volume "`pwd`/sql/init_data.sql:/init_data.sql"  --network furutsu_network  -it --rm jbergknoff/postgresql-client postgres://postgres:postgres@furutsu_db:5432/furutsu \
	-f /init_data.sql
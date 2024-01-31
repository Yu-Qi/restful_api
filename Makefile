default:
	@go build -o build/ app/main.go

format:
	@go fmt `go list ./... | grep -v 'vendor'`

vet:
	@go vet `go list ./... | grep -v 'vendor'`

test:
	go test -v -vet "" ./... -race $(TEST_OPTS) ./...

lint:
	golangci-lint run --timeout 5m

gosec:
	gosec -quiet -include=G601 ./...

run:
	ENV=local \
	APP_PORT=8130 \
	go run app/main.go

live:
	APP_PORT=8129 \
	gin -i -p 8130 -a 8129 -t ./ -d app run

build:
	@go build -o build/ app/main.go

docker-build:
	docker build --build-arg CMD_DIR=app --build-arg APP_PORT=8130 -t restful-api:latest .

docker-run:
	docker run --rm -p 8130:8130 restful-api:latest
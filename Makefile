lint:
	golangci-lint run ./...

build_client:
	cd cmd/client && go build -v -ldflags="-X 'main.buildVersion = 1.0.0' -X 'main.buildDate = $(date)'" client

run_client:
	cd cmd/client && go run client -ldflags="-X 'main.buildVersion=v1.0.0' -X 'main.buildDate=$(date)'"

build_server:
	cd cmd/server && go build server

run_server:
	cd cmd/server && go run server
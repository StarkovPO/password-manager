lint:
	golangci-lint run ./...

build_client:
	cd cmd/client
	go build -v -ldflags="-X 'main.buildVersion = 1.0.0' -X 'main.buildDate = $(date)'"



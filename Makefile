test:
	CGO_ENABLED=1 go test --race ./...

cover:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

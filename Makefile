build:
	CGO_ENABLED=0 go build ./cmd/cassowary

test:
	go test -race -v ./...

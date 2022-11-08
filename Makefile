build:
	CGO_ENABLED=0 go build ./cmd/cassowary

build-linux:
	GOOS=linux CGO_ENABLED=0 go build ./cmd/cassowary

test:
	go test -race -v ./...

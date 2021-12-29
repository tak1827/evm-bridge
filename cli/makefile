install:
	go install -mod=readonly

.PHONY: test
test:
	go test -race ./...

bench:
	go test ./... -bench=. -benchtime=10s

fmt:
	go fmt ./...

lint:
	go vet ./...

build:
	go build -o bridgecli -gcflags '-m'

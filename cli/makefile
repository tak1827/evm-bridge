PROTO_SRC_FILES=$(shell find ./proto -type f -name "*.proto" | sed 's/\/proto//g')

.PHONY: proto
proto:
	rm ./pb/*.pb.go; \
	cd proto; \
	protoc -I=. -I=${GOPATH}/src/github.com/protobuf \
		--gofast_out=paths=source_relative:../pb  \
		$(PROTO_SRC_FILES)

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

PROTO_SRC_FILES=$(shell find ./proto -type f -name "*.proto" | sed 's/\/proto//g')

.PHONY: proto
proto:
	rm ./pb/*.pb.go; \
	cd proto; \
	protoc -I=. -I=${GOPATH}/src/github.com/protobuf \
		--gofast_out=paths=source_relative:../pb  \
		$(PROTO_SRC_FILES)

install:
	cd ./cmd/bridgecli/; \
	go install -mod=readonly

.PHONY: test
test:
	MallocNanoZone=0 go test -race ./...
# 	go test ./...

bench:
	go test ./... -bench=. -benchtime=10s

fmt:
	go fmt ./...

lint:
	go vet ./...

build:
	cd ./cmd/bridgecli/; \
	go build -o bridgecli -gcflags '-m'

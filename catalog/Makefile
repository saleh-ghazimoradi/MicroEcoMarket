generate-grpc:
	protoc --go_out=. --go-grpc_out=. gateway/proto/catalog.proto

fmt:
	go fmt ./...

vet:
	go vet ./...

run: fmt vet
	go run main.go

test:
	go test -v ./...
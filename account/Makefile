generate-grpc:
	protoc --go_out=. --go-grpc_out=. gateway/proto/account.proto

fmt:
	go fmt ./...

vet:
	go vet ./...

run: fmt vet
	go run main.go

test:
	go test -v ./...
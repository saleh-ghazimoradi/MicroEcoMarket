dockerUp:
	docker compose up -d

dockerDown:
	docker compose down

fmt:
	go fmt ./...

vet:
	go vet ./...

run: fmt vet
	go run main.go

test:
	go test -v ./...

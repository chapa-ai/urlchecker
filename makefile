run:
	go run cmd/main.go

build:
	go build -o urlchecker cmd/main.go

dc:
	docker compose build && docker compose up

test:
	go test -v ./...

mod:
	go mod tidy

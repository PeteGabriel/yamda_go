run:
	docker-compose up -d
	go build -o bin/main ./cmd/api && ./bin/main

test:
	go test ./...
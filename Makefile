.Phony: build run test docker-run clean 

build:
	go build -o bun/agent ./cmd/agent

run:
	go run ./cmd/agent

test:
	go test ./... -v

docker-run:
	docker-compose -f deployments/docker-compose.yml up --build

clean:
	rm -rf bin/
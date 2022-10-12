all:
	@go build -o app ./cmd/nats-service/main.go

publish:
	@go run publisher.go

test:
	@go test -v internal/app/app_test.go
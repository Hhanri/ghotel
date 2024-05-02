build:
	@go build -o bin/ghotel

run: build
	@./bin/ghotel

test:
	@go test ./... -v
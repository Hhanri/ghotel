build:
	@go build -o bin/ghotel

seed:
	@go run scripts/seed.go -dbUri=mongodb://ghotel:secret@localhost:27017/

run: build
	@./bin/ghotel

test:
	@go test ./... -v -count=1
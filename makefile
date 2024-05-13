build:
	@go build -o bin/ghotel

seed:
	@go run scripts/seed.go -dbUri=mongodb://ghotel:secret@localhost:27017/

run: build
	@./bin/ghotel -dbUri=mongodb://ghotel:secret@localhost:27017/

test:
	@go test ./... -v -count=1
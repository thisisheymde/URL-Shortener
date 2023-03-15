build:
	@go build -o bin/urlshortener

run: build
	@./bin/urlshortener
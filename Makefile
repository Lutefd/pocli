build:
	@go build -o bin/poke-repl cmd/main.go

run: build
	@./bin/poke-repl 

clean:
	@rm -rf bin
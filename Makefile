all: node

node:
	go build -o build/node cmd/node/main.go

clean:
	rm build/node
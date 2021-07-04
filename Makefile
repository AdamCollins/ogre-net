all: node

node:
	go build -o build/node cmd/node/main.go

client:
	go build -o build/local_proxy cmd/client/main.go

clean:
	rm build/node
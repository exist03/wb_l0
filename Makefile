nats:
	nats-server -js

publisher:
	go build -o publisher cmd/publisher/main.go

server:
	go run cmd/server/main.go
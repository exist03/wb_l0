publisher:
	go build -o publisher cmd/publisher/main.go

service:
	go run cmd/publisher/main.go
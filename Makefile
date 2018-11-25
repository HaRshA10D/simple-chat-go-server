test: 
	go test -p 1 -cover -v ./...

build:
	mkdir -p out/
	go build -o out/chat-server
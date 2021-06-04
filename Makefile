all: apiserver manageserver

apiserver:
	go build -o bin/apiserver cmd/apiserver/main.go

manageserver:
	go build -o bin/manageserver cmd/manageserver/main.go
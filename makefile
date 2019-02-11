build_api:
	-protoc --go_out=plugins=grpc:. api/*.proto

build_server:
	go build -o server/server ./server

run_server:
	go build -o server/server ./server && server/server

build_client:
	go build -o client/client ./client

build_client:
	go build -o client/client ./client && client/client


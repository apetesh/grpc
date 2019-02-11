package main

import (
	api "github.com/apetesh/grpc/api"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	phonebookClient := api.NewPhonebookClient(conn)
	_ = phonebookClient
}

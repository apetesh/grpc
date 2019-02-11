package main

import (
	"log"
	"os"
	"os/signal"
	"github.com/apetesh/grpc/server"
)

func main() {
	s := server.NewServer(8000)
	go func(){
		err := s.Start()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}()
	log.Printf("listening on port 8000")
	stopChan := make(chan os.Signal, 0)
	signal.Notify(stopChan, os.Kill, os.Interrupt)
	<-stopChan
	s.Stop()
}

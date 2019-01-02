package main

import (
	"os"
	"os/signal"
)

func main() {
	pb := NewPhonebook()
	s := NewServer(pb, 8000)
	go s.Start()
	stopChan := make(chan os.Signal, 0)
	signal.Notify(stopChan, os.Kill, os.Interrupt)
	<-stopChan
	s.Stop()
}

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"wqy/config"
)

func main() {
	client := NewWeQu(config.Name, config.Name, config.Key).Join()
	client.Subscribe()
	go client.Ping()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	client.Disconnect()
	fmt.Println("Disconnected from MQTT server")
}

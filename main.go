package main

import (
	"fmt"
	"log"
	"neural-network/server"
	"os"
)

func main() {
	fmt.Println("Starting the Neural Network server...")
	if err := server.StartServer(); err != nil {
		log.Fatalf("Error starting server: %v", err)
		os.Exit(1)
	}
}

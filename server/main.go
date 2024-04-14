package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sidsharma96/chatapp-backend/api"
)

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Port is not found in the environment!")
	}

	server := api.NewServer(port)
	server.Start()
}

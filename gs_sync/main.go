package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nrf24l01/rerandom/gs_sync/modes"
)

func main() {
	if os.Getenv("PRODUCTION_ENV") != "true" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("failed to load .env: %v", err)
		}
	}

	if len(os.Args) > 1 && os.Args[1] == "sync" {
		modes.RunSync()
	} else {
		panic("unknown command, use 'sync' to sync the sheet")
	}
}

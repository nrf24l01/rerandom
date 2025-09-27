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
	log.Printf("Arguments: %v\n", os.Args)

	if len(os.Args) > 1 && os.Args[1] == "sync" {
		modes.RunSync()
	} else if len(os.Args) > 1 && os.Args[1] == "webserver" {
		modes.RunWebserver()
	} else {
		panic("unknown command, use 'sync' to sync the sheet, or 'webserver' to start the webserver")
	}
}

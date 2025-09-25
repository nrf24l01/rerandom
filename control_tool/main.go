package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nrf24l01/rerandom/control_tool/core"
	"github.com/nrf24l01/rerandom/control_tool/database"
	"github.com/nrf24l01/rerandom/control_tool/tasks"
)

func main() {
	if os.Getenv("PRODUCTION_ENV") != "true" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("failed to load .env: %v", err)
		}
	}

	cfg := core.BuildConfigFromEnv()
	db := database.RegisterPostgres(cfg)

	h := tasks.Handler{
		DB: db,
		Config: cfg,
	}

	argsWithoutProg := os.Args[1:]

	if (len(argsWithoutProg) == 3 || len(argsWithoutProg) == 2) && argsWithoutProg[0] == "create-user" {
		log.Printf("Creating user %s", argsWithoutProg[1])
		if len(argsWithoutProg) == 2 {
			log.Printf("Password not provided as argument, reading from stdin")
			var password string
			_, err := fmt.Scanln(&password)
			if err != nil {
				log.Fatalf("Failed to read password from stdin: %v", err)
			}
			argsWithoutProg = append(argsWithoutProg, password)
		}
		user, password := argsWithoutProg[1], argsWithoutProg[2]
		if err := h.CreateUser(user, password); err != nil {
			log.Printf("Error creating user: %v", err)
		} else {
			log.Printf("User %s created successfully", user)
		}
	}
}
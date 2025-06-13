package main

import (
	"book_rest_api/internal/config"
	"fmt"
	"log"
)

func init() {
	fmt.Println("Initializing application...")
	_, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
}

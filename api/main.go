package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/gc"
	"github.com/sarmerer/forum/api/logger"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/router"

	_ "github.com/mattn/go-sqlite3"
)

// Init does necessary preparations to successfully run the API
func Init() {
	log.Println("Starting server...")

	gc.Start()
	logger.InitLogs("Garbage collector", nil)

	err := repository.InitDB()
	logger.InitLogs("Database init", err)

	err = repository.CheckDBIntegrity()
	logger.InitLogs("Database integrity", err)
}

// Run sets up API routes and then starts to listen for requests
func Run() {
	mux := router.New()
	mux.SetupRoutes()
	port := config.APIPort
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	log.Printf("Listening https://localhost:%s\n", port)
	////log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.APIPort), "./ssl/cert.pem", "./ssl/key.pem", mux))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}

func main() {
	Init()
	Run()
}

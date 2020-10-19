package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/gc"
	"github.com/sarmerer/forum/api/logger"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/router"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	Init()
	Run()
}

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
	log.Printf("Listening https://localhost:%d\n", config.APIPort)
	////log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.APIPort), "./ssl/cert.pem", "./ssl/key.pem", mux))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.APIPort), mux))
}

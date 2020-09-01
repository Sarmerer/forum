package api

import (
	"fmt"
	"forum/api/gc"
	"forum/api/logger"
	"forum/api/router"
	"forum/api/utils"
	"forum/config"
	"forum/database"
	"log"
	"net/http"
)

// Init does necessary preparations to successfully run the API
func Init() {
	log.Println("Starting server...")

	err := utils.LoadEnv(".env")
	logger.InitLogs("Environment", err)

	gc.Start()
	logger.InitLogs("Garbage collector", nil)

	err = database.CheckIntegrity()
	logger.InitLogs("Database integrity", err)
}

// Run sets up API routes and then starts to listen for requests
func Run() {
	mux := router.New()
	mux.SetupRoutes()
	log.Printf("Listening https://localhost:%d\n", config.APIPort)
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.APIPort), "./ssl/cert.pem", "./ssl/key.pem", mux))
}

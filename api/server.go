package api

import (
	"fmt"
	"forum/api/logger"
	"forum/api/router"
	"forum/api/sessions"
	"forum/api/utils"
	"forum/config"
	"forum/database"
	"log"
	"net/http"
)

func Init() {
	log.Println("Starting server...")

	err := utils.LoadEnv(".env")
	logger.InitLogs("Environment", err)

	sessions.StartGC()
	logger.InitLogs("Garbage collector", nil)

	err = database.CheckIntegrity()
	logger.InitLogs("Database integrity", err)
}

func Run() {
	mux := router.New()
	log.Printf("Listening https://localhost:%d\n", config.APIPort)
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.APIPort), "./ssl/cert.pem", "./ssl/key.pem", mux))
}

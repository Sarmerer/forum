package api

import (
	"fmt"
	"forum/api/cache"
	"forum/api/database"
	"forum/api/router"
	"forum/api/utils"
	"forum/config"
	"log"
	"net/http"
	"os"
	"time"
)

func Init() {
	log.Println("Starting server...")

	if err := utils.LoadEnv(".env"); err != nil {
		log.Printf("Environment:\t|%s|\tError: %s\n", utils.Red("FAIL"), err)
		os.Exit(1)
	}
	log.Printf("Environment:\t|%s|\n", utils.Green("OK"))

	cache.Sessions = cache.NewManager(14*24*time.Hour, 2*time.Hour)
	log.Printf("Sessions manager:\t|%s|\n", utils.Green("OK"))

	if err := database.CheckIntegrity(); err != nil {
		log.Printf("Database integrity:\t|%s|\tError: %s\n", utils.Red("FAIL"), err)
		os.Exit(1)
	}
	log.Printf("Databse integrity:\t|%s|\n", utils.Green("OK"))
}

func Run() {
	log.Printf("Listening https://localhost:%d\n", config.APIPort)
	mux := router.New()
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.APIPort), "./ssl/cert.pem", "./ssl/key.pem", mux))
}

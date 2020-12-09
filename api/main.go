package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/gc"
	"github.com/sarmerer/forum/api/logger"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/router"
	"github.com/sarmerer/forum/api/utils"

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

	err = utils.SetupEnv()
	logger.CheckErrAndLog("Environment", "loaded", err)

	flags := utils.ParseFlags(os.Args[1:])
	for _, flag := range flags {
		logger.Log("Config", flag)
	}
	config.Init()
}

// Run sets up API routes and then starts to listen for requests
func Run() {
	mux := router.New()
	mux.SetupRoutes()
	port := config.APIPort
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	log.Printf("Listening %s\n", config.APIURL)
	if config.Production {
		log.Fatal(http.ListenAndServeTLS(":"+port, "./ssl/cert.pem", "./ssl/key.pem", mux))
	} else {
		log.Fatal(http.ListenAndServe(":"+port, mux))
	}
}

func main() {
	Init()
	Run()
}

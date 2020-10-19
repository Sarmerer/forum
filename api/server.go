// package api

// import (
// 	"fmt"
// 	"forum/api/config"
// 	"forum/api/gc"
// 	"forum/api/logger"
// 	"forum/api/repository"
// 	"forum/api/router"
// 	"log"
// 	"net/http"
// )

// // Init does necessary preparations to successfully run the API
// func Init() {
// 	log.Println("Starting server...")

// 	gc.Start()
// 	logger.InitLogs("Garbage collector", nil)

// 	err := repository.InitDB()
// 	logger.InitLogs("Database init", err)

// 	err = repository.CheckDBIntegrity()
// 	logger.InitLogs("Database integrity", err)
// }

// // Run sets up API routes and then starts to listen for requests
// func Run() {
// 	mux := router.New()
// 	mux.SetupRoutes()
// 	log.Printf("Listening https://localhost:%d\n", config.APIPort)
// 	////log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.APIPort), "./ssl/cert.pem", "./ssl/key.pem", mux))
// 	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.APIPort), mux))
// }

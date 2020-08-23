package api

import (
	"fmt"
	"forum/api/cache"
	"forum/api/router"
	"forum/api/utils"
	"forum/config"
	"log"
	"net/http"
	"time"
)

func Init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := utils.LoadEnv(".env"); err != nil {
		log.Fatal("Could not launch the server. Error:", err)
	}
	cache.Sessions = cache.NewManager(14*24*time.Hour, 2*time.Hour)
}

func Run() {
	fmt.Printf("\nListening https://localhost:%d\n", config.APIPort)
	mux := router.New()
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.APIPort), "./ssl/cert.pem", "./ssl/key.pem", mux))
}

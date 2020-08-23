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

func Run() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Printf("\nListening https://localhost:%d\n", config.APIPort)
	if err := utils.LoadEnv(".env"); err != nil {
		log.Fatal("Could not launch the server. Error:", err)
	}
	cache.Sessions = cache.NewManager(14*24*time.Hour, 30*time.Minute)
	fmt.Println(cache.Sessions)
	mux := router.New()
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.APIPort), "./ssl/cert.pem", "./ssl/key.pem", mux))
}

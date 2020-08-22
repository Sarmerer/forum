package api

import (
	"fmt"
	"forum/api/router"
	"forum/api/utils"
	"forum/config"
	"log"
	"net/http"
)

func Run() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Printf("\nListening https://localgost:%d\n", config.APIPort)
	if err := utils.LoadEnv(".env"); err != nil {
		log.Fatal("Could not launch the server. Error:", err)
	}
	mux := router.New()
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.APIPort), "./ssl/cert.pem", "./ssl/key.pem", mux))
}

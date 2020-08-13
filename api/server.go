package api

import (
	"fmt"
	"forum/api/controllers"
	"forum/api/router/routes"
	"forum/config"
	"log"
	"net/http"
)

func Run() {
	go http.ListenAndServe(fmt.Sprintf(":%d", config.HTTPport), http.HandlerFunc(controllers.RedirectToHTTPS))
	fmt.Printf("\nListening [::]:%d\n", config.HTTPSport)
	mux := routes.New()
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.HTTPSport), "./ssl/cert.pem", "./ssl/key.pem", mux))
}

package api

import (
	routers "forum/api/controllers/router"
	"log"
	"net/http"
)

func Run() {
	mux := routers.New()
	log.Fatal(http.ListenAndServe(":8080", mux))
}

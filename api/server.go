package api

import (
	"fmt"
	"forum/api/controllers/router"
	"forum/config"
	"forum/database"
	"log"
	"net/http"
)

func Run() {
	//initDB()
	fmt.Printf("\nListening [::]:%d\n", config.Port)
	mux := router.New()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), mux))
}

func initDB() {
	db, err := database.Connect()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(db)
}

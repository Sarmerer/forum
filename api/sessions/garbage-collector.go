package sessions

import (
	"forum/api/database"
	"forum/api/models"
	"forum/config"
	"log"
	"time"
)

//StartGC starts garbage collector
func StartGC() {
	go GarbageCollector()
}

func GarbageCollector() {
	for {
		<-time.After(config.GCInterval)
		db, dbErr := database.Connect()
		if dbErr != nil {
			log.Fatal("Garbage collector could not start. Error: ", dbErr)
		}
		um, umErr := models.NewUserModel(db)
		if umErr != nil {
			log.Fatal("Garbage collector could not start. Error: ", umErr)
			return
		}
		users, uErr := um.FindAll()
		if uErr != nil {
			log.Fatal("Garbage collector could not start. Error: ", uErr)
		}
		for _, user := range users {
			if time.Now().After(user.LastOnline.Add(14 * 24 * time.Hour)) {
				um.UpdateSession(user.ID, "")
			}
		}
	}
}

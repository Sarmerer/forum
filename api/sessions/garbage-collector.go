package sessions

import (
	models "forum/api/models/user"
	"forum/config"
	"forum/database"
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
		start := time.Now()
		counter := 0
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
			if time.Now().After(user.LastOnline.Add(config.SessionExpiration)) {
				um.UpdateSession(user.ID, "")
				counter++
			}
		}
		log.Printf("Garbage collector has destroyed expired sessions. Time took: %s, sessions closed: %v", time.Since(start), counter)
	}
}

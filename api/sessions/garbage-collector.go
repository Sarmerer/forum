package sessions

import (
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
		<-time.After(config.SessionExpiration)
		um, umErr := models.NewUserModel()
		defer um.DB.Close()
		if umErr != nil {
			log.Fatal("Garbage collector could not start. Error: ", umErr)
			return
		}
		users, uErr := um.FindAll()
		if uErr != nil {
			log.Fatal("Garbage collector could not start. Error: ", umErr)
		}
		for _, user := range users {
			if time.Now().After(user.LastOnline.Add(14 * 24 * time.Hour)) {
				um.UpdateSession(user.ID, "")
			}
		}
	}
}

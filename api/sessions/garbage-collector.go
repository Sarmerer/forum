package sessions

import (
	"fmt"
	"forum/api/logger"
	models "forum/api/models/user"
	"forum/config"
	"forum/database"
	"time"
)

//StartGC starts garbage collector
func StartGC() {
	go GarbageCollector()
}

func GarbageCollector() {
	for {
		<-time.After(config.GCInterval)
		counter := 0
		start := time.Now()
		db, dbErr := database.Connect()
		if dbErr != nil {
			logger.ServerLogs("Garbage collector", "", dbErr)
			return
		}
		um, umErr := models.NewUserModel(db)
		if umErr != nil {
			logger.ServerLogs("Garbage collector", "", umErr)
			return
		}
		users, uErr := um.FindAll()
		if uErr != nil {
			logger.ServerLogs("Garbage collector", "", uErr)
			return
		}
		for _, user := range users {
			if time.Now().After(user.LastOnline.Add(config.SessionExpiration)) {
				um.UpdateSession(user.ID, "")
				counter++
			}
		}
		logger.ServerLogs("Garbage collector", fmt.Sprintf("closed %v sessions. Took %s", counter, time.Since(start)), nil)
	}
}

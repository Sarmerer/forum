package gc

import (
	"fmt"
	"forum/api/logger"
	"forum/api/repository/crud"
	"forum/config"
	"forum/database"
	"time"
)

// Start starts the garbage collector service,
// it will remove expired sessions from database with the interval, specified in the config file
func Start() {
	go garbageCollector()
}

func garbageCollector() {
	for {
		<-time.After(config.GCInterval)
		counter := 0
		start := time.Now()
		db, dbErr := database.Connect()
		if dbErr != nil {
			logger.ServerLogs("Garbage collector", "", dbErr)
			return
		}
		um := crud.NewUserModel(db)
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

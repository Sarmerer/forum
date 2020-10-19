package gc

import (
	"fmt"
	"time"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/logger"
	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
)

// Start starts the garbage collector service,
// it will remove expired sessions from database with the interval, specified in the config file
func Start() {
	go garbageCollector()
}

func garbageCollector() {
	var (
		um         repository.UserRepo
		lastOnline time.Time
		users      []models.User
		err        error
	)
	um = crud.NewUserRepoCRUD()
	for {
		<-time.After(config.GCInterval)
		counter := 0
		start := time.Now()
		if users, err = um.FindAll(); err != nil {
			logger.ServerLogs("Garbage collector", "", err)
			return
		}
		for _, user := range users {
			if lastOnline, err = time.Parse(config.TimeLayout, user.LastOnline); err != nil {
				logger.ServerLogs("Garbage collector", "time parsing error", err)
			}
			if time.Now().After(lastOnline.Add(config.SessionExpiration)) {
				if err = um.UpdateSession(user.ID, ""); err != nil {
					logger.ServerLogs("Garbage collector", "session closing error", err)
				}
				counter++
			}
		}
		logger.ServerLogs("Garbage collector", fmt.Sprintf("closed %v sessions. Took %s", counter, time.Since(start)), nil)
	}
}

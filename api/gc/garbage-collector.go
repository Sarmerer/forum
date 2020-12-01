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
		repo       repository.UserRepo = crud.NewUserRepoCRUD()
		lastOnline time.Time
		users      []models.User
		err        error
	)
	for {
		<-time.After(config.GCInterval)
		counter := 0
		start := time.Now()
		if users, err = repo.FindAll(); err != nil {
			logger.CheckErrAndLog("Garbage collector", "", err)
			return
		}
		for _, user := range users {
			if lastOnline, err = time.Parse(config.TimeLayout, user.LastActive); err != nil {
				logger.CheckErrAndLog("Garbage collector", "time parsing error", err)
			}
			if time.Now().After(lastOnline.Add(config.SessionExpiration)) {
				if err = repo.UpdateSession(user.ID, ""); err != nil {
					logger.CheckErrAndLog("Garbage collector", "session closing error", err)
				}
				counter++
			}
		}
		logger.CheckErrAndLog("Garbage collector", fmt.Sprintf("closed %v sessions. Took %s", counter, time.Since(start)), nil)
	}
}

package garbagecollector

import (
	"time"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/logger"
	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
)

// Start launches the garbage collector service,
// it will remove expired sessions from database with the interval, specified in the config file
func Start() {
	go garbageCollector()
}

func garbageCollector() {
	go closeSessions()
}

func closeSessions() {
	var (
		repo  repository.UserRepo = crud.NewUserRepoCRUD()
		users []models.User
		err   error
	)
	for {
		time.Sleep(config.GCInterval)
		deauthedUsers := 0
		if users, err = repo.FindAll(); err != nil {
			logger.CheckErrAndLog("Garbage collector", "", err)
			return
		}
		for _, user := range users {
			if time.Now().After(time.Unix(user.LastActive, 0).Add(config.SessionExpiration)) {
				if err = repo.UpdateSession(user.ID, ""); err != nil {
					logger.CheckErrAndLog("Garbage collector", "session closing error", err)
				}
				deauthedUsers++
			}
		}
	}
}

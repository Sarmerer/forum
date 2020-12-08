package crud

import (
	"fmt"
	"net/http"

	"github.com/sarmerer/forum/api/database"
	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (UserRepoCRUD) GetPassword(userID int64) (password string, status int, err error) {
	var (
		conn *database.MongoDatastore
	)
	if conn, err = database.Connect("users"); err != nil {
		return "", http.StatusInternalServerError, err
	}
	type fwa struct {
		Password string `bson:"password"`
	}
	fwafwa := fwa{}
	ctx, cancel := utils.Ctx()
	defer cancel()
	f := conn.Collection.FindOne(ctx, bson.D{{Key: "login", Value: "banana"}}, options.FindOne().SetProjection(bson.D{{"_id", false}, {"password", true}}))
	fmt.Println(f.Decode(&fwafwa))
	fmt.Println(fwafwa)
	return password, http.StatusOK, err
}

func (UserRepoCRUD) ValidateSession(sessionID string) (user models.UserCtx, status int, err error) {
	return user, http.StatusOK, nil
}

func (UserRepoCRUD) UpdateSession(id int64, newSession string) error {

	return nil
}

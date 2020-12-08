package crud

import (
	"errors"
	"net/http"

	"github.com/sarmerer/forum/api/database"
	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//UserRepoCRUD helps performing CRUD operations
type UserRepoCRUD struct{}

var DeletedUser = &models.User{
	ID:          -1,
	Login:       "deleted",
	Email:       "deleted",
	Avatar:      "https://avatars.dicebear.com/api/male/deleted",
	DisplayName: "deleted",
	Role:        0,
}

//NewUserRepoCRUD creates an instance of UserModel
func NewUserRepoCRUD() UserRepoCRUD {
	return UserRepoCRUD{}
}

//FindAll returns all users in the database
func (UserRepoCRUD) FindAll() ([]models.User, error) {
	var (
		users []models.User
	)

	return users, nil
}

//FindByID returns a specific user from the database
func (UserRepoCRUD) FindByID(userID int64) (user *models.User, status int, err error) {
	var conn *database.MongoDatastore
	if conn, err = database.Connect("users"); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	ctx, cancel := utils.Ctx()
	defer cancel()
	if err = conn.Collection.
		FindOne(ctx, bson.D{{Key: "id", Value: userID}}, options.FindOne().SetProjection(bson.D{{Key: "password", Value: false}, {Key: "session_id", Value: false}})).Decode(&user); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return user, http.StatusOK, err
}

//Create adds a new user to the database
func (UserRepoCRUD) Create(user *models.User) (*models.User, int, error) {
	var (
		conn        *database.MongoDatastore
		loginExists int64
		emailExists int64
		result      *mongo.InsertOneResult
		newUser     *models.User
		err         error
	)
	if conn, err = database.Connect("users"); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	ctx, cancel := utils.Ctx()
	defer cancel()
	if loginExists, err = conn.Collection.CountDocuments(ctx, bson.D{
		{Key: "login", Value: user.Login},
	},
	); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if emailExists, err = conn.Collection.CountDocuments(ctx, bson.D{
		{Key: "email", Value: user.Email},
	},
	); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if loginExists > 0 || emailExists > 0 {
		var message string
		if loginExists > 0 && emailExists > 0 {
			message = "login and email are already taken"
		} else if loginExists > 0 {
			message = "login is already taken"
		} else if emailExists > 0 {
			message = "email is already taken"
		}
		return nil, http.StatusBadRequest, errors.New(message)
	}

	if result, err = conn.Collection.InsertOne(ctx, user); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if err = conn.Collection.FindOne(ctx, bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&newUser); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return newUser, http.StatusOK, nil
}

//Update updates existing user in the database
func (UserRepoCRUD) Update(user *models.User) (int, error) {

	return http.StatusNotModified, errors.New("could not update the user")
}

func (UserRepoCRUD) UpdateLastActivity(userID int64) error {
	var (
		conn *database.MongoDatastore
		err  error
	)
	if conn, err = database.Connect("users"); err != nil {
		return err
	}
	ctx, cancel := utils.Ctx()
	defer cancel()
	conn.Collection.FindOneAndUpdate(ctx, bson.D{{Key: "id", Value: userID}}, bson.D{{Key: "last_activity", Value: utils.CurrentUnixTime()}})
	return nil
}

//Delete deletes user from the database
func (UserRepoCRUD) Delete(userID int64) (int, error) {
	return http.StatusNotModified, errors.New("could not delete the user")
}

//FindByNameOrEmail finds a user by name or email in the database
func (UserRepoCRUD) FindByLoginOrEmail(login string) (*models.User, int, error) {
	var (
		u   models.User
		err error
	)

	return &u, http.StatusOK, err
}

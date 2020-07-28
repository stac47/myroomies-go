package data

import (
	"context"
	"time"

	"github.com/stac47/myroomies/pkg/models"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	usersCollection = "users"
)

type MongoUserDataAccess struct {
	db *mongo.Database
}

func (dao *MongoUserDataAccess) getCollection() *mongo.Collection {
	return dao.db.Collection(usersCollection)
}

func GetMongoUserDataAccess(db *mongo.Database) UserDataAccess {
	return &MongoUserDataAccess{db}
}

func (dao *MongoUserDataAccess) RetrieveUsers(ctx context.Context) []models.User {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout*time.Second)
	defer cancel()
	cur, err := dao.getCollection().Find(ctx, bson.D{})
	if err != nil {
		log.Errorf("Cannot retrieve the users: %s", err)
		return nil
	}
	defer cur.Close(ctx)
	users := make([]models.User, 0)
	for cur.Next(ctx) {
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			log.Errorf("Cannot decode a user: %s", err)
			return nil
		}
		users = append(users, user)
	}
	return users
}

func (dao *MongoUserDataAccess) CreateUser(ctx context.Context, newUser models.User) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout*time.Second)
	defer cancel()
	_, err := dao.getCollection().InsertOne(ctx, newUser)
	return err
}

func (dao *MongoUserDataAccess) RetrieveUser(ctx context.Context, login string) *models.User {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout*time.Second)
	defer cancel()
	filter := bson.M{"login": login}
	var user models.User
	err := dao.getCollection().FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Errorf("Cannot retrieve user [%s]: %s", login, err)
		return nil
	}
	return &user
}

func (dao *MongoUserDataAccess) DeleteUser(ctx context.Context, user models.User) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout*time.Second)
	defer cancel()
	filter := bson.M{"login": user.Login}
	result, err := dao.getCollection().DeleteOne(ctx, filter)
	if err != nil {
		log.Errorf("Cannot delete the user [%s]: %s", user.Login, err)
		return err
	}
	log.Debugf("Number of deleted users: %d", result.DeletedCount)
	return nil
}

func (dao *MongoUserDataAccess) UpdateUser(ctx context.Context, user models.User) (err error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout*time.Second)
	defer cancel()
	filter := bson.M{"login": user.Login}
	update := bson.D{
		{"$set", bson.D{
			{"firstname", user.Firstname},
			{"lastname", user.Lastname},
			{"password", user.Password}}}}
	result := dao.getCollection().FindOneAndUpdate(ctx, filter, update)
	return result.Err()
}

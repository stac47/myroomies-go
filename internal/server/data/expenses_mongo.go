package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/stac47/myroomies/pkg/models"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	expensesCollection = "expenses"
)

func (dao *MongoExpenseDataAccess) getCollection() *mongo.Collection {
	return dao.db.Collection(expensesCollection)
}

type MongoExpenseDataAccess struct {
	db *mongo.Database
}

func GetMongoExpenseDataAccess(db *mongo.Database) ExpenseDataAccess {
	return &MongoExpenseDataAccess{db}
}

func (dao *MongoExpenseDataAccess) CreateExpense(ctx context.Context, newExpense models.Expense) (models.Expense, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout*time.Second)
	defer cancel()
	result, err := dao.getCollection().InsertOne(ctx, newExpense)
	if err != nil {
		log.Errorf("Cannot insert expense: %v", newExpense)
		return models.Expense{}, err
	}
	newExpense.Id = result.InsertedID.(primitive.ObjectID).Hex()
	return newExpense, nil
}

func (dao *MongoExpenseDataAccess) RetrieveExpenseFromId(ctx context.Context, id string) *models.Expense {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout*time.Second)
	defer cancel()
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}
	var expense models.Expense
	err := dao.getCollection().FindOne(ctx, filter).Decode(&expense)
	if err != nil {
		log.Errorf("Cannot retrieve expense [%s]: %s", id, err)
		return nil
	}
	return &expense
}

func (dao *MongoExpenseDataAccess) UpdateExpense(ctx context.Context, updatedExpense models.Expense) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout*time.Second)
	defer cancel()
	objectId, _ := primitive.ObjectIDFromHex(updatedExpense.Id)
	filter := bson.M{"_id": objectId}
	// Clear the _id
	updatedExpense.Id = ""
	result := dao.getCollection().FindOneAndReplace(ctx, filter, updatedExpense)
	if result == nil {
		msg := fmt.Sprintf("Expense [%s] was not updated", updatedExpense.Id)
		return errors.New(msg)
	} else if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (dao *MongoExpenseDataAccess) DeleteExpense(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout*time.Second)
	defer cancel()
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}
	result, err := dao.getCollection().DeleteOne(ctx, filter)
	if err != nil {
		log.Errorf("Cannot delete expense [%s]: %s", id, err)
		return err
	}
	if result.DeletedCount == 0 {
		log.Warnf("The expense [%s] was not found: cannot delete it", id)
		return errors.New("Deletion failed")
	}
	return nil
}

func (dao *MongoExpenseDataAccess) RetrieveExpenses(ctx context.Context) []models.Expense {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout*time.Second)
	defer cancel()
	cur, err := dao.getCollection().Find(ctx, bson.D{})
	if err != nil {
		log.Errorf("Cannot retrieve the expenses: %s", err)
		return nil
	}
	defer cur.Close(ctx)
	expenses := make([]models.Expense, 0)
	for cur.Next(ctx) {
		var expense models.Expense
		err := cur.Decode(&expense)
		if err != nil {
			log.Errorf("Cannot decode a expense: %s", err)
			return nil
		}
		expenses = append(expenses, expense)
	}
	return expenses
}

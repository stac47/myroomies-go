package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	log "github.com/sirupsen/logrus"
)

const (
	myroomiesDatabase = "myroomies"
)

// Factory to access data from MongoDB

type MongoDataAccessParams struct {
	Server string
}

type mongoDataAccessFactory struct {
	params      MongoDataAccessParams
	mongoClient *mongo.Client
}

func connectMongo(url string, timeoutSeconds int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return mongo.Connect(ctx, options.Client().ApplyURI(url))
}

func NewMongoDataAccessFactory(params MongoDataAccessParams) *mongoDataAccessFactory {
	const (
		// TODO: Make this parameterable
		connectionTimeoutSeconds = 1
		maxConnectionTries       = 5
		retryPeriodSeconds       = 5
	)

	tryCounter := 1
	currentRetryPeriodSeconds := retryPeriodSeconds
	var mongoClient *mongo.Client
	for {
		if tryCounter == maxConnectionTries {
			panic("Error connecting to MongoDB")
		}
		log.Infof("Connection #%d to MongoDB [%s]...", tryCounter, params.Server)
		var err error
		mongoClient, err = connectMongo(params.Server, connectionTimeoutSeconds)
		if err != nil {
			tryCounter++
			currentRetryPeriodSeconds *= 2
			log.Warnf("Connection #%d to MongoDB [%s] failed. Next try in %d seconds",
				tryCounter, params.Server, currentRetryPeriodSeconds)
		} else {
			log.Infof("Connection to MongoDB [%s] is OK (handle=%p)", params.Server, mongoClient)
			break
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err := mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		panic("Error connecting to MongoDB: " + err.Error())
	}
	return &mongoDataAccessFactory{
		params:      params,
		mongoClient: mongoClient,
	}
}

func (f mongoDataAccessFactory) GetUserDataAccess() UserDataAccess {
	return GetMongoUserDataAccess(f.mongoClient.Database(myroomiesDatabase))
}

func (f mongoDataAccessFactory) GetExpenseDataAccess() ExpenseDataAccess {
	return GetMongoExpenseDataAccess(f.mongoClient.Database(myroomiesDatabase))
}

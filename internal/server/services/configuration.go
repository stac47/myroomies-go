package services

import (
	"github.com/stac47/myroomies/internal/server/data"

	log "github.com/sirupsen/logrus"
)

var (
	dataAccessFactory data.DataAccessFactory
)

func GetDataAccess() data.DataAccessFactory {
	if dataAccessFactory == nil {
		log.Fatal("Error: you have to initialize the services with ConfigureServices")
	}
	return dataAccessFactory
}

func Configure(dataAccessParams interface{}) {
	log.Info("Configuring service layer")
	if dataAccessFactory != nil {
		log.Warnf("Data access layer already configured [%p]", dataAccessFactory)
		return
	}
	log.Info("Creating data access layer")
	dataAccessFactory = data.NewDataAccessFactory(dataAccessParams)
	log.Infof("Data access factory created at [%p]", dataAccessFactory)
}

func Shutdown() {
	dataAccessFactory.Close()
	log.Info("Database connection closed.")
}

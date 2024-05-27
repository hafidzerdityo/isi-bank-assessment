package services

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hafidzresttemplate.com/datastore"
)


type ServiceSetup struct{
	Logger *logrus.Logger
	Datastore *datastore.DatastoreSetup
	Db		*gorm.DB
}




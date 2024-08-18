package services

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hafidzresttemplate.com/datastore"
	"hafidzresttemplate.com/startup"
)


type ServiceSetup struct{
	Logger *logrus.Logger
	Datastore *datastore.DatastoreSetup
	Db		*gorm.DB
	EventPub startup.EventStreamConfig
}


func (s *ServiceSetup)TransactionStatusHandler(tx *gorm.DB, err *error) {
	panicMessage := recover()
	if panicMessage != nil {
		s.Logger.Warn(logrus.Fields{"panic": panicMessage}, nil, "transaction failed, rollback")
		tx.Rollback()
		*err = fmt.Errorf("panic recovered")
	}
	if *err != nil {
		s.Logger.Warn(logrus.Fields{"error": (*err).Error()}, nil, "transaction failed, rollback")
		tx.Rollback()
		return
	} 
	s.Logger.Info(logrus.Fields{}, nil, "transaction success, committed ")
	tx.Commit()
}


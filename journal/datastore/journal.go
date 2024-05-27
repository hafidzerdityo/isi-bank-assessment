package datastore

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hafidzresttemplate.com/dao"
)


func(d *DatastoreSetup) InsertJournal(tx *gorm.DB, reqPayload dao.Journal)(err error){
	d.Logger.Info(
		logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayload)}, nil, "START: InsertJournal Datastore",
	)

	if err = tx.Create(&reqPayload).Error; err != nil {
		d.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, err.Error(),
		)
		return
	}
	remark := "END: InsertJournal Datastore"
	d.Logger.Info(
		logrus.Fields{"response": fmt.Sprintf("%+v", map[string]interface{}{})}, nil, remark,
	)

	return
}
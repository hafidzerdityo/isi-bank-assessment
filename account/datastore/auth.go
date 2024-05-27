package datastore

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hafidzresttemplate.com/dao"
)


func (d *DatastoreSetup) GetAccount(tx *gorm.DB, noRekening string) (datastoreResponse dao.Account, err error) {
	reqPayloadForLog := map[string]interface{}{
		"noRekening": noRekening,
	}
	d.Logger.Info(
		logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayloadForLog)}, nil, "START: GetAccount Datastore",
	)

	err = tx.Where("no_rekening = ?", noRekening).Take(&datastoreResponse).Error
	if err != nil {
		d.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, err.Error(),
		)
		return
	}

	datastoreResponseForLog := datastoreResponse
	datastoreResponseForLog.HashedPin = "*REDACTED*"
	remark := "END: GetAccount Datastore"
	d.Logger.Info(
		logrus.Fields{"response": fmt.Sprintf("%+v", datastoreResponseForLog)}, nil, remark,
	)

	return
}

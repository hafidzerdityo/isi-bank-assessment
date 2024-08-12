package datastore

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hafidzresttemplate.com/dao"
)


func (d *DatastoreSetup) GetMutasi(tx *gorm.DB, ID int, limit int, page int) (datastoreResponse []dao.Transaction, count int64 , err error) {
	reqPayloadForLog := map[string]interface{}{
	  "ID":    ID,
	  "limit": limit,
	  "offset": page,
	}
	d.Logger.Info(
	  logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayloadForLog)}, nil, "START: GetMutasi Datastore",
	)
  
	err = tx.Where("id_rekening = ?", ID).Offset((page - 1) * limit).Limit(limit).Find(&datastoreResponse).Error
	
	if err != nil {
	  d.Logger.Error(
		logrus.Fields{"error": err.Error()}, nil, err.Error(),
	  )
	  return
	}

	err = tx.Model(&dao.Transaction{}).Where("id_rekening = ?", ID).Count(&count).Error

	if err != nil {
		d.Logger.Error(
		logrus.Fields{"error": err.Error()}, nil, err.Error(),
		)
		return
	}

	remark := "END: GetMutasi Datastore"
	d.Logger.Info(
	  logrus.Fields{"response": fmt.Sprintf("%+v", datastoreResponse)}, nil, remark,
	)
  
	return
  }


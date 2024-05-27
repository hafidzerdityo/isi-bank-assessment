package datastore

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hafidzresttemplate.com/dao"
)


func(d *DatastoreSetup) UpdateSaldo(tx *gorm.DB, reqPayload dao.CreateTabungTarikUpdate)(saldo float64, err error){
	d.Logger.Info(
		logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayload)}, nil, "START: UpdateSaldo Datastore",
	)
	var model dao.Account
	if err = tx.Model(&model).Where("no_rekening = ?", reqPayload.NoRekening).UpdateColumn("saldo", gorm.Expr("saldo + ?", reqPayload.Nominal)).Error; err != nil{
		d.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, err.Error(),
		)
		return
	}
    
    if err = tx.Model(&model).Select("saldo").Where("no_rekening = ?", reqPayload.NoRekening).Scan(&saldo).Error; err != nil {
		d.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, err.Error(),
		)
		return
    }

	remark := "END: UpdateSaldo Datastore"
	d.Logger.Info(
		logrus.Fields{"response": fmt.Sprintf("%+v",  map[string]interface{}{})}, nil, remark,
	)

	return
}

func(d *DatastoreSetup) InsertCatatan(tx *gorm.DB, reqPayload dao.Transaction)(err error){
	d.Logger.Info(
		logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayload)}, nil, "START: InsertCatatan Datastore",
	)

	if err = tx.Create(&reqPayload).Error; err != nil {
		d.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, err.Error(),
		)
		return
	}

	remark := "END: InsertCatatan Datastore"
	d.Logger.Info(
		logrus.Fields{"response": fmt.Sprintf("%+v",  map[string]interface{}{})}, nil, remark,
	)

	return
}
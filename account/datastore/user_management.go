package datastore

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hafidzresttemplate.com/dao"
)


func(d *DatastoreSetup) InsertCustomer(tx *gorm.DB, reqPayload dao.Customer)(err error){
	d.Logger.Info(
		logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayload)}, nil, "START: InsertCustomer Datastore",
	)

	if err = tx.Create(&reqPayload).Error; err != nil {
		d.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, err.Error(),
		)
		return
	}
	remark := "END: InsertCustomer Datastore"
	d.Logger.Info(
		logrus.Fields{"response": fmt.Sprintf("%+v", map[string]interface{}{})}, nil, remark,
	)

	return
}

func(d *DatastoreSetup) InsertAccount(tx *gorm.DB, reqPayload dao.Account)(err error){
	reqPayloadForLog := reqPayload
	reqPayloadForLog.HashedPin = "*REDACTED*"
	d.Logger.Info(
		logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayloadForLog)}, nil, "START: InsertAccount Datastore",
	)

	if err = tx.Create(&reqPayload).Error; err != nil {
		d.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, err.Error(),
		)
		return
	}

	remark := "END: InsertAccount Datastore"
	d.Logger.Info(
		logrus.Fields{"response": fmt.Sprintf("%+v",  map[string]interface{}{})}, nil, remark,
	)

	return
}

func (d *DatastoreSetup) GetCustomer(tx *gorm.DB, reqPayload dao.Customer) (datastoreResponse dao.Customer, err error) {
	d.Logger.Info(
		logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayload)}, nil, "START: GetCustomer Datastore",
	)

	err = tx.Where("no_hp = ? OR nik = ? OR email = ?", reqPayload.NoHp, reqPayload.Nik, reqPayload.Email).Take(&datastoreResponse).Error
	if err != nil {
		d.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, err.Error(),
		)
		return
	}

	remark := "END: GetCustomer Datastore"
	d.Logger.Info(
		logrus.Fields{"response": fmt.Sprintf("%+v", datastoreResponse)}, nil, remark,
	)

	return
}

func (d *DatastoreSetup) CheckCustomerNotExist(tx *gorm.DB, reqPayload dao.Customer) (datastoreResponse bool, err error) {
	d.Logger.Info(
		logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayload)}, nil, "START: CheckCustomerNotExist Datastore",
	)
	isNotExist := false
	var customerModel dao.Customer
	err = tx.Where("no_hp = ? OR nik = ?", reqPayload.NoHp, reqPayload.Nik).Take(&customerModel).Error
	if err == nil {
		err = fmt.Errorf("customer is exist")
		d.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, err.Error(),
		)
		datastoreResponse = isNotExist
		return
	}
	if err == gorm.ErrRecordNotFound {
		err = nil
		isNotExist = true
		datastoreResponse = isNotExist
	}
	if err != nil{
		d.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, err.Error(),
		)
		isNotExist = true
		datastoreResponse = isNotExist
		return
	}

	remark := "END: CheckCustomerNotExist Datastore"
	d.Logger.Info(
		logrus.Fields{"response": fmt.Sprintf("CheckCustomerNotExist: %+v", datastoreResponse)}, nil, remark,
	)

	return
}


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


func (d *DatastoreSetup) GetAccountAndCif(tx *gorm.DB, noRekening string) (datastoreResponse dao.AccountWithCustomerQuery, err error) {
	reqPayloadForLog := map[string]interface{}{
		"noRekening": noRekening,
	}
	d.Logger.Info(
		logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayloadForLog)}, nil, "START: GetAccount Datastore",
	)

	err = tx.Table("account"). 
		Select("account.id, account.id_nasabah, account.no_rekening, account.hashed_pin, account.saldo, account.created_at, account.updated_at, account.is_deleted, " +
			"customer.id as customer_id, customer.nama, customer.nik, customer.no_hp, customer.email, customer.created_at as customer_created_at, customer.updated_at as customer_updated_at, customer.is_deleted as customer_is_deleted").
		Joins("inner join customer on account.id_nasabah = customer.id").
		Where("account.no_rekening = ?", noRekening).
		Take(&datastoreResponse).Error
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


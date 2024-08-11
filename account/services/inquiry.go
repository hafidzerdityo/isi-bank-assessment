package services

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hafidzresttemplate.com/dao"
)


func (s *ServiceSetup)GetSaldo(reqPayload dao.NoRekeningReq) (appResponse dao.SaldoRes, remark string, err error) {
	s.Logger.Info(
		logrus.Fields{"req_payload": fmt.Sprintf("%+v", appResponse)}, nil, "START: GetSaldo Service",
	)
	tx := s.Db.Begin()
	if tx.Error != nil {
		remark = "Error When Initializing DB"
		return
	}

	// check if customer exist
	customerData, err := s.Datastore.GetAccount(s.Db, reqPayload.NoRekening)
	if err == gorm.ErrRecordNotFound{
		tx.Rollback()
		err = fmt.Errorf("account not exist error")
		remark = "no rekening tidak terdaftar" 
		s.Logger.Error(
			logrus.Fields{"validation_error": err.Error()}, nil, remark,
		)
		return
	}
	if err != nil {
		tx.Rollback()
		remark = "database get data error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}
	// if customerData.ID == 0{
	// 	tx.Rollback()
	// 	err = fmt.Errorf("customer not exist error")
	// 	remark = "No Rekening Tidak Dikenali" 
	// 	s.Logger.Error(
	// 		logrus.Fields{"validation_error": err.Error()}, nil, remark,
	// 	)
	// 	return
	// }

	// Get Saldo
	appResponse.Saldo = &customerData.Saldo

	tx.Commit()

	remark = "END: GetSaldo Service"
	s.Logger.Info(
		logrus.Fields{"response": fmt.Sprintf("%+v", appResponse)}, nil, remark,
	)
	return
}


func (s *ServiceSetup)GetMutasi(reqPayload dao.MutasiReq) (appResponse dao.MutasiRes, remark string, err error) {
	s.Logger.Info(
		logrus.Fields{"req_payload": fmt.Sprintf("%+v", appResponse)}, nil, "START: GetMutasi Service",
	)
	tx := s.Db.Begin()
	if tx.Error != nil {
		remark = "Error When Initializing DB"
		return
	}

	// check if customer exist
	customerData, err := s.Datastore.GetAccount(s.Db, reqPayload.NoRekening)
	if err == gorm.ErrRecordNotFound{
		tx.Rollback()
		err = fmt.Errorf("account not exist error")
		remark = "no rekening tidak terdaftar" 
		s.Logger.Error(
			logrus.Fields{"validation_error": err.Error()}, nil, remark,
		)
		return
	}
	if err != nil {
		tx.Rollback()
		remark = "database get data error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	
	getMutasiParamLimit := int(15)
	if reqPayload.Limit != nil{
		getMutasiParamLimit = *reqPayload.Limit
	}
	getMutasiParamOffset := int(1)
	if reqPayload.Page != nil{
		getMutasiParamOffset = *reqPayload.Page
	}

	// get mutasi
	mutasiData, countMutasi, err := s.Datastore.GetMutasi(s.Db, customerData.ID, getMutasiParamLimit, getMutasiParamOffset)
	if (len(mutasiData) == 0){
		tx.Rollback()
		err = fmt.Errorf("mutasi data not exist error")
		remark = "mutasi tidak ditemukan" 
		s.Logger.Error(
			logrus.Fields{"validation_error": err.Error()}, nil, remark,
		)
		return
	}
	if err != nil {
		tx.Rollback()
		remark = "database select data error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	// Get List Mutasi
	for _, val := range mutasiData {
		derNominalIn := float64(0)
		if val.NominalIn != nil{
			derNominalIn = *val.NominalIn
		}
		derNominalOut := float64(0)
		if val.NominalOut != nil{
			derNominalOut = *val.NominalOut
		}
		mutasiRes := dao.MutasiData{
			Waktu:          val.Waktu,
			JenisTransaksi: val.JenisTransaksi,
			IdJurnal: val.IdJurnal,
			NominalIn:        derNominalIn,
			NominalOut:       derNominalOut,
		}
		appResponse.ListData = append(appResponse.ListData, mutasiRes)
	}

	appResponse.TotalData = countMutasi

	tx.Commit()

	remark = "END: GetMutasi Service"
	s.Logger.Info(
		logrus.Fields{"response": fmt.Sprintf("%+v", appResponse)}, nil, remark,
	)
	return
}


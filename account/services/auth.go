package services

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hafidzresttemplate.com/dao"
	"hafidzresttemplate.com/pkg/utils"
)


func (s *ServiceSetup)CheckAccountAndPin(reqPayload dao.CheckAccountAndPinReq) (isExist bool, remark string, err error) {
	reqPayloadForLog := reqPayload
	reqPayloadForLog.Pin = "*REDACTED*"
	s.Logger.Info(
		logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayloadForLog)}, nil, "START: CheckAccountAndPin Service",
	)
	tx := s.Db.Begin()
	if tx.Error != nil {
		remark = "Error When Initializing DB"
		return
	}

	// check if account exist
	accountData, err := s.Datastore.GetAccount(s.Db, reqPayloadForLog.NoRekening)
	if err != nil {
		tx.Rollback()
		remark = "Data Get Error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}
	if accountData.ID == 0{
		tx.Rollback()
		err = fmt.Errorf("account not exist error")
		remark = "no_rekening belum terdaftar" 
		s.Logger.Error(
			logrus.Fields{"validation_error": err.Error()}, nil, remark,
		)
		return
	}

	// check if pin correct
	err = utils.VerifyPassword(reqPayload.Pin, accountData.HashedPin)
	if err != nil{
		tx.Rollback()
		err = fmt.Errorf("wrong pin error")
		remark = "pin salah" 
		s.Logger.Error(
			logrus.Fields{"validation_error": err.Error()}, nil, remark,
		)
		return
	}

	tx.Commit()
	isExist = true

	remark = "END: CheckAccountAndPin Service"
	s.Logger.Info(
		logrus.Fields{"response": "Account Exist and Pin is correct"}, nil, remark,
	)
	return
}


func (s *ServiceSetup)AccountLogin(reqPayload dao.AccountLoginReq) (appResponse dao.AccountLoginRes, remark string, err error) {
	reqPayloadForLog := reqPayload
	reqPayloadForLog.Pin = "*REDACTED*"
	s.Logger.Info(
		logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayloadForLog)}, nil, "START: AccountLogin Service",
	)
	tx := s.Db.Begin()
	if tx.Error != nil {
		remark = "Error When Initializing DB"
		return
	}

	if !utils.IsDigit(reqPayload.Pin){
		tx.Rollback()
		err = fmt.Errorf("pin validation error")
		remark = "pin must be a string of digit" 
		s.Logger.Error(
			logrus.Fields{"validation_error": err.Error()}, nil, remark,
		)
		return
	}

	// check if no_rekening exist
	getAccountData := reqPayload.NoRekening
	loginData, err := s.Datastore.GetAccount(s.Db, getAccountData)
	if err == gorm.ErrRecordNotFound{
		tx.Rollback()
		err = fmt.Errorf("no_rekening not found error")
		remark = "no_rekening tidak ditemukan" 
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

	if reqPayload.Pin == ""{
		tx.Rollback()
		err = fmt.Errorf("pin empty error")
		remark = "silahkan input pin" 
		s.Logger.Error(
			logrus.Fields{"validation_error": err.Error()}, nil, remark,
		)
		return
	}
	
	// check if pin correct
	if reqPayload.Pin != ""{
		err = utils.VerifyPassword(reqPayload.Pin, loginData.HashedPin)
		if err != nil{
			tx.Rollback()
			err = fmt.Errorf("wrong pin error")
			remark = "pin salah" 
			s.Logger.Error(
				logrus.Fields{"validation_error": err.Error()}, nil, remark,
			)
			return
		}
	}

	// generate token
	var JWTFieldParam dao.JWTField
	JWTFieldParam.NoRekening = loginData.NoRekening
	tokenString, err := utils.CreateJWTToken(JWTFieldParam)
	if err != nil{
		tx.Rollback()
		err = fmt.Errorf("wrong password error")
		remark = "password salah" 
		s.Logger.Error(
			logrus.Fields{"validation_error": err.Error()}, nil, remark,
		)
		return
	}

	appResponse.AccessToken = tokenString
	appResponse.TokenType = "bearer"

	tx.Commit()

	remark = "END: AccountLogin Service"
	s.Logger.Info(
		logrus.Fields{"response": fmt.Sprintf("%+v", appResponse)}, nil, remark,
	)
	return
}


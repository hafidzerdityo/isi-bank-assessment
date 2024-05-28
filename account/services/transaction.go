package services

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	customerrors "hafidzresttemplate.com/customErrors"
	"hafidzresttemplate.com/dao"
	"hafidzresttemplate.com/pkg/utils"
)



func (s *ServiceSetup)CreateTabung(reqPayload dao.CreateTabungTarikUpdate) (appResponse dao.SaldoRes, remark string, err error) {
	s.Logger.Info(
		logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayload)}, nil, "START: CreateTabung Service",
	)
	tx := s.Db.Begin()
	if tx.Error != nil {
		remark = "Error When Initializing DB"
		return
	}

	defer s.TransactionStatusHandler(tx, &err)

	// check if account exist
	accountData, err := s.Datastore.GetAccount(s.Db, reqPayload.NoRekening)
	if err == gorm.ErrRecordNotFound{
		err = customerrors.ErrAccountNotFound
		remark = "no rekening tidak ditemukan" 
		s.Logger.Error(
			logrus.Fields{"validation_error": err.Error()}, nil, remark,
		)
		return
	}
	if err != nil {
		remark = "database get data error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	if reqPayload.Pin == ""{
		err = fmt.Errorf("pin empty error")
		remark = "silahkan input pin" 
		s.Logger.Error(
			logrus.Fields{"validation_error": err.Error()}, nil, remark,
		)
		return
	}
	
	// check if pin correct
	if reqPayload.Pin != ""{
		err = utils.VerifyPassword(reqPayload.Pin, accountData.HashedPin)
		if err != nil{
			err = customerrors.ErrWrongPassword
			remark = "pin salah" 
			s.Logger.Error(
				logrus.Fields{"validation_error": err.Error()}, nil, remark,
			)
			return
		}
	}


	// update account saldo (increase)
	saldo, err := s.Datastore.UpdateSaldo(tx, reqPayload)
	if err != nil {
		remark = "database update data error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	// create catatan transfer
	var insertCatatanParam dao.Transaction
	insertCatatanParam.IdRekening = accountData.ID
	insertCatatanParam.JenisTransaksi = "C"
	insertCatatanParam.Nominal = reqPayload.Nominal
	insertCatatanParam.Waktu = time.Now()
	insertCatatanParam.NomorRekeningTujuan = nil
	err = s.Datastore.InsertCatatan(tx, insertCatatanParam)
	if err != nil {
		remark = "database insert data error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	var PubMessageParam dao.PubStruct
	PubMessageParam.TanggalTransaksi = insertCatatanParam.Waktu
	PubMessageParam.NoRekening= accountData.NoRekening
	PubMessageParam.JenisTransaksi= "C"
	PubMessageParam.Nominal = insertCatatanParam.Nominal

	// send message to the event stream
	err = s.EventPub.PublishJournal(PubMessageParam)
	if err != nil {
		remark = "Failed to send message to the event stream"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	appResponse.Saldo = &saldo

	remark = "END: CreateTabung Service"
	s.Logger.Info(
		logrus.Fields{"response": fmt.Sprintf("%+v", appResponse)}, nil, remark,
	)
	return
}


func (s *ServiceSetup)CreateTarik(reqPayload dao.CreateTabungTarikUpdate) (appResponse dao.SaldoRes, remark string, err error) {
	s.Logger.Info(
		logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayload)}, nil, "START: CreateTarik Service",
	)
	tx := s.Db.Begin()
	if tx.Error != nil {
		remark = "Error When Initializing DB"
		return
	}

	defer s.TransactionStatusHandler(tx, &err)

	// check if account exist
	accountData, err := s.Datastore.GetAccount(s.Db, reqPayload.NoRekening)
	if err == gorm.ErrRecordNotFound{
		err = customerrors.ErrAccountNotFound
		remark = "no rekening tidak ditemukan" 
		s.Logger.Error(
			logrus.Fields{"validation_error": err.Error()}, nil, remark,
		)
		return
	}
	if err != nil {
		remark = "database get data error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	if reqPayload.Pin == ""{
		err = customerrors.ErrWrongPassword
		remark = "silahkan input pin" 
		s.Logger.Error(
			logrus.Fields{"validation_error": err.Error()}, nil, remark,
		)
		return
	}
	
	// check if pin correct
	if reqPayload.Pin != ""{
		err = utils.VerifyPassword(reqPayload.Pin, accountData.HashedPin)
		if err != nil{
			err = customerrors.ErrWrongPassword
			remark = "pin salah" 
			s.Logger.Error(
				logrus.Fields{"validation_error": err.Error()}, nil, remark,
			)
			return
		}
	}

	if accountData.Saldo < reqPayload.Nominal{
		err = customerrors.ErrInsufficientBalance
		remark = "maaf, saldo tidak cukup" 
		s.Logger.Error(
			logrus.Fields{"validation_error": err.Error()}, nil, remark,
		)
		return
	}

	// update account saldo (decrease)
	var decreaseSaldoParam dao.CreateTabungTarikUpdate
	decreaseSaldoParam.Nominal = -1 * reqPayload.Nominal
	decreaseSaldoParam.NoRekening = reqPayload.NoRekening
	saldo, err := s.Datastore.UpdateSaldo(tx, decreaseSaldoParam)
	if err != nil {
		remark = "database update data error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	// create catatan transfer
	var insertCatatanParam dao.Transaction
	insertCatatanParam.IdRekening = accountData.ID
	insertCatatanParam.JenisTransaksi = "D"
	insertCatatanParam.Nominal = reqPayload.Nominal
	insertCatatanParam.Waktu = time.Now()
	insertCatatanParam.NomorRekeningTujuan = nil
	err = s.Datastore.InsertCatatan(tx, insertCatatanParam)
	if err != nil {
		remark = "database insert data error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	var PubMessageParam dao.PubStruct
	PubMessageParam.TanggalTransaksi = insertCatatanParam.Waktu
	PubMessageParam.NoRekening = accountData.NoRekening
	PubMessageParam.Nominal = insertCatatanParam.Nominal
	PubMessageParam.JenisTransaksi = "D"

	// send message to the event stream
	err = s.EventPub.PublishJournal(PubMessageParam)
	if err != nil {
		remark = "failed to send message to the event stream"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	appResponse.Saldo = &saldo

	remark = "END: CreateTarik Service"
	s.Logger.Info(
		logrus.Fields{"response": fmt.Sprintf("%+v", appResponse)}, nil, remark,
	)
	return
}


func (s *ServiceSetup)CreateTransfer(reqPayload dao.CreateTransferUpdate) (appResponse dao.SaldoRes, remark string, err error) {
	s.Logger.Info(
		logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayload)}, nil, "START: CreateTransfer Service",
	)
	tx := s.Db.Begin()
	if tx.Error != nil {
		remark = "Error When Initializing DB"
		return
	}

	defer s.TransactionStatusHandler(tx, &err)

	if reqPayload.NoRekeningAsal == reqPayload.NoRekeningTujuan{
		err = fmt.Errorf("same sender and receiver account number error")
		remark = "momor rekening pengirim dan penerima tidak boleh sama"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	// check if sender account exist
	accountSenderData, err := s.Datastore.GetAccount(s.Db, reqPayload.NoRekeningAsal)
	if err == gorm.ErrRecordNotFound{
		err = fmt.Errorf("account not exist error")
		remark = "no rekening asal tidak ditemukan" 
		s.Logger.Error(
			logrus.Fields{"validation_error": err.Error()}, nil, remark,
		)
		return
	}
	if err != nil {
		remark = "database get data error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	// check if benefciary account exist
	accountReceiverData, err := s.Datastore.GetAccount(s.Db, reqPayload.NoRekeningTujuan)
	if err == gorm.ErrRecordNotFound{
		err = fmt.Errorf("account not exist error")
		remark = "no rekening tujuan tidak ditemukan" 
		s.Logger.Error(
			logrus.Fields{"validation_error": err.Error()}, nil, remark,
		)
		return
	}
	if err != nil {
		remark = "database get data error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	if accountSenderData.Saldo < reqPayload.Nominal{
		err = fmt.Errorf("insufficient balance error")
		remark = "maaf, saldo tidak cukup" 
		s.Logger.Error(
			logrus.Fields{"validation_error": err.Error()}, nil, remark,
		)
		return
	}

	// update sender account saldo (decrease)
	var decreaseSaldoParam dao.CreateTabungTarikUpdate 
	decreaseSaldoParam.NoRekening = reqPayload.NoRekeningAsal
	decreaseSaldoParam.Nominal = -1 * reqPayload.Nominal

	saldo, err := s.Datastore.UpdateSaldo(tx, decreaseSaldoParam)
	if err != nil {
		remark = "database update data error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}
	
	// update account saldo (increase)
	var increaseSaldoParam dao.CreateTabungTarikUpdate 
	increaseSaldoParam.NoRekening = reqPayload.NoRekeningTujuan
	increaseSaldoParam.Nominal = reqPayload.Nominal
	_, err = s.Datastore.UpdateSaldo(tx, increaseSaldoParam)
	if err != nil {
		remark = "database update data error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	catatanTime := time.Now()

	// create catatan transfer sender
	var insertCatatanParam dao.Transaction
	insertCatatanParam.IdRekening = accountSenderData.ID
	insertCatatanParam.JenisTransaksi = "T"
	insertCatatanParam.Nominal = reqPayload.Nominal
	insertCatatanParam.Waktu = catatanTime
	insertCatatanParam.NomorRekeningTujuan = &accountReceiverData.NoRekening
	err = s.Datastore.InsertCatatan(tx, insertCatatanParam)
	if err != nil {
		remark = "database insert data error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	// create catatan transfer receiver
	insertCatatanParam.IdRekening = accountReceiverData.ID
	insertCatatanParam.JenisTransaksi = "T"
	insertCatatanParam.Nominal = reqPayload.Nominal
	insertCatatanParam.Waktu = catatanTime
	insertCatatanParam.NomorRekeningTujuan = nil
	err = s.Datastore.InsertCatatan(tx, insertCatatanParam)
	if err != nil {
		remark = "database insert data error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	// var PubMessageParamSender dao.PubStruct
	// PubMessageParamSender.TanggalTransaksi = insertCatatanParam.Waktu
	// PubMessageParamSender.NoRekening = accountSenderData.NoRekening
	// PubMessageParamSender.Nominal = insertCatatanParam.Nominal
	// PubMessageParamSender.JenisTransaksi = "TC"


	// // send message to the event stream
	// err = s.EventPub.PublishJournal(PubMessageParamSender)
	// if err != nil {
	// 	remark = "Failed to send message to the event stream"
	// 	s.Logger.Error(
	// 		logrus.Fields{"error": err.Error()}, nil, remark,
	// 	)
	// 	return
	// }
	// var PubMessageParamReceiver dao.PubStruct
	// PubMessageParamReceiver.TanggalTransaksi = insertCatatanParam.Waktu
	// PubMessageParamReceiver.NoRekening = accountReceiverData.NoRekening
	// PubMessageParamReceiver.Nominal = insertCatatanParam.Nominal
	// PubMessageParamSender.JenisTransaksi = "TR"
	// // send message to the event stream
	// err = s.EventPub.PublishJournal(PubMessageParamReceiver)
	// if err != nil {
	// 	remark = "Failed to send message to the event stream"
	// 	s.Logger.Error(
	// 		logrus.Fields{"error": err.Error()}, nil, remark,
	// 	)
	// 	return
	// }

	appResponse.Saldo = &saldo

	remark = "END: CreateTransfer Service"
	s.Logger.Info(
		logrus.Fields{"response": fmt.Sprintf("%+v", appResponse)}, nil, remark,
	)
	return
}

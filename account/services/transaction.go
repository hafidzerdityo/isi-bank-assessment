package services

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hafidzresttemplate.com/dao"
	"hafidzresttemplate.com/pkg/utils"
)



func (s *ServiceSetup)CreateTabung(reqPayload dao.CreateTabungTarikUpdate) (appResponse dao.SaldoRes, remark string, err error) {
	ctx := context.Background()
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
		err = utils.ErrAccountNotFound
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
	idJurnal := utils.GenerateNumericUUID(12)
	var insertCatatanParam dao.Transaction
	insertCatatanParam.IdJurnal = idJurnal
	insertCatatanParam.IdRekening = accountData.ID
	insertCatatanParam.JenisTransaksi = "C"
	insertCatatanParam.NominalIn = &reqPayload.Nominal
	wibTime, _ := time.LoadLocation("Asia/Jakarta")
	insertCatatanParam.Waktu = time.Now().In(wibTime)
	err = s.Datastore.InsertCatatan(tx, insertCatatanParam)
	if err != nil {
		remark = "database insert data error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	var PubMessageParam dao.PubStruct
	PubMessageParam.Waktu = insertCatatanParam.Waktu
	PubMessageParam.NoRekening= accountData.NoRekening
	PubMessageParam.JenisTransaksi= "C"
	PubMessageParam.NominalIn = *insertCatatanParam.NominalIn
	PubMessageParam.IdJurnal = idJurnal

	// send message to the event stream
	err = s.EventPub.PublishJournal(ctx, PubMessageParam)
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
	ctx := context.Background()
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
		err = utils.ErrAccountNotFound
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

	if accountData.Saldo < reqPayload.Nominal{
		err = utils.ErrInsufficientBalance
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
	idJurnal := utils.GenerateNumericUUID(12)
	insertCatatanParam.IdJurnal = idJurnal
	insertCatatanParam.IdRekening = accountData.ID
	insertCatatanParam.JenisTransaksi = "D"
	insertCatatanParam.NominalOut= &reqPayload.Nominal
	insertCatatanParam.Waktu = time.Now()
	err = s.Datastore.InsertCatatan(tx, insertCatatanParam)
	if err != nil {
		remark = "database insert data error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	var PubMessageParam dao.PubStruct
	PubMessageParam.Waktu = insertCatatanParam.Waktu
	PubMessageParam.NoRekening = accountData.NoRekening
	PubMessageParam.NominalOut = *insertCatatanParam.NominalOut
	PubMessageParam.JenisTransaksi = "D"
	PubMessageParam.IdJurnal = idJurnal

	// send message to the event stream
	err = s.EventPub.PublishJournal(ctx, PubMessageParam)
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
	ctx := context.Background()
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
	idJurnal := utils.GenerateNumericUUID(12)
	var insertCatatanParamSender dao.Transaction
	insertCatatanParamSender.IdJurnal = idJurnal
	insertCatatanParamSender.IdRekening = accountSenderData.ID
	insertCatatanParamSender.JenisTransaksi = "T"
	insertCatatanParamSender.NominalOut = &reqPayload.Nominal
	insertCatatanParamSender.NominalIn = nil
	insertCatatanParamSender.Waktu = catatanTime
	err = s.Datastore.InsertCatatan(tx, insertCatatanParamSender)
	if err != nil {
		remark = "database insert data error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	var insertCatatanParamReceiver dao.Transaction

	// create catatan transfer receiver
	insertCatatanParamReceiver.IdJurnal = idJurnal
	insertCatatanParamReceiver.IdRekening = accountReceiverData.ID
	insertCatatanParamReceiver.JenisTransaksi = "T"
	insertCatatanParamReceiver.NominalIn = &reqPayload.Nominal
	insertCatatanParamReceiver.NominalOut = nil
	insertCatatanParamReceiver.Waktu = catatanTime
	err = s.Datastore.InsertCatatan(tx, insertCatatanParamReceiver)
	if err != nil {
		remark = "database insert data error"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	var PubMessageParamSender dao.PubStruct
	PubMessageParamSender.Waktu = catatanTime
	PubMessageParamSender.NoRekening = accountSenderData.NoRekening
	PubMessageParamSender.NominalOut = reqPayload.Nominal
	PubMessageParamSender.JenisTransaksi = "T"
	PubMessageParamSender.IdJurnal = idJurnal


	// send message to the event stream
	err = s.EventPub.PublishJournal(ctx, PubMessageParamSender)
	if err != nil {
		remark = "Failed to send message to the event stream"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}
	var PubMessageParamReceiver dao.PubStruct
	PubMessageParamReceiver.Waktu = catatanTime
	PubMessageParamReceiver.NoRekening = accountReceiverData.NoRekening
	PubMessageParamReceiver.NominalIn = reqPayload.Nominal
	PubMessageParamSender.JenisTransaksi = "T"
	PubMessageParamSender.IdJurnal = idJurnal
	// send message to the event stream
	
	err = s.EventPub.PublishJournal(ctx, PubMessageParamReceiver)
	if err != nil {
		remark = "Failed to send message to the event stream"
		s.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	appResponse.Saldo = &saldo

	remark = "END: CreateTransfer Service"
	s.Logger.Info(
		logrus.Fields{"response": fmt.Sprintf("%+v", appResponse)}, nil, remark,
	)
	return
}

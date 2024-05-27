package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"hafidzresttemplate.com/customErrors"
	"hafidzresttemplate.com/dao"
	"hafidzresttemplate.com/pkg/utils"
)

func (a *ApiSetup) CreateTabung(c *fiber.Ctx) error {
	var reqPayload dao.CreateTabungTarikReq
    if err := c.BodyParser(&reqPayload); err != nil {
        a.Logger.Error(
            logrus.Fields{"error": err.Error()}, nil, err.Error(),
        )
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "resp_msg" : err.Error(),
            "resp_data" : dao.SaldoRes{
                Saldo: nil,
            },
        })
    }
	
	// Validate request payload
	if errMsg, err := utils.ValidateStruct(&reqPayload); err != nil {
		a.Logger.Error(
			logrus.Fields{"error": errMsg}, nil, errMsg,
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"resp_msg" : errMsg,
			"resp_data" : dao.SaldoRes{
				Saldo: nil,
			},
		})
	}

	authPayload := c.Locals("auth_payload").(map[string]interface{})
	noRekening := authPayload["no_rekening"].(string)
	pin := authPayload["pin"].(string)

    a.Logger.Info(
        logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayload)}, nil, "START: CreateTabung API",
    )

	tabungUpdateParam := dao.CreateTabungTarikUpdate{
		Nominal: reqPayload.Nominal,
		NoRekening: noRekening,
		Pin: pin,
	}
    data, remark, err := a.Services.CreateTabung(tabungUpdateParam)
    if err != nil {
        a.Logger.Error(
            logrus.Fields{"error": err.Error()}, nil, remark,
        )

		statusCode := fiber.StatusInternalServerError
		if err == customErrors.ErrInsufficientBalance{
			statusCode = fiber.StatusBadRequest
		}
		if err == customErrors.ErrAccountNotFound{
			statusCode = fiber.StatusNotFound
		}
		if err == customErrors.ErrWrongPassword{
			statusCode = fiber.StatusUnauthorized
		}
        return c.Status(statusCode).JSON(fiber.Map{
            "resp_msg" : remark,
            "resp_data" : data,
        })
    }

	response := map[string]interface{}{
        "resp_msg" : "Tabung Succeed",
        "resp_data" : data,
	}

	remark = "END: CreateTabung API"
    a.Logger.Info(
        logrus.Fields{"response": fmt.Sprintf("%+v", response)}, nil, remark,
    )
    return c.JSON(response)
}

func (a *ApiSetup) CreateTarik(c *fiber.Ctx) error {
	var reqPayload dao.CreateTabungTarikReq
    if err := c.BodyParser(&reqPayload); err != nil {
        a.Logger.Error(
            logrus.Fields{"error": err.Error()}, nil, err.Error(),
        )
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "resp_msg" : err.Error(),
            "resp_data" : dao.SaldoRes{
                Saldo: nil,
            },
        })
    }

	// Validate request payload
	if errMsg, err := utils.ValidateStruct(&reqPayload); err != nil {
		a.Logger.Error(
			logrus.Fields{"error": errMsg}, nil, errMsg,
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"resp_msg" : errMsg,
			"resp_data" : dao.SaldoRes{
				Saldo: nil,
			},
		})
	}

	authPayload := c.Locals("auth_payload").(map[string]interface{})
	noRekening := authPayload["no_rekening"].(string)
	pin := authPayload["pin"].(string)

    a.Logger.Info(
        logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayload)}, nil, "START: CreateTabung API",
    )

	tabungUpdateParam := dao.CreateTabungTarikUpdate{
		Nominal: reqPayload.Nominal,
		NoRekening: noRekening,
		Pin: pin,
	}


    data, remark, err := a.Services.CreateTarik(tabungUpdateParam)
    if err != nil {
        a.Logger.Error(
            logrus.Fields{"error": err.Error()}, nil, remark,
        )

		statusCode := fiber.StatusInternalServerError
		if err == customErrors.ErrInsufficientBalance{
			statusCode = fiber.StatusBadRequest
		}
		if err == customErrors.ErrAccountNotFound{
			statusCode = fiber.StatusNotFound
		}
		if err == customErrors.ErrWrongPassword{
			statusCode = fiber.StatusUnauthorized
		}
        return c.Status(statusCode).JSON(fiber.Map{
            "resp_msg" : remark,
            "resp_data" : data,
        })

    }

	response := map[string]interface{}{
        "resp_msg" : "Tarik Succeed",
        "resp_data" : data,
	}

	remark = "END: CreateTarik API"
    a.Logger.Info(
        logrus.Fields{"response": fmt.Sprintf("%+v", response)}, nil, remark,
    )
    return c.JSON(response)
}

func (a *ApiSetup) CreateTransfer(c *fiber.Ctx) error {
	var reqPayload dao.CreateTransferReq
    if err := c.BodyParser(&reqPayload); err != nil {
        a.Logger.Error(
            logrus.Fields{"error": err.Error()}, nil, err.Error(),
        )
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "resp_msg" : err.Error(),
            "resp_data" : dao.SaldoRes{
                Saldo: nil,
            },
        })
    }

	decodedJWT := c.Locals("decodedJWT").(map[string]interface{})
	noRekeningAsal := decodedJWT["no_rekening"].(string)

    a.Logger.Info(
        logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayload)}, nil, "START: CreateTransfer API",
    )

    // Validate request payload
    if errMsg, err := utils.ValidateStruct(&reqPayload); err != nil {
        a.Logger.Error(
            logrus.Fields{"error": errMsg}, nil, errMsg,
        )
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "resp_msg" : errMsg,
            "resp_data" : dao.SaldoRes{
                Saldo: nil,
            },
        })
    }

	tabungUpdateParam := dao.CreateTransferUpdate{
		Nominal: reqPayload.Nominal,
		NoRekeningTujuan: reqPayload.NoRekeningTujuan,
		NoRekeningAsal: noRekeningAsal,
	}

    data, remark, err := a.Services.CreateTransfer(tabungUpdateParam)
    if err != nil {
        a.Logger.Error(
            logrus.Fields{"error": err.Error()}, nil, remark,
        )
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "resp_msg" : remark,
            "resp_data" : data,
        })
    }

	response := map[string]interface{}{
        "resp_msg" : "Transfer Succeed",
        "resp_data" : data,
	}

	remark = "END: CreateTransfer API"
    a.Logger.Info(
        logrus.Fields{"response": fmt.Sprintf("%+v", response)}, nil, remark,
    )
    return c.JSON(response)
}

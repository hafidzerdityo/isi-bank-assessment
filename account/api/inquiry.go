package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"hafidzresttemplate.com/dao"
	"hafidzresttemplate.com/pkg/utils"
)

func (a *ApiSetup) GetSaldo(c *fiber.Ctx) error {
	var reqPayload dao.NoRekeningReq
	decodedJWT := c.Locals("decodedJWT").(map[string]interface{})
	noRekening := decodedJWT["no_rekening"].(string)
	reqPayload.NoRekening = noRekening

    a.Logger.Info(
        logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayload)}, nil, "START: GetSaldo API",
    )

    // Validate request payload
    if errMsg, err := utils.ValidateStruct(&reqPayload); err != nil {
        a.Logger.Error(
            logrus.Fields{"error": errMsg}, nil, errMsg,
        )
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "resp_msg" : errMsg,
            "resp_data" : dao.CreateCustRes{
                NoRekening: nil,
            },
        })
    }

    data, remark, err := a.Services.GetSaldo(reqPayload)
    if err != nil {
        a.Logger.Error(
            logrus.Fields{"error": err.Error()}, nil, remark,
        )
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "resp_msg" : remark,
            "resp_data" : data,
        })
    }

	response := map[string]interface{}{
        "resp_msg" : "GetSaldo Succeed",
        "resp_data" : data,
	}

	remark = "END: GetSaldo API"
    a.Logger.Info(
        logrus.Fields{"response": fmt.Sprintf("%+v", response)}, nil, remark,
    )
    return c.JSON(response)
}

func (a *ApiSetup) GetMutasi(c *fiber.Ctx) error {
	var reqPayload dao.NoRekeningReq
	decodedJWT := c.Locals("decodedJWT").(map[string]interface{})
	noRekening := decodedJWT["no_rekening"].(string)
	reqPayload.NoRekening = noRekening

    a.Logger.Info(
        logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayload)}, nil, "START: GetMutasi API",
    )

    // Validate request payload
    if errMsg, err := utils.ValidateStruct(&reqPayload); err != nil {
        a.Logger.Error(
            logrus.Fields{"error": errMsg}, nil, errMsg,
        )
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "resp_msg" : errMsg,
            "resp_data" : dao.CreateCustRes{
                NoRekening: nil,
            },
        })
    }

    data, remark, err := a.Services.GetMutasi(reqPayload)
    if err != nil {
        a.Logger.Error(
            logrus.Fields{"error": err.Error()}, nil, remark,
        )
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "resp_msg" : remark,
            "resp_data" : data,
        })
    }

	response := map[string]interface{}{
        "resp_msg" : "GetMutasi Succeed",
        "resp_data" : data,
	}

	remark = "END: GetMutasi API"
    a.Logger.Info(
        logrus.Fields{"response": fmt.Sprintf("%+v", response)}, nil, remark,
    )
    return c.JSON(response)
}

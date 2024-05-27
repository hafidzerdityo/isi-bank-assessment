package api

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"hafidzresttemplate.com/dao"
)


func (a *ApiSetup) AccountLogin(c *fiber.Ctx) error {

	authorization := c.Get("Authorization")
	if authorization == "" {
	  a.Logger.Error(logrus.Fields{"error": "Missing Authorization header"}, "Missing Authorization header in login request")
	  return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"resp_msg": "Authorization header is required",
		"resp_data": dao.CreateCustRes{
		  NoRekening: nil,
		},
	  })
	}


	if(len(authorization) == 0) || (!strings.Contains(authorization, ".")){
		err := fmt.Errorf("authorization length cannot be 0 and must contain dot to seperate")
        a.Logger.Error(
            logrus.Fields{"error": err.Error()}, nil,  err.Error(),
        )
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "resp_msg" :  err.Error(),
            "resp_data" : dao.CreateCustRes{
                NoRekening: nil,
            },
        })
	}

	auth_splitted := strings.Split(authorization, ".")

	var authPayload dao.AccountLoginReq
	authPayload.NoRekening = auth_splitted[0]
	authPayload.Pin = auth_splitted[1]

	reqPayloadForLog := authPayload
	reqPayloadForLog.Pin = "*REDACTED*"
    a.Logger.Info(
        logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayloadForLog)}, nil, "START: AccountLogin API",
    )

    data, remark, err := a.Services.AccountLogin(authPayload)
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
        "resp_msg" : "Login Succeed",
        "resp_data" : data,
	}

	remark = "END: AccountLogin API"
    a.Logger.Info(
        logrus.Fields{"response": fmt.Sprintf("%+v", response)}, nil, remark,
    )
    return c.JSON(response)
}

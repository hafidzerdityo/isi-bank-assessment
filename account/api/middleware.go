package api

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"hafidzresttemplate.com/dao"
	"hafidzresttemplate.com/pkg/utils"
)

func (a *ApiSetup) JwtDecode() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            err := fmt.Errorf("auth header format error")
            remark := "Authorization header is not in the expected format"
            a.Logger.Error(
                logrus.Fields{"error": err.Error()}, nil, remark,
            )
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "resp_data": nil,
                "resp_msg":  remark,
            })
        }
        token := parts[1]
		a.Logger.Info(logrus.Fields{"JWT_MIDDLEWARE_REQUEST": token}, nil, "START: JwtDecode Middleware",)


        // Check JWT Token
        isValid, remark, decodedJWT, err := utils.ValidateJWTToken(token)
        if err != nil || !isValid {
            a.Logger.Error(
                logrus.Fields{"error": err.Error()}, nil, remark,
            )
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "resp_data": nil,
                "resp_msg":  remark,
            })
        }

		a.Logger.Info(logrus.Fields{"JWT_MIDDLEWARE_RESPONSE": fmt.Sprintf("%+v", decodedJWT)}, nil, "END: JwtDecode Middleware",)

        // Store decodedJWT in context locals for later use
        c.Locals("decodedJWT", decodedJWT)

        // Call next middleware or handler
        return c.Next()
    }
}


func (a *ApiSetup) PinDecode() fiber.Handler {
    return func(c *fiber.Ctx) error {
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
				"resp_data" : nil,
			})
		}
	
		auth_splitted := strings.Split(authorization, ".")
	
		authPayload := make(map[string]interface{})
		authPayload["no_rekening"] = auth_splitted[0]
		authPayload["pin"]  = auth_splitted[1]
	
		reqPayloadForLog := make(map[string]interface{})
		reqPayloadForLog["pin"] = "*REDACTED*"
		reqPayloadForLog["no_rekening"] = auth_splitted[0]
		a.Logger.Info(
			logrus.Fields{"req_payload": fmt.Sprintf("%+v", reqPayloadForLog)}, nil, "START: AccountLogin API",
		)
        // Store decodedJWT in context locals for later use
        c.Locals("auth_payload", authPayload)

        // Call next middleware or handler
        return c.Next()
    }
}

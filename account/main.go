package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"hafidzresttemplate.com/api"
	"hafidzresttemplate.com/logs"
	"hafidzresttemplate.com/startup"
)


func main() {

	loggerInit := logs.InitLog()
	dbInit, envInit, evenStreamInit, err := startup.Startup(loggerInit)
	if err != nil{
		remark := "Start Up Failed"
		loggerInit.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}
	app := api.InitApi(loggerInit, dbInit, evenStreamInit)
	app.Listen(fmt.Sprintf("%v:%v", envInit.Host, envInit.Port))
	loggerInit.Info("Service Started!")
}

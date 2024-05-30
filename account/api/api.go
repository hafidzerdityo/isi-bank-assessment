package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hafidzresttemplate.com/datastore"
	"hafidzresttemplate.com/services"
	"hafidzresttemplate.com/startup"
)

type ApiSetup struct {
    Logger *logrus.Logger
	Services *services.ServiceSetup
}

func NewApiSetup(loggerInit *logrus.Logger, db *gorm.DB, EventPubInit startup.EventStreamConfig)(apiSet ApiSetup) {
	apiSet = ApiSetup{
		Logger: loggerInit,
		Services: &services.ServiceSetup{
			Logger: loggerInit,
			Db: db,
			Datastore: &datastore.DatastoreSetup{
				Logger: loggerInit,
			},
		EventPub: EventPubInit,
		},
		
	}
    return 
}

func InitApi(loggerInit *logrus.Logger, dbInit *gorm.DB, EventPubInit startup.EventStreamConfig)(app *fiber.App) {
	app = fiber.New()
	app.Use(logger.New())
	app.Use(recover.New()) // Enable recover middleware

	apiSetup := NewApiSetup(loggerInit, dbInit, EventPubInit)
	apiSetup.Logger.Info("Setting up api routes...")

	api := app.Group("/api")
	v1 := api.Group("/v1")

	trx := v1.Group("/transaction")

	trx.Use(apiSetup.PinDecode())

	trx.Post("/tabung", apiSetup.CreateTabung)

	
	return 
}
package startup

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type envValue struct{
	Host         string
	Port         string
	DbUser      string
	DbName      string
	DbPassword  string
	DbHost      string
	DbPort      string
	Broker string
	Topic  string
}


func Startup(loggerInit *logrus.Logger)(dbInit *gorm.DB, envInit envValue, EventStreamInit EventStreamConfig,  err error){
	loggerInit.Info("Startup...")
	godotenv.Load("service.env")
	godotenv.Load("config.env")
	godotenv.Load("redis_stream_config.env")

	envInit = envValue{
		Host		: os.Getenv("SERVICE_HOST"),
		Port        : os.Getenv("SERVICE_PORT"),
		DbUser      : os.Getenv("POSTGRES_USER"),
		DbName      : os.Getenv("POSTGRES_DB"),
		DbPassword  : os.Getenv("POSTGRES_PASSWORD"),
		DbHost      : os.Getenv("POSTGRES_HOST"),
		DbPort      : os.Getenv("POSTGRES_PORT"),
		Broker : os.Getenv("BROKER"),
		Topic  : os.Getenv("TOPIC"),
	}


	loggerInit.Info("Initialize Database...")
	dbInit, err = InitializeDB(
		envInit.DbUser,
		envInit.DbName,
		envInit.DbPassword,
		envInit.DbHost,
		envInit.DbPort,
	)
	if err != nil{
		remark := "Error when initializing Database"
		loggerInit.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}
	loggerInit.Info("Connected to the Database!")
	
	loggerInit.Info("Initialize Redis Stream Client...")
	redistClient, err := InitiateRedisStream(
		loggerInit,
		envInit.Broker,
		
	)
	if err != nil{
		remark := "error when initializing Redis Stream Client"
		loggerInit.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}
	loggerInit.Info("Connected to Redis!")
	
	EventStreamInit = EventStreamConfig{
		Logger: loggerInit,
		RedisClient: redistClient,
		Topic: envInit.Topic,
	}

	return
}


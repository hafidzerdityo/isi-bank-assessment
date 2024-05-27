package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"hafidzresttemplate.com/consumer"
	"hafidzresttemplate.com/dao"
	"hafidzresttemplate.com/logs"
	"hafidzresttemplate.com/startup"
)



func main() {
	loggerInit := logs.InitLog()
	godotenv.Load("service.env")
	godotenv.Load("config.env")
	godotenv.Load("redis_stream_config.env")

	dbUser := os.Getenv("POSTGRES_USER")
	dbName := os.Getenv("POSTGRES_DB")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")

	dbInit, err := dao.InitializeDB(
		dbUser,
		dbName,
		dbPassword,
		dbHost,
		dbPort,
	)
	if err != nil{
		remark := "Error when initializing Database"
		loggerInit.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}


	brokerRedis := os.Getenv("BROKER")
	topicRedis := os.Getenv("TOPIC")
	groupRedis := os.Getenv("GROUP")

	redisClient, err := startup.InitiateRedisStream(loggerInit, brokerRedis)
	if err != nil{
		remark := "error when initializing Redis Stream Client"
		loggerInit.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		return
	}

	loggerInit.Info("Connected to Redis!")

	consumer.InitConsumer(loggerInit, 
		redisClient, 
		dbInit, 
		topicRedis, 
		groupRedis)
}


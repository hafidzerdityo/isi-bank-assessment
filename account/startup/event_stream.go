package startup

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"hafidzresttemplate.com/dao"
)

type EventStreamConfig struct {
	Logger *logrus.Logger
    RedisClient *redis.Client
	Topic string
}


func (k *EventStreamConfig)PublishJournal(ctx context.Context, message dao.PubStruct) (err error) {
	k.Logger.Info(
		logrus.Fields{"pub_payload": fmt.Sprintf("%+v", message)}, nil, "START: PublishJournal Stream",
	)
	// Use XAdd to publish the message
	messageJson, err := json.Marshal(message)
	if err != nil {
		remark := "error marshaling message" 
		k.Logger.Error(
			logrus.Fields{"pub_message_error": err.Error()}, nil, remark,
		)
		return
	}
	_, err = k.RedisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: k.Topic,
		Values: map[string]string{
			"message": string(messageJson),
		},
	}).Result()
	if err != nil {
		remark := "error when publishing message" 
		k.Logger.Error(
			logrus.Fields{"pub_message_error": err.Error()}, nil, remark,
		)
		return
	}

	remark := "END: PublishJournal Stream"
	k.Logger.Info(
		logrus.Fields{"response": fmt.Sprintf("%+v", string(messageJson))}, nil, remark,
	)

    return
}


func InitiateRedisStream(loggerInit *logrus.Logger, broker string)(redisClient *redis.Client, err error) {
    redisClient = redis.NewClient(&redis.Options{
        Addr: "redis:6379", // Connect to the Redis container using its service name
        DB:   0,            // Use the default database (0)
    })

	ctx := context.Background()
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		remark := "Unable to connect to redis!" 
		loggerInit.Fatal(
			logrus.Fields{"pub_message_error": err.Error()}, nil, remark,
		)
		return
	}
	return
}

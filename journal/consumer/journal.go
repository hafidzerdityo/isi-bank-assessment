package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"hafidzresttemplate.com/dao"
)

func (a *ConsumerSetup) CreateJournalLoop(ctx context.Context, redisClient *redis.Client, topic string, group string) {
  for {
	entries, err := redisClient.XRead(ctx, &redis.XReadArgs{
		Streams:  []string{topic, "$"},
		Count:    2,
		Block:    0,
	}).Result()
	if err != nil {
		remark := "failed to read the stream"
		a.Logger.Error(
			logrus.Fields{"error": err.Error()}, nil, remark,
		)
		continue
	}

	for i := 0; i < len(entries[0].Messages); i++ {
		_ = entries[0].Messages[i].ID
		values := entries[0].Messages[i].Values
		a.Logger.Info(logrus.Fields{
			"values_message": fmt.Sprintf("%+v", values),
		},nil,"subscribed redis message")

		messageJson := []byte(values["message"].(string))

		var reqRedis dao.SubStruct
		err = json.Unmarshal(messageJson, &reqRedis)
		if err != nil {
			fmt.Println("Error unmarshalling data:", err)
			return
		}

		a.Logger.Info(logrus.Fields{
			"message":reqRedis,
		},nil,"subscribed redis message")

		
		_, remark, err := a.Services.CreateJournal(reqRedis)
		if err != nil {
			a.Logger.Error(
				logrus.Fields{"error": err.Error()}, nil, remark,
			)
			continue
		}
	}
 
}}

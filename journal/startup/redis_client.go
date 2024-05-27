package startup

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

func InitiateRedisStream(loggerInit *logrus.Logger, broker string) (redisClient *redis.Client, err error) {
    redisClient = redis.NewClient(&redis.Options{
        Addr: "redis:6379", // Connect to the Redis container using its service name
        DB:   0,            // Use the default database (0)
    })

	ctx := context.Background()
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		remark := "Unable to connect to redis!" 
		loggerInit.Fatal(
			logrus.Fields{"connection_redis_error": err.Error()}, nil, remark,
		)
		return
	}

	
	return
}

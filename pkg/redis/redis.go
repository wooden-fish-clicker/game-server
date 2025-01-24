package redis

import (
	"context"

	"game-server/pkg/logger"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client *redis.Client
}

func NewRedisClient(addr string, password string, db int) *Redis {
	rd := connectRedis(addr, password, db)
	return &Redis{rd}
}

func connectRedis(addr string, password string, db int) *redis.Client {
	rd := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db, // 使用默認的資料庫
	})

	_, err := rd.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatal("無法連接到Redis: ", err)
		return nil
	}
	return rd

}

func (rd *Redis) CloseRedis() {
	defer rd.Client.Close()
}

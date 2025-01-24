package main

import (
	"game-server/configs"
	"game-server/internal/endpoints/grpc_handlers"
	"game-server/internal/repository"
	"game-server/internal/server"
	"game-server/internal/services"
	"game-server/pkg/db"
	"game-server/pkg/logger"
	"game-server/pkg/redis"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
)

func init() {
	configs.Setup()
	logger.Setup()

}
func main() {

	app := fx.New(
		// 依賴注入
		fx.Provide(
			func() *mongo.Database {
				mongo := db.NewMongoDBClient(configs.C.MongoDB.ConnString)
				return mongo.Client.Database(configs.C.MongoDB.Name)
			},

			func() *redis.Redis {
				return redis.NewRedisClient(configs.C.Redis.Addr, configs.C.Redis.Password, configs.C.Redis.DB)
			},

			repository.NewUserRepository,
			repository.NewUserCacheRepository,

			services.NewAttackService,
			services.NewClickService,

			grpc_handlers.NewGameServer,
		),

		// 啟動
		fx.Invoke(server.StartServer),
	)

	app.Run()

}

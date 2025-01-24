package server

import (
	"context"
	"net"

	"game-server/configs"
	"game-server/grpc_proto/game_server"
	"game-server/internal/endpoints/grpc_handlers"
	"game-server/pkg/logger"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func StartServer(lc fx.Lifecycle, gameServer *grpc_handlers.GameServer) {

	addr := configs.C.Service.Addr
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	game_server.RegisterGameServerServiceServer(s, gameServer)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("server start listening at", addr)

			go func() {
				if err := s.Serve(lis); err != nil {
					logger.Fatal("Failed to serve: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping server")
			s.GracefulStop()
			return nil
		},
	})
}

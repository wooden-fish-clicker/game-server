package grpc_handlers

import (
	"context"
	"game-server/grpc_proto/game_server"
	"game-server/internal/services"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GameServer struct {
	game_server.UnimplementedGameServerServiceServer
	clickService  *services.ClickService
	attackService *services.AttackService
}

func NewGameServer(clickService *services.ClickService, attackService *services.AttackService) *GameServer {
	return &GameServer{clickService: clickService, attackService: attackService}
}

func (a *GameServer) Attack(ctx context.Context, req *game_server.AttackRequest) (*game_server.AttackResponse, error) {

	userHp, userPoints, targetHp, targetPoints, err := a.attackService.Attack(ctx, req.GetBase().UserId, req.GetTargetId(), int(req.GetBase().Type))

	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	return &game_server.AttackResponse{
		UserInfoData: &game_server.UserInfoData{
			Hp:     userHp,
			Points: userPoints,
		},
		TargetInfoData: &game_server.TargetInfoData{
			Hp:     targetHp,
			Points: targetPoints,
		},
	}, status.Errorf(codes.OK, "")
}
func (a *GameServer) Click(ctx context.Context, req *game_server.ClickRequest) (*game_server.ClickResponse, error) {

	hp, points, err := a.clickService.Click(ctx, req.GetBase().UserId, int(req.GetBase().Type))

	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	return &game_server.ClickResponse{
		UserInfoData: &game_server.UserInfoData{
			Hp:     hp,
			Points: points,
		},
	}, status.Errorf(codes.OK, "")
}

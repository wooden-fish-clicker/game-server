package services

import (
	"context"
	"fmt"
	"game-server/internal/models"
	"game-server/internal/repository"

	"github.com/go-redis/redis/v8"
)

type attackType int

const (
	basicAttack attackType = 1 // 1

)

type AttackService struct {
	userRepository      *repository.UserRepository
	userCacheRepository *repository.UserCacheRepository
}

func NewAttackService(userRepository *repository.UserRepository, userCacheRepository *repository.UserCacheRepository) *AttackService {
	return &AttackService{userRepository: userRepository, userCacheRepository: userCacheRepository}
}

func (a *AttackService) Attack(ctx context.Context, userId string, targetId string, attackType int) (int32, int64, int32, int64, error) {
	switch attackType {
	case int(basicAttack):
		return a.basicType(ctx, userId, targetId)
	default:
		return 0, 0, 0, 0, fmt.Errorf("invalid attack type: %d", attackType)
	}

}

func (a *AttackService) basicType(ctx context.Context, userId string, targetId string) (int32, int64, int32, int64, error) {

	var (
		consumeHp     = 0
		ConsumePoints = -2
		DamageHp      = 0
		DamagePoint   = -1
	)
	userHp, userPoints, err := a.getUserState(ctx, userId)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	targetHp, targetPoints, err := a.getUserState(ctx, targetId)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	if userPoints <= int64(ConsumePoints) {
		return userHp, userPoints, targetHp, targetPoints, fmt.Errorf("user %s points 不足", userId)
	}
	if targetHp <= 0 {
		return userHp, userPoints, targetHp, targetPoints, fmt.Errorf("target %s is dead", targetId)
	}

	userHp, userPoints, targetHp, targetPoints, err = a.userCacheRepository.AdjustState(ctx, repository.Attack{
		Type:          int(basicAttack),
		UserId:        userId,
		ConsumePoints: ConsumePoints,
		ConsumeHp:     consumeHp,
		TargetId:      targetId,
		DamagePoint:   DamagePoint,
		DamageHp:      DamageHp,
	})
	if err != nil {
		return userHp, userPoints, targetHp, targetPoints, err
	}

	return userHp, userPoints, targetHp, targetPoints, nil

}

func (a *AttackService) getUserState(ctx context.Context, userId string) (int32, int64, error) {
	userHp, userPoints, _, err := a.userCacheRepository.GetUserState(ctx, userId)

	if err == redis.Nil {
		user := &models.User{
			ID: userId,
		}
		if err := a.userRepository.GetDeteil(ctx, user); err != nil {
			return 0, 0, err
		}

		if err := a.userCacheRepository.SetUserState(ctx, userId, user.UserInfo.Hp, user.UserInfo.Points); err != nil {
			return 0, 0, err
		}
		userHp = user.UserInfo.Hp
		userPoints = user.UserInfo.Points
	} else if err != nil {
		return 0, 0, err
	}

	return userHp, userPoints, err
}

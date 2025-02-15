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
	basicAttack attackType = iota + 1 // 1
	strawDoll
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
		return a.basicAttack(ctx, userId)
	case int(strawDoll):
		return a.strawDollAttack(ctx, userId, targetId)
	default:
		return 0, 0, 0, 0, fmt.Errorf("invalid attack type: %d", attackType)
	}

}

func (a *AttackService) basicAttack(ctx context.Context, userId string) (int32, int64, int32, int64, error) {

	var (
		consumePoints = -2
	)

	userHp, userPoints, err := a.getUserState(ctx, userId)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	if err := a.validateUserState(userId, userHp, userPoints, consumePoints); err != nil {
		return 0, 0, 0, 0, err
	}

	userHp, userPoints, err = a.userCacheRepository.AdjustPoints(ctx, userId, consumePoints)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return userHp, userPoints, 0, 0, nil

}

func (a *AttackService) strawDollAttack(ctx context.Context, userId string, targetId string) (int32, int64, int32, int64, error) {

	var (
		consumeHp     = 0
		consumePoints = -4
		damageHp      = 0
		damagePoint   = -1
	)
	userHp, userPoints, err := a.getUserState(ctx, userId)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	targetHp, targetPoints, err := a.getUserState(ctx, targetId)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	if err := a.validateUserState(userId, userHp, userPoints, consumePoints); err != nil {
		return 0, 0, 0, 0, err
	}

	if err := a.validateTargetState(targetId, targetHp, targetPoints); err != nil {
		return 0, 0, 0, 0, err
	}

	userHp, userPoints, targetHp, targetPoints, err = a.userCacheRepository.AdjustState(ctx, repository.Attack{
		Type:          int(basicAttack),
		UserId:        userId,
		ConsumePoints: consumePoints,
		ConsumeHp:     consumeHp,
		TargetId:      targetId,
		DamagePoint:   damagePoint,
		DamageHp:      damageHp,
	})
	if err != nil {
		return 0, 0, 0, 0, err
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

func (a *AttackService) validateUserState(userId string, userHp int32, userPoints int64, consumePoints int) error {
	if userPoints < int64(-consumePoints) {
		return fmt.Errorf("user %s points 不足", userId)
	}
	if userHp <= 0 {
		return fmt.Errorf("user %s is dead", userId)
	}
	return nil
}

func (a *AttackService) validateTargetState(targetId string, targetHp int32, targetPoints int64) error {
	if targetPoints <= 0 { //TODO 目前先不扣到hp , 所以points為0時就不繼續 , 等有重生機制的時候再改
		return fmt.Errorf("target %s points 不足", targetId)
	}
	if targetHp <= 0 {
		return fmt.Errorf("target %s is dead", targetId)
	}
	return nil
}

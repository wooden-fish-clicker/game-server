package services

import (
	"context"
	"fmt"

	"game-server/internal/models"
	"game-server/internal/repository"
)

type clickType int

const (
	basicClick clickType = 1 // 1

)

type ClickService struct {
	userRepository      *repository.UserRepository
	userCacheRepository *repository.UserCacheRepository
}

func NewClickService(userRepository *repository.UserRepository, userCacheRepository *repository.UserCacheRepository) *ClickService {
	return &ClickService{userRepository: userRepository, userCacheRepository: userCacheRepository}
}

func (c *ClickService) Click(ctx context.Context, userId string, clickType int) (int32, int64, error) {
	switch clickType {
	case int(basicClick):
		return c.basicClick(ctx, userId)
	default:
		return 0, 0, fmt.Errorf("invalid click type: %d", clickType)
	}

}

func (c *ClickService) basicClick(ctx context.Context, userId string) (int32, int64, error) {
	if err := c.checkKeyExist(ctx, userId); err != nil {
		return 0, 0, err
	}

	hp, points, err := c.userCacheRepository.AdjustPoints(ctx, userId, 1)
	if err != nil {
		return 0, 0, err
	}

	// hp, points, _, err := c.userCacheRepository.GetUserState(ctx, userId)
	// if err != nil {
	// 	return 0, 0, err
	// }

	return hp, points, nil
}

func (c *ClickService) checkKeyExist(ctx context.Context, userId string) error {
	var (
		exist bool
		err   error
	)
	if exist, err = c.userCacheRepository.CheckKeyExist(ctx, userId); err != nil {
		return err
	} else if !exist {
		//在redis裡面沒找的的話，就去mongodb撈出來，並且寫到redis中
		user := &models.User{
			ID: userId,
		}
		if err := c.userRepository.GetDeteil(ctx, user); err != nil {
			return err
		}

		if err := c.userCacheRepository.SetUserState(ctx, userId, user.UserInfo.Hp, user.UserInfo.Points); err != nil {
			return err
		}
	}

	return nil
}

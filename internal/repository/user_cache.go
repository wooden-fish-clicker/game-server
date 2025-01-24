package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	redisClient "game-server/pkg/redis"
)

const (
	userStateKey = "user:state:"
)

type Attack struct {
	Type          int
	UserId        string
	ConsumePoints int
	ConsumeHp     int
	TargetId      string
	DamagePoint   int
	DamageHp      int
}

type UserCacheRepository struct {
	redis *redisClient.Redis
}

func NewUserCacheRepository(redis *redisClient.Redis) *UserCacheRepository {
	return &UserCacheRepository{redis: redis}
}

func (u *UserCacheRepository) CheckKeyExist(ctx context.Context, userID string) (bool, error) {
	exists, err := u.redis.Client.Exists(ctx, userStateKey+userID).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (u *UserCacheRepository) AdjustPoints(ctx context.Context, userID string, count int) (int32, int64, error) {
	key := userStateKey + userID

	hp, points, err := u.adjust(ctx, key, "points", count)
	if err != nil {
		return 0, 0, err
	}

	return hp, points, nil
}

func (u *UserCacheRepository) AdjustHp(ctx context.Context, userID string, count int) (int32, int64, error) {
	key := userStateKey + userID

	hp, points, err := u.adjust(ctx, key, "hp", count)
	if err != nil {
		return 0, 0, err
	}

	return hp, points, nil
}

func (u *UserCacheRepository) adjust(ctx context.Context, key string, field string, count int) (int32, int64, error) {
	script := getHIncrByLuaScript()
	now := time.Now().Unix()

	result, err := u.redis.Client.Eval(ctx, script, []string{key},
		field,
		count,
		now).Result()

	if err != nil {
		return 0, 0, err
	}

	resultData := result.([]interface{})

	hp, points, err := hpAndPointsStrToInt(resultData[0].(string), resultData[1].(string))
	if err != nil {
		return 0, 0, err
	}
	return hp, points, nil
}
func getHIncrByLuaScript() string {
	return `local userKey = KEYS[1]
	local field = ARGV[1]
	local count = tonumber(ARGV[2])
	local now = ARGV[3]

	redis.call("HINCRBY", userKey, field, count)

	redis.call("HSET", userKey, "last_modified", now)


	local newUserState = redis.call("HMGET", userKey, "hp", "points")

	return newUserState`
}

func (u *UserCacheRepository) AdjustState(ctx context.Context, attackInfo Attack) (int32, int64, int32, int64, error) {
	script := getAdjustStateLuaScript()

	userKey := userStateKey + attackInfo.UserId
	targetKey := userStateKey + attackInfo.TargetId
	now := time.Now().Unix()

	result, err := u.redis.Client.Eval(ctx, script, []string{userKey, targetKey},
		attackInfo.ConsumePoints,
		attackInfo.ConsumeHp,
		attackInfo.DamagePoint,
		attackInfo.DamageHp,
		now).Result()

	if err != nil {
		return 0, 0, 0, 0, err
	}

	data := result.([]interface{})
	newUserState := data[0].([]interface{})
	newTargetState := data[1].([]interface{})

	newUserHp, newUserPoints, err := hpAndPointsStrToInt(newUserState[0].(string), newUserState[1].(string))
	if err != nil {
		return 0, 0, 0, 0, err
	}

	newTargetHp, newTargetPoints, err := hpAndPointsStrToInt(newTargetState[0].(string), newTargetState[1].(string))
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return newUserHp,
		newUserPoints,
		newTargetHp,
		newTargetPoints,
		nil
}

func getAdjustStateLuaScript() string {
	return `local userKey = KEYS[1]
	local targetKey = KEYS[2]
	local consumePoints = tonumber(ARGV[1])
	local consumeHp = tonumber(ARGV[2])
	local damagePoint = tonumber(ARGV[3])
	local damageHp = tonumber(ARGV[4])
	local now = ARGV[5]

	if consumePoints ~= 0 then
		redis.call("HINCRBY", userKey, "points", consumePoints)
	end

	if consumeHp ~= 0 then
		redis.call("HINCRBY", userKey, "hp", consumeHp)
	end

	if damagePoint ~= 0 then
		redis.call("HINCRBY", targetKey, "points", damagePoint)
	end

	if damageHp ~= 0 then
		redis.call("HINCRBY", targetKey, "hp", damageHp)
	end

	redis.call("HSET", userKey, "last_modified", now)
	redis.call("HSET", targetKey, "last_modified", now)

	local newUserState = redis.call("HMGET", userKey, "hp", "points")
	local newTargetState = redis.call("HMGET", targetKey, "hp", "points")

	return {newUserState, newTargetState}`
}

func (u *UserCacheRepository) SetUserState(ctx context.Context, userId string, hp int32, points int64) error {
	key := userStateKey + userId
	values := map[string]interface{}{
		"hp":            hp,
		"points":        points,
		"last_modified": time.Now().Unix(),
	}
	err := u.redis.Client.HSet(ctx, key, values).Err()
	if err != nil {
		return err
	}

	return nil
}

func (u *UserCacheRepository) GetUserState(ctx context.Context, userId string) (int32, int64, time.Time, error) {
	key := userStateKey + userId
	results, err := u.redis.Client.HMGet(ctx, key, "hp", "points", "last_modified").Result()
	if err != nil {
		return 0, 0, time.Time{}, err
	}

	hpStr, _ := results[0].(string)
	pointsStr, _ := results[1].(string)

	hp, points, err := hpAndPointsStrToInt(hpStr, pointsStr)
	if err != nil {
		return 0, 0, time.Time{}, err
	}

	lastModifiedStr, _ := results[2].(string)
	lastModifiedTimestamp, err := strconv.Atoi(lastModifiedStr)
	if err != nil {
		return 0, 0, time.Time{}, fmt.Errorf("failed to parse last_modified as time: %v", err)
	}
	lastModified := time.Unix(int64(lastModifiedTimestamp), 0)

	return hp, points, lastModified, nil
}

func hpAndPointsStrToInt(hpStr string, pointStr string) (int32, int64, error) {
	hp, err := strconv.Atoi(hpStr)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert hp to int: %v", err)
	}

	points, err := strconv.ParseInt(pointStr, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert hp to int64: %v", err)
	}

	return int32(hp), points, nil
}

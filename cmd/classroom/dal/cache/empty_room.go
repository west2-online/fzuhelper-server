package cache

import (
	"context"
	"encoding/json"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// SetEmptyRoomCache 设置空教室缓存 key 为 date + campus + startTime + endTime，value直接采用bytes存储
func SetEmptyRoomCache(ctx context.Context, key string, emptyRoomList []string) {

	emptyRoomJson, err := json.Marshal(emptyRoomList)
	// 两天过期
	err = RedisClient.Set(ctx, key, emptyRoomJson, constants.ClassroomKeyExpire).Err()
	if err != nil {
		logger.LoggerObj.Fatalf("dal.cache.SetEmptyRoomCache failed, err is %v", err)
	}
}
func GetEmptyRoomCache(ctx context.Context, key string) (emptyRoomList []string) {
	data, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		logger.LoggerObj.Fatalf("dal.cache.GetEmptyRoomCache failed, err is %v", err)
		return nil
	}
	err = json.Unmarshal([]byte(data), &emptyRoomList)
	if err != nil {
		logger.LoggerObj.Fatalf("dal.cache.GetEmptyRoomCache Unmarshal failed, err is %v", err)
		return nil
	}
	return
}
func IsExistRoomInfo(ctx context.Context, key string) bool {
	return RedisClient.Exists(ctx, key).Val() == 1
}

package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
)

type EmptyRoomCache struct {
	roomName []string
	lastTime int64
}

func SetEmptyRoomCache(ctx context.Context, key string, emptyRoomList []string) {
	emptyRoomCache := &EmptyRoomCache{
		roomName: emptyRoomList,
		lastTime: time.Now().Unix(),
	}
	emptyRoomJson, err := json.Marshal(emptyRoomCache)
	// 10分钟过期
	err = RedisClient.Set(ctx, key, emptyRoomJson, time.Minute*10).Err()
	if err != nil {
		klog.Error(err)
	}
}
func GetEmptyRoomCache(ctx context.Context, key string) (emptyRoomList []string, err error) {
	data, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(data), &emptyRoomList)
	if err != nil {
		return nil, err
	}
	return
}
func IsExistRoomInfo(ctx context.Context, key string) (exist int64, err error) {
	exist, err = RedisClient.Exists(ctx, key).Result()
	return
}

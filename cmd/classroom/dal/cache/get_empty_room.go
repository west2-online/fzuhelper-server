package cache

import (
	"context"
	"encoding/json"
	"fmt"
)

func GetEmptyRoomCache(ctx context.Context, key string) (emptyRoomList []string, err error) {
	data, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("dal.GetEmptyRoomCache: Get rooms info failed: %w", err)
	}
	err = json.Unmarshal([]byte(data), &emptyRoomList)
	if err != nil {
		return nil, fmt.Errorf("dal.GetEmptyRoomCache: Unmarshal rooms info failed: %w", err)
	}
	return emptyRoomList, nil
}

package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

func SetEmptyRoomCache(ctx context.Context, key string, emptyRoomList []string) error {
	emptyRoomJson, err := json.Marshal(emptyRoomList)
	if err != nil {
		return fmt.Errorf("dal.SetEmptyRoomCache: Marshal rooms info failed: %w", err)
	}
	err = RedisClient.Set(ctx, key, emptyRoomJson, constants.ClassroomKeyExpire).Err()
	if err != nil {
		return fmt.Errorf("dal.SetEmptyRoomCache: Set rooms info failed: %w", err)
	}
	return nil
}

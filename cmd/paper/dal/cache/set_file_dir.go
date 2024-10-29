package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func SetFileDirCache(ctx context.Context, path string, dir model.UpYunFileDir) error {
	key := getFileDirKey(path)
	data, err := sonic.Marshal(dir)
	if err != nil {
		return fmt.Errorf("dal.SetFileDirCache: Unmarshal dir info failed: %w", err)
	}
	return RedisClient.Set(ctx, key, data, 5*time.Second).Err()
}

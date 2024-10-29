package cache

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func GetFileDirCache(ctx context.Context, path string) (bool, *model.UpYunFileDir, error) {
	key := getFileDirKey(path)
	ret := &model.UpYunFileDir{}
	if !IsExists(ctx, key) {
		return false, ret, nil
	}
	data, err := RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return false, ret, fmt.Errorf("dal.GetFileDirCache: get dir info failed: %w", err)
	}
	err = sonic.Unmarshal(data, &ret)
	if err != nil {
		return false, ret, fmt.Errorf("dal.GetFileDirCache: Unmarshal dir info failed: %w", err)
	}
	return true, ret, nil
}

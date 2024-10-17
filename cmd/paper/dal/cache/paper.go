/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func SetFileDirCache(ctx context.Context, path string, dir model.UpYunFileDir) error {
	key := getFileDirKey(path)
	data, err := json.Marshal(dir)
	if err != nil {
		return fmt.Errorf("dal.SetFileDirCache: Unmarshal dir info failed: %w", err)
	}
	return RedisClient.Set(ctx, key, data, 5*time.Second).Err()
}

func GetFileDirCache(ctx context.Context, path string) (bool, *model.UpYunFileDir, error) {
	key := getFileDirKey(path)
	ret := &model.UpYunFileDir{}
	if !IsExists(ctx, key) {
		return false, ret, nil // 缓存中不存在该路径对应的目录
	}
	data, err := RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return false, ret, fmt.Errorf("dal.GetFileDirCache: get dir info failed: %w", err)
	}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return false, ret, fmt.Errorf("dal.GetFileDirCache: Unmarshal dir info failed: %w", err)
	}
	return true, ret, nil
}

func IsExists(ctx context.Context, key string) bool {
	return RedisClient.Exists(ctx, key).Val() == 1
}

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

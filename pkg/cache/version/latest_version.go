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

package version

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/pkg/base/environment"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
)

func (c *CacheVersion) GetLatestVersionCache(ctx context.Context, versionType string) (*model.VersionHistory, error) {
	data, err := c.client.Get(ctx, constants.LatestVersionCachePrefix+versionType).Bytes()
	if err != nil {
		return nil, fmt.Errorf("dal.GetLatestVersionCache: cache failed: %w", err)
	}
	vh := new(model.VersionHistory)
	if err = sonic.Unmarshal(data, vh); err != nil {
		return nil, fmt.Errorf("dal.GetLatestVersionCache: Unmarshal failed: %w", err)
	}
	return vh, nil
}

func (c *CacheVersion) SetLatestVersionCache(ctx context.Context, vh *model.VersionHistory) error {
	if environment.IsTestEnvironment() {
		return nil
	}
	data, err := sonic.Marshal(vh)
	if err != nil {
		return fmt.Errorf("dal.SetLatestVersionCache: Marshal failed: %w", err)
	}
	key := constants.LatestVersionCachePrefix + vh.Type
	if err = c.client.Set(ctx, key, data, constants.LatestVersionKeyExpire).Err(); err != nil {
		return fmt.Errorf("dal.SetLatestVersionCache: Set cache failed: %w", err)
	}
	return nil
}

func (c *CacheVersion) DeleteLatestVersionCache(ctx context.Context, versionType string) error {
	if err := c.client.Del(ctx, constants.LatestVersionCachePrefix+versionType).Err(); err != nil {
		return fmt.Errorf("dal.DeleteLatestVersionCache: cache failed: %w", err)
	}
	return nil
}

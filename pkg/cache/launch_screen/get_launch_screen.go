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

package launch_screen

import (
	"context"
	"fmt"
	"strings"

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

func (c *CacheLaunchScreen) GetLaunchScreenCache(ctx context.Context, key string) (pictureIdList []int64, err error) {
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("dal.GetLaunchScreenCache: Get pictureIdList cache failed: %w", err)
	}

	if err = sonic.Unmarshal([]byte(data), &pictureIdList); err != nil {
		return nil, fmt.Errorf("dal.GetLaunchScreenCache: Unmarshal pictureIdList failed: %w", err)
	}

	return pictureIdList, nil
}

func (c *CacheLaunchScreen) GetLastLaunchScreenIdCache(ctx context.Context, device string) (int64, error) {
	id, err := c.client.Get(ctx, strings.Join([]string{device, constants.LastLaunchScreenIdKey}, ":")).Int64()
	if err != nil {
		return -1, fmt.Errorf("dal.GetLaunchScreenCache: Get LastLaunchScreenId cache failed: %w", err)
	}

	return id, nil
}

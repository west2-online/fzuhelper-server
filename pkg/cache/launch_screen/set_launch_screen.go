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

func (c *CacheLaunchScreen) SetLaunchScreenCache(ctx context.Context, key string, pictureIdList *[]int64) error {
	pictureIdListJson, err := sonic.Marshal(pictureIdList)
	if err != nil {
		return fmt.Errorf("dal.SetLaunchScreenCache: Marshal pictureIdList failed: %w", err)
	}
	if err = c.client.Set(ctx, key, pictureIdListJson, constants.LaunchScreenKeyExpire).Err(); err != nil {
		return fmt.Errorf("dal.SetLaunchScreenCache: Set pictureIdList cache failed: %w", err)
	}
	return nil
}

func (c *CacheLaunchScreen) SetLastLaunchScreenIdCache(ctx context.Context, id int64, device string) error {
	if err := c.client.Set(ctx, strings.Join([]string{device, constants.LastLaunchScreenIdKey}, ":"), id, constants.LaunchScreenKeyExpire).Err(); err != nil {
		return fmt.Errorf("dal.SetTotalLaunchScreenCountCache failed: %w", err)
	}
	return nil
}

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

package classroom

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

// SetAvailabilityCache 保存当天教室数据，nil 表示刷新失败，空切片表示没有教室
func (c *CacheClassroom) SetAvailabilityCache(ctx context.Context, key string, rooms []*model.Classroom) error {
	data, err := json.Marshal(rooms)
	if err != nil {
		return fmt.Errorf("dal.SetAvailabilityCache: Marshal failed: %w", err)
	}
	if err := c.client.Set(ctx, key, data, constants.ClassroomKeyExpire).Err(); err != nil {
		return fmt.Errorf("dal.SetAvailabilityCache: set cache failed: %w", err)
	}
	return nil
}

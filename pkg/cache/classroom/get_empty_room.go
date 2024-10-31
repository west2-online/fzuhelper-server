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
)

func (c *CacheClassroom) GetEmptyRoomCache(ctx context.Context, key string) (emptyRoomList []string, err error) {
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("dal.GetEmptyRoomCache: Get rooms info failed: %w", err)
	}
	err = json.Unmarshal([]byte(data), &emptyRoomList)
	if err != nil {
		return nil, fmt.Errorf("dal.GetEmptyRoomCache: Unmarshal rooms info failed: %w", err)
	}
	return emptyRoomList, nil
}

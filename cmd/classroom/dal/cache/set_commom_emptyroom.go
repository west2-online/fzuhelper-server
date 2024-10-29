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

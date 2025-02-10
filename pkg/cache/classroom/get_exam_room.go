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
	"github.com/west2-online/fzuhelper-server/pkg/errno"

	"github.com/bytedance/sonic"

	"github.com/west2-online/jwch"
)

func (c *CacheClassroom) GetExamRoom(ctx context.Context, key string) ([]*jwch.ExamRoomInfo, error) {
	ret := make([]*jwch.ExamRoomInfo, 0)
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "dal.GetExamRoom: Get exam rooms info failed: %w", err)
	}
	err = sonic.Unmarshal([]byte(data), &ret)
	if err != nil {
		return nil, errno.Errorf(errno.InternalJSONErrorCode, "dal.GetExamRoom: Unmarshal exam rooms info failed: %w", err)
	}
	return ret, nil
}

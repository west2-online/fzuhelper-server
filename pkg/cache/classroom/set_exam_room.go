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

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/jwch"
)

func (c *CacheClassroom) SetExamRoom(ctx context.Context, key string, value []*jwch.ExamRoomInfo) {
	examRoomJson, err := sonic.Marshal(value)
	if err != nil {
		logger.Errorf("dal.SetExamRoom: marshal exam room info failed, err: %v", err)
		return
	}
	err = c.client.Set(ctx, key, examRoomJson, constants.ExamRoomKeyExpire).Err()
	if err != nil {
		logger.Errorf("dal.SetExamRoom: set exam room failed, err: %v", err)
		return
	}
}

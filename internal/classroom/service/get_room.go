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

package service

import (
	"errors"
	"fmt"

	"github.com/west2-online/fzuhelper-server/internal/classroom/dal/cache"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
)

func (s *ClassroomService) GetEmptyRoom(req *classroom.EmptyRoomRequest) ([]string, error) {
	// 从redis中获取数据
	key := fmt.Sprintf("%s.%s.%s.%s", req.Date, req.Campus, req.StartTime, req.EndTime)
	if ok := cache.IsExistRoomInfo(s.ctx, key); !ok {
		return nil, errors.New("service.GetEmptyRoom: room info not exist")
	}
	emptyRoomList, err := cache.GetEmptyRoomCache(s.ctx, key)
	if err != nil {
		return nil, fmt.Errorf("service.GetEmptyRoom: Get room info failed: %w", err)
	}
	return emptyRoomList, nil
}

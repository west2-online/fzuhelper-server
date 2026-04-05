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
	"fmt"
	"time"

	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

const (
	HoursInADay = 24
	MinDateDiff = 0
	MaxDateDiff = 7
)

func (s *ClassroomService) GetEmptyRoom(req *classroom.EmptyRoomRequest) ([]string, error) {
	// 判断req.date只能从今天开始的七天内，在当前日期前或超过 7 天则报错
	// 首先判断date的格式是否符合要求
	requestDate, err := utils.TimeParse(req.Date)
	if err != nil {
		return nil, errno.ErrNoWithPreMessage(err, "Classroom.GetEmptyRoom: Date format failed")
	}
	now := time.Now().Truncate(constants.ONE_DAY)
	requestDate = requestDate.Truncate(constants.ONE_DAY)
	dateDiff := requestDate.Sub(now).Hours() / HoursInADay
	if dateDiff < MinDateDiff || dateDiff > MaxDateDiff {
		return nil, errno.Errorf(errno.InternalServiceErrorCode, "Classroom.GetEmptyRoom: Date out of range: %v", req.Date)
	}
	// 从redis中获取数据
	key := fmt.Sprintf("%s.%s.%s.%s", req.Date, req.Campus, req.StartTime, req.EndTime)
	if !s.cache.IsKeyExist(s.ctx, key) {
		return nil, errno.Errorf(errno.InternalRedisErrorCode, "Classroom.GetEmptyRoom: Room info not exist")
	}
	emptyRoomList, err := s.cache.Classroom.GetEmptyRoomCache(s.ctx, key)
	if err != nil {
		return nil, errno.ErrNoWithPreMessage(err, "Classroom.GetEmptyRoom: Get room info failed")
	}
	return emptyRoomList, nil
}
